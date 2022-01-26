package miner

import (
	"os/exec"
	"strings"

	"xdt.com/hm-diag/diag/jsonrpc"
	"xdt.com/hm-diag/util"
)

type PeerBookParams struct {
	Addr string `json:"addr"`
}

type PeerBook struct {
	Address         string   `json:"address"`
	ConnectionCount uint     `json:"connectionCount"`
	LastUpdated     uint     `json:"lastUpdated"`
	ListenAddrCount uint     `json:"listenAddrCount"`
	ListenAddresses []string `json:"listenAddresses"`
	Name            string   `json:"name"`
	Nat             string   `json:"nat"`
	Sessions        struct {
		Local  string `json:"local"`
		Name   string `json:"name"`
		P2p    string `json:"p2p"`
		Remote string `json:"remote"`
	} `json:"sessions"`
}

type InfoSummary struct {
	FirmwareVersion string `json:"firmwareVersion"`
	Height          string `json:"height"`
	Name            string `json:"name"`
	Version         string `json:"version"`
}

type InfoP2pStatus struct {
	Connected string `json:"connected"`
	Dialable  string `json:"dialable"`
	Height    uint64 `json:"height"`
	NatType   string `json:"natType"`
}

type MinerInfo struct {
	InfoHeight    uint64        `json:"infoHeight"`
	InfoP2pStatus InfoP2pStatus `json:"infoP2pStatus"`
	InfoRegion    string        `json:"infoRegion"`
	InfoSummary   InfoSummary   `json:"infoSummary"`
	PeerAddr      string        `json:"peerAddr"`
	PeerBook      PeerBook      `json:"peerBook"`
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
	ver, _ := Version()
	res, _ = client.Call("info_name", nil)
	name := ""
	if r, ok := res.(map[string]interface{}); ok {
		name = r["name"].(string)
	}
	resMap["infoSummary"] = map[string]interface{}{
		"firmwareVersion": fwVer,
		"version":         ver,
		"name":            name,
		"height":          resMap["infoHeight"],
	}

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

func Version() (string, error) {
	s := "cat /etc/lsb_release | grep DISTRIB_RELEASE | awk -F = '{print $2}'"
	cmd := exec.Command("sh", "-c", s)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(string(out), "\n", ""), nil
}
