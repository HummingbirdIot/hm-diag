package link

import (
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type Closehandler func(code int, text string) error
type Connecthandler func() error
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

func (c *UplinkConn) Start() error {
	return c.connect()
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

func (c *UplinkConn) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	wsConn, resp, err := websocket.DefaultDialer.Dial(c.config.server, c.config.header)
	if err != nil {
		c.httpResp = resp
		c.connectErr = err
		return errors.WithMessage(err, "dial connection")
	}
	wsConn.SetCloseHandler(c.wsCloseHandler)
	wsConn.SetPongHandler(func(d string) error {
		log.Println("got pong", d)
		return nil
	})

	c.connected = true
	c.ws = wsConn
	go c.onConnectHandler()
	return nil
}

func (c *UplinkConn) wsCloseHandler(code int, text string) error {
	log.Printf("ws closed, code: %d, text: %s", code, text)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.connected = false
	if c.config.autoReconnect {
		c.reconnect()
	}
	if c.onCloseHandler != nil {
		c.onCloseHandler(code, text)
	}
	return nil
}

func (c *UplinkConn) reconnect() error {
	// TODO: reconnect
	log.Println("reconnect...")
	err := c.connect()
	return err
}
