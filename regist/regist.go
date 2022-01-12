package regist

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"xdt.com/hm-diag/config"
	"xdt.com/hm-diag/diag/device"
)

var defaultRegistryHost = "hiot-registry"
var httpClient = http.Client{}

type Register struct {
	ApiPort          int
	RegistryApiPort  int
	ReistIntervalSec int64
}

var reg *Register = nil

func initInstance(conf *config.GlobalConfig) {
	reg = &Register{
		ApiPort:          conf.ApiPort,
		RegistryApiPort:  6753,
		ReistIntervalSec: 30,
	}
}

func Instance() *Register {
	return reg
}

func StartRegistJob() {
	if reg == nil {
		initInstance(config.Config())
	}
	log.Printf("regist job scheduler start")

	ticker := time.NewTicker(time.Duration(reg.ReistIntervalSec) * time.Second)
	quitTask := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				reg.Do()
			case <-quitTask:
				ticker.Stop()
			}
		}
	}()

	reg.Do()
}

func (r *Register) Do() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("do regist error", r)
		}
	}()
	registry, err := registryIp()
	if err == nil {
		registryApi := "http://" + registry.String() + ":" + strconv.Itoa(r.RegistryApiPort) + "/regist"
		info, _ := r.GetRegistInfo()
		j, _ := json.Marshal(info)
		req, err := http.NewRequest("POST", registryApi, bytes.NewReader(j))
		if err != nil {
			log.Println("[error] request regist api error:", err)
			return
		}
		resp, err := httpClient.Do(req)
		if err != nil {
			log.Println("[error] request error:", err)
			return
		}
		if resp.StatusCode >= 300 {
			log.Println("[error] request regist api error, wrong http status", resp.StatusCode)
		}
	} else {
		// log.Println("[error] get default registry error:", err)
	}
}

func registryIp() (*net.IPAddr, error) {
	ip, err := net.ResolveIPAddr("", defaultRegistryHost)
	return ip, err
}

func (r *Register) GetRegistInfo() (interface{}, error) {
	m := make(map[string]interface{})
	hardAddr, ipAddrs, err := device.GetAddrs("eth0")
	if err != nil {
		return nil, err
	}
	m["eth0"] = map[string]interface{}{"device_addr": hardAddr.String(), "ip": ipAddrs}
	hardAddr, ipAddrs, err = device.GetAddrs("wlan0")
	if err != nil {
		return nil, err
	}
	m["wlan0"] = map[string]interface{}{"device_addr": hardAddr.String(), "ip": ipAddrs}

	sn, err := device.GetSn()
	if err != nil {
		sn = "unknown"
		log.Println("[error] can't get sn:", err)
	}
	m["serialNumber"] = sn
	m["api"] = map[string]interface{}{
		"port":   r.ApiPort,
		"schema": "http",
	}

	return m, nil
}
