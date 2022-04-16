package link

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/pkg/errors"
	"xdt.com/hm-diag/link/message"
)

type Client struct {
	config ClientConfig
	conn   *UplinkConn
}

const (
	SECURITY_KEY_HEADER = "X-SECRET-KEY"
	SECURITY_ID         = "X-SECRET-ID"
)

func newClient(config ClientConfig) (*Client, error) {
	return &Client{config: config}, nil
}

func (c *Client) Start(ctx context.Context) error {
	if c.config.Secret == "" || c.config.Server == "" || c.config.ID == "" {
		return fmt.Errorf("secret and server ad id must be provided")
	}

	if c.conn != nil {
		go c.conn.Close()
	}

	conn, err := NewConn(ConnConfig{
		header: map[string][]string{
			SECURITY_KEY_HEADER: {c.config.Secret},
			SECURITY_ID:         {c.config.ID},
		},
		server:        c.config.Server,
		autoReconnect: true,
	})
	if err != nil {
		return errors.WithMessage(err, "NewConn start")
	}

	conn.SetConnectHandler(c.onConnectHandler)
	conn.SetReadMessageHandler(c.read)
	c.conn = conn
	err = c.conn.Start(ctx)
	if err != nil {
		return errors.WithMessage(err, "client start")
	}
	go c.conn.keepPing(ctx)
	return nil
}

func (c *Client) WriteMessage(msg interface{}) error {
	if !c.conn.Connected() {
		return fmt.Errorf("connection is not established")
	}
	c.conn.mu.Lock()
	defer c.conn.mu.Unlock()
	err := c.conn.WriteJSON(msg)
	if err != nil {
		log.Println("error writing message: ", err)
		fmt.Println(msg)
		return err
	}
	return nil
}

func (c *Client) read(ctx context.Context) {
	log.Println("started read message loop")
	// TODO: for error
	for {
		if !c.conn.Connected() {
			log.Println("disconnected, give up reading")
			break
		}
		if ctx.Err() != nil {
			log.Println("break read message loop cause context:", ctx.Err())
			break
		}

		log.Println("to read message ...")
		_, buf, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("error reading message: ", err)
			// c.conn.wsCloseHandler(ctx, 1006, err.Error())
			break
		} else {
			// log.Println("go ws message:", string(buf))
			var msg map[string]interface{}
			err := json.Unmarshal(buf, &msg)
			if err != nil {
				log.Printf("invalid message: %s, err: %s", string(buf), err)
				continue
			}
			// TODO: pool
			go c.handleMessage(&msg, buf)
		}
	}
}

func (c *Client) handleMessage(msg *map[string]interface{}, rawBuf []byte) {
	t, err := message.Typeof(msg)
	if err != nil {
		log.Printf("unknown message type for %#v, error: %s\n", msg, err)
	}
	switch t {
	case message.TYPE_HTTP_REQUEST:
		var r message.HttpRequest
		err := json.Unmarshal(rawBuf, &r)
		if err != nil {
			log.Println("invalid message for http request")
			return
		}
		c.handleRpcRequest(&r)
	default:
		log.Println("unknown message type:", t)
	}
}

func (c *Client) handleRpcRequest(r *message.HttpRequest) {
	// log.Printf("handle message: %#v", r)
	resp, err := requestLocal(r)
	if err != nil {
		log.Printf("do request error: %v", err)
		return
	}
	c.WriteMessage(resp)
}

func (c *Client) onConnectHandler(ctx context.Context) error {
	log.Println("in client onConnectHandler")
	return nil
}
