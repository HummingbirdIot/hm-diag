package link

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
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
	return c.connect(ctx)
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
	fmt.Println("ws server: " + c.config.server)
	wsConn, resp, err := websocket.DefaultDialer.Dial(c.config.server, c.config.header)
	if err != nil {
		fmt.Println("connect ws server error: ", err)
		c.httpResp = resp
		c.connectErr = err
		return errors.WithMessage(err, "dial connection")
	}
	wsConn.SetCloseHandler(func(code int, text string) error {
		return c.wsCloseHandler(ctx, code, text)
	})
	wsConn.SetPongHandler(func(d string) error {
		log.Println("got pong", d)
		return nil
	})

	c.connected = true
	c.ws = wsConn
	go c.keepPing(ctx)
	go c.onConnectHandler(ctx)

	return nil
}

//ping
func (c *UplinkConn) keepPing(ctx context.Context) {
	// Keep connection alive
	log.Println("ws connection keep ping start")
	for {
		select {
		case <-ctx.Done():
			log.Println("kepp ping break case context done")
			break
		case <-time.After(30 * time.Second):
		}

		time.Sleep(30 * time.Second)
		ctx.Done()
		log.Println("ping ws server:", c.config.server)
		err := c.ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second))
		if err != nil {
			log.Println("ws connection ping error", err)
			c.ws.Close()
			if c.config.autoReconnect {
				c.reconnect(ctx)
			}
		}
	}
}

//ws关闭时处理函数
func (c *UplinkConn) wsCloseHandler(ctx context.Context, code int, text string) error {
	log.Printf("ws closed, code: %d, text: %s", code, text)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.connected = false
	if c.config.autoReconnect {
		c.reconnect(ctx)
	}
	if c.onCloseHandler != nil {
		c.onCloseHandler(code, text)
	}
	return nil
}

//ws重连
func (c *UplinkConn) reconnect(ctx context.Context) error {
	// TODO: reconnect
	log.Println("reconnect...")
	err := c.connect(ctx)
	return err
}
