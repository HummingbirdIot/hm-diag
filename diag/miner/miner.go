package miner

import (
	_ "embed"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"xdt.com/hm-diag/config"
	"xdt.com/hm-diag/diag/jsonrpc"
	"xdt.com/hm-diag/util"
)

var (
	//go:embed hotspot-info.sh
	hotspotInfoScript string
	hotspotInfoReg    = regexp.MustCompile(`"location: (\S*)",.*"owner: /p2p/(.*)"`)
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

type BaseInfo struct {
	Owner    string `json:"owner"`
	Location string `json:"location"`
}

func FetchData(url string) map[string]interface{} {
	// TODO 去掉 map，采用强类型

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

	baseInfo, err := HotspotBaseInfo()
	if err != nil {
		log.Println("get hotspot base info error:", err)
	}

	infoSummary := map[string]interface{}{
		"firmwareVersion": fwVer,
		"version":         ver,
		"name":            name,
		"height":          resMap["infoHeight"],
		"location":        baseInfo.Location,
		"owner":           baseInfo.Owner,
	}

	resMap["infoSummary"] = infoSummary
	// res, _ = client.Call("print_keys", nil)
	// resMap["print_keys"] = res

	uptime, err := Uptime()
	if err == nil {
		upsec := int64(uptime.Seconds())
		infoSummary["uptime"] = upsec
	} else {
		log.Println("get miner uptime error:", err)
	}

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

func Uptime() (time.Duration, error) {
	resultPre := ">>>result:"
	cmd := exec.Command(config.MAIN_SCRIPT, "minerStartTime")
	cmd.Dir = config.Config().GitRepoDir
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	lines := strings.Split(string(out), "\n")
	var startedAt time.Time
	for _, line := range lines {
		if strings.HasPrefix(line, resultPre) {
			str := strings.TrimPrefix(line, resultPre)
			str = str[:19] // remove ms、zone
			t, err := util.ParseTimeInLocation("UTC", "2006-01-02T15:04:05", str)
			if err != nil {
				log.Println("get miner start time error", err)
				return 0, err
			} else {
				startedAt = t
			}
			break
		}
	}
	sec := time.Since(startedAt)
	return sec, nil
}

func PacketForwardVersion() (string, error) {
	resultPre := ">>>result:"
	log.Println("exec cmd:", config.MAIN_SCRIPT, "pktfwdVersion")
	cmd := exec.Command(config.MAIN_SCRIPT, "pktfwdVersion")
	cmd.Dir = config.Config().GitRepoDir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, resultPre) {
			v := strings.TrimPrefix(line, resultPre)
			return v, nil
		}
	}
	return "", fmt.Errorf("can't get version")
}

func HotspotBaseInfo() (BaseInfo, error) {
	cmd := exec.Command("bash", "-c", hotspotInfoScript)
	buf, err := cmd.Output()
	if err != nil {
		return BaseInfo{}, err
	}

	// str := `[\loc: 8c309948840dbff", "owner: /p2p/144ZZ7GRjFuQNhSuPz7UQzagGXoD7ohiMyHVW6LbY3SEyjNvEjG"]`
	str := string(buf)
	matches := hotspotInfoReg.FindStringSubmatch(str)
	if len(matches) != 3 {
		return BaseInfo{}, fmt.Errorf("wrong result")
	}

	return BaseInfo{Location: matches[1], Owner: matches[2]}, nil
}
