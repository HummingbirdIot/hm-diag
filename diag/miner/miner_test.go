package miner

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/kpango/glg"
)

const (
	url      = "http://127.0.0.1:4467"
	peerAddr = "/p2p/11d4avV1z1Zz6ACbXvVpz4L8zn8VUmxHJ7GnLLEW8dKyC1nQ5qV"
)

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

func TestGetPeerAddr(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	glg.Info("httpMock")
	httpmock.RegisterResponder("POST", url, func(r *http.Request) (*http.Response, error) {
		buffer, _ := ioutil.ReadAll(r.Body)
		var req Req
		json.Unmarshal(buffer, &req)
		glg.Info("<<<<<<<<<<<<", req.Method)
		switch req.Method {
		case "info_height":
			return httpmock.NewStringResponse(200, string("13456")), nil
		case "peer_addr":
			res := `{
				"id":"1",
				"jsonrpc":"2.0",
				"result":{"peer_addr":"/p2p/11d4avV1z1Zz6ACbXvVpz4L8zn8VUmxHJ7GnLLEW8dKyC1nQ5qV"}
			}`
			return httpmock.NewStringResponse(200, res), nil
		}
		return nil, nil
	})
	addr := GetPeerAddr(url)
	if addr != peerAddr {
		t.Error("invaild perrAddr")
	}
}
