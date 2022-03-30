package link

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"xdt.com/hm-diag/diag"
	"xdt.com/hm-diag/link/message"
)

type Client struct {
	config ClientConfig
	conn   *UplinkConn
}

type ClientConfig struct {
	ID     string `json:"id"`
	Auth   string `json:"auth"`
	Server string `json:"server"`
}

const (
	SECURITY_KEY_HEADER = "X-SECRET-KEY"
	SECURITY_ID         = "X-SECRET-ID"
)

func NewClient(config ClientConfig) (*Client, error) {
	if config.Auth == "" || config.Server == "" {
		return nil, fmt.Errorf("auth and server must be provided")
	}
	return &Client{config: config}, nil
}

func (c *Client) Start() error {
	conn, err := NewConn(ConnConfig{
		header: map[string][]string{
			SECURITY_KEY_HEADER: {c.config.Auth},
			SECURITY_ID:         {c.config.ID},
		},
		server:        c.config.Server,
		autoReconnect: true,
	})
	if err != nil {
		return errors.WithMessage(err, "client start")
	}
	conn.SetConnectHandler(c.onConnectHandler)
	c.conn = conn
	err = c.conn.Start()
	if err != nil {
		return errors.WithMessage(err, "client start")
	}
	return nil
}

func (c *Client) WriteMessage(msg interface{}) error {
	err := c.conn.WriteJSON(msg)
	if err != nil {
		log.Println("error writing message: ", err)
	}
	return nil
}

func (c *Client) read() {
	// TODO: for error
	for {
		if !c.conn.Connected() {
			log.Println("disconnected, give up reading")
			break
		}
		_, buf, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("error reading message: ", err)
		} else {
			var msg message.Msg[any]
			err := json.Unmarshal(buf, &msg)
			if err != nil {
				log.Printf("invalid message: %s, err: %s", string(buf), err)
				continue
			}
			// TODO: pool
			go c.handleMessage(&msg)
		}
	}
}

func (c *Client) handleMessage(msg *message.Msg[any]) {
	t := message.Typeof(msg)
	switch t {
	case message.TYPE_HTTP_REQUEST:
		c.handleRpcRequest(&msg)
	default:
		log.Println("unknown message type:", t)
	}
}

func (c *Client) handleRpcRequest(msg interface{}) {
	m, ok := msg.(*message.HttpRequest)
	if !ok {
		log.Printf("invalid rpc request message, type: %#v", msg)
		return
	}
	resp, err := requestLocal(m)
	if err != nil {
		log.Printf("do request error: %v", err)
		return
	}
	c.WriteMessage(resp)
}

func (c *Client) onConnectHandler() error {
	if diag.TaskInstance() == nil {
		// wait for init diag task
		time.Sleep(time.Second * 1)
		c.onConnectHandler()
	}
	d := diag.TaskInstance().Data()
	if d.FetchTime.IsZero() {
		// wait for load data
		time.Sleep(time.Second * 1)
		c.onConnectHandler()
	} else {
		c.WriteMessage(message.OfReportData(d))
	}

	return nil
}
