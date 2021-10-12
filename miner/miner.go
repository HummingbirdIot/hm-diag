package miner

import (
	"xdt.com/hm-diag/jsonrpc"
)

type PeerBookParams struct {
	Addr string `json:"addr"`
}

func FetchData(url string) map[string]interface{} {
	client := &jsonrpc.Client{Url: url}
	resMap := make(map[string]interface{})
	res, _ := client.Call("info_height", nil)
	resMap["info_height"] = res.(map[string]interface{})["height"]

	res, _ = client.Call("info_region", nil)
	resMap["info_region"] = res.(map[string]interface{})["region"]

	res, _ = client.Call("peer_addr", nil)
	resMap["peer_addr"] = res.(map[string]interface{})["peer_addr"]

	res, _ = client.Call("peer_book", PeerBookParams{Addr: "self"})
	resMap["peer_book"] = res

	// res, _ = client.Call("info_summary", nil)
	// resMap["info_summary"] = res

	return resMap
}
