package link

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kpango/glg"
	"github.com/pkg/errors"
)

var (
	connCtx       context.Context
	connCtxCancel context.CancelFunc
)

type Closehandler func(code int, text string) error
type Connecthandler func(ctx context.Context) error
type UplinkConn struct {
	mu               sync.RWMutex
	config           ConnConfig
	ws               *websocket.Conn
	connectErr       error
	httpResp         *http.Response
	connected        bool
	onCloseHandler   Closehandler
	onConnectHandler Connecthandler
}

type ConnConfig struct {
	header        http.Header
	server        string
	autoReconnect bool
}

func NewConn(conf ConnConfig) (*UplinkConn, error) {
	_, err := url.Parse(conf.server)
	if err != nil {
		return nil, errors.WithMessage(err, "server address")
	}

	c := &UplinkConn{
		config: conf,
	}
	return c, nil
}

func (c *UplinkConn) Start(ctx context.Context) error {
	err := c.connect(ctx)
	if err != nil {
		return err
	}
	go c.keepPing(ctx)
	return nil

}

func (c *UplinkConn) Close() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in close client connection", r)
		}
	}()
	c.config.autoReconnect = false
	c.ws.Close()
}

func (c *UplinkConn) SetCloseHandler(h Closehandler) {
	c.onCloseHandler = h
}

func (c *UplinkConn) SetConnectHandler(h Connecthandler) {
	c.onConnectHandler = h
}

func (c *UplinkConn) WriteJSON(v interface{}) error {
	return c.ws.WriteJSON(v)
}

func (c *UplinkConn) ReadJSON(v interface{}) error {
	return c.ws.ReadJSON(v)
}

func (c *UplinkConn) ReadMessage() (messageType int, p []byte, err error) {
	return c.ws.ReadMessage()
}

func (c *UplinkConn) Connected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

//ws连接
func (c *UplinkConn) connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	glg.Info("ws server: " + c.config.server)
	wsConn, resp, err := websocket.DefaultDialer.Dial(c.config.server, c.config.header)
	if err != nil {
		glg.Error("connect ws server error: ", err)
		c.httpResp = resp
		c.connectErr = err
		return errors.WithMessage(err, "dial connection")
	}
	wsConn.SetCloseHandler(func(code int, text string) error {
		return c.wsCloseHandler(ctx, code, text)
	})
	wsConn.SetPongHandler(func(d string) error {
		glg.Debug("got pong", d)
		return nil
	})

	c.connected = true
	c.ws = wsConn
	connCtx, connCtxCancel = context.WithCancel(ctx)
	//connect 函数内的协程都有connCtx管理，属于ctx下的context（context销毁时子context也会销毁）
	go c.onConnectHandler(connCtx)
	return nil
}

//ping
func (c *UplinkConn) keepPing(ctx context.Context) {
	// Keep connection alive
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	intervalTime := 30 + r.Intn(10) + r.Intn(3)*10
	glg.Info("ws connection keep ping start interval : ", intervalTime)

	for {
		select {
		case <-ctx.Done():
			glg.Warn("kepp ping break case context done")
			return
		case <-time.After(time.Duration(intervalTime) * time.Second):
			glg.Debug("ping ws server:", c.config.server)
			err := c.ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second))
			if err != nil {
				glg.Error("ws connection ping error", err)
				c.wsCloseHandler(ctx, 1006, err.Error())
			}
		}
	}
}

//ws关闭时处理函数
func (c *UplinkConn) wsCloseHandler(ctx context.Context, code int, text string) error {
	glg.Info("ws closed, code: %d, text: %s", code, text)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.connected = false
	c.ws.Close()
	connCtxCancel()
	if c.config.autoReconnect {
		go func() {
			err := c.reconnect(ctx)
			if err != nil {
				glg.Error("reconnect error : ", err)
			}
			return
		}()
	}
	if c.onCloseHandler != nil {
		c.onCloseHandler(code, text)
	}
	return nil
}

//ws重连
func (c *UplinkConn) reconnect(ctx context.Context) error {
	// TODO: reconnect
	glg.Info("reconnect...")
	err := c.connect(ctx)
	return err
}
