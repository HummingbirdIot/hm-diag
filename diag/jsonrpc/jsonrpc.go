package jsonrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Client struct {
	Url string
}

type Req struct {
	Id      string      `json:"id"`
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type Resp struct {
	Id      string      `json:"id"`
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   interface{} `json:"error"`
}

func (c Client) Call(method string, params interface{}) (result interface{}, err error) {
	// jsonBuf := []byte(`{"method": "info_height", "jsonrpc": "2.0", "id": 1}`)
	req := Req{Id: "1", Jsonrpc: "2.0", Method: method, Params: params}

	jsonBuf, _ := json.Marshal(req)
	resp, err := http.Post(c.Url, "application/json", bytes.NewReader(jsonBuf))
	if err != nil {
		log.Println("jsonrpc error resp: ", err, " jsonrpc call: ", c.Url, string(jsonBuf))
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bodyStr := string(body)
	var res map[string]interface{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	if res["error"] != nil {
		log.Println("jsonrpc error: ", string(bodyStr), " jsonrpc call: ", c.Url, string(jsonBuf))
		return nil, fmt.Errorf(bodyStr)
	}

	return res["result"], nil
}
