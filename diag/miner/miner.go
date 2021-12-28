package miner

import (
	"os/exec"

	"xdt.com/hm-diag/diag/jsonrpc"
	"xdt.com/hm-diag/util"
)

type PeerBookParams struct {
	Addr string `json:"addr"`
}

func FetchData(url string) map[string]interface{} {
	client := &jsonrpc.Client{Url: url}
	resMap := make(map[string]interface{})
	res, _ := client.Call("info_height", nil)
	resMap["infoHeight"] = res.(map[string]interface{})["height"]

	res, _ = client.Call("info_region", nil)
	resMap["infoRegion"] = res.(map[string]interface{})["region"]

	res, _ = client.Call("peer_addr", nil)
	resMap["peerAddr"] = res.(map[string]interface{})["peer_addr"]

	res, _ = client.Call("peer_book", PeerBookParams{Addr: "self"})
	resMap["peerBook"] = util.ToLowerCamelObj(res)

	res, _ = client.Call("info_p2p_status", nil)
	resMap["infoP2pStatus"] = util.ToLowerCamelObj(res)

	//res, _ = client.Call("info_summary", nil)
	//resMap["infoSummary"] = util.ToLowerCamelObj(res)

	fwVer, _ := FirmwareVersion()
	resMap["infoSummary"] = map[string]string{"firmwareVersion": fwVer}

	// res, _ = client.Call("print_keys", nil)
	// resMap["print_keys"] = res

	return resMap
}

func FirmwareVersion() (string, error) {
	cmd := exec.Command("cat", "/etc/lsb_release")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
