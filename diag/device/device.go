package device

import (
	"log"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/godbus/dbus/v5"
	"xdt.com/hm-diag/util"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	unet "github.com/shirou/gopsutil/v3/net"
)

type DiskInfo struct {
	Free        uint64  `json:"free"`
	Fstype      string  `json:"fstype"`
	Path        string  `json:"path"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

type NetInterfaceInfo struct {
	Index        uint   `json:"index"`
	Mtu          uint   `json:"mtu"`
	Name         string `json:"name"`
	HardwareAddr string `json:"hardwareAddr"`
	Addrs        []struct {
		Addr string `json:"addr"`
	} `json:"addrs"`
}

type WifiInfo struct {
	Connected bool   `json:"connected"`
	Name      string `json:"name"`
	Powered   bool   `json:"powered"`
}
type MemInfo struct {
	Available uint64 `json:"available"`
	Buffers   uint64 `json:"buffers"`
	Cached    uint64 `json:"cached"`
	Free      uint64 `json:"free"`
	Shared    uint64 `json:"shared"`
	Total     uint64 `json:"total"`
	Used      uint64 `json:"used"`
}
type HostInfo struct {
	Hostname string `json:"hostname"`
	Uptime   uint64 `json:"uptime"`
	BootTime uint64 `json:"bootTime"`
}

type DeviceInfo struct {
	CpuFreq      uint               `json:"cpuFreq"`
	CpuPercent   []float64          `json:"cpuPercent"`
	CpuTemp      string             `json:"cpuTemp"`
	Disk         DiskInfo           `json:"disk"`
	Host         HostInfo           `json:"host"`
	Mem          MemInfo            `json:"mem"`
	NetInterface []NetInterfaceInfo `json:"netInterface"`
	Wifi         WifiInfo           `json:"wifi"`
}

func GetWifiInfo() (map[string]interface{}, error) {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	obj := conn.Object("net.connman", "/net/connman/technology/wifi")

	call := obj.Call("net.connman.Technology.GetProperties", 0)
	if call.Err != nil {
		return nil, call.Err
	}

	r := call.Body[0].(map[string]dbus.Variant)
	resMap := make(map[string]interface{})
	for k, v := range r {
		a := v.Value()
		ka := util.FisrtLower(k)
		resMap[ka] = a
	}

	return resMap, nil
}

func GetCpuFreq() (interface{}, error) {
	cmd := exec.Command("vcgencmd", "measure_clock", "arm")
	data, err := cmd.Output()
	if err != nil {
		return "", err
	}
	str := strings.Split(strings.ReplaceAll(string(data), "\n", ""), "=")[1]
	return str, nil
}

func GetCpuTemp() (string, error) {
	cmd := exec.Command("vcgencmd", "measure_temp")
	data, err := cmd.Output()
	if err != nil {
		return "", err
	}

	str := strings.Split(strings.ReplaceAll(string(data), "\n", ""), "=")[1]
	return str, nil
}

func GetInfo() map[string]interface{} {

	res := make(map[string]interface{})

	hostInfo, _ := host.Info()
	res["host"] = HostInfo{
		Hostname: hostInfo.Hostname,
		Uptime:   hostInfo.Uptime,
		BootTime: hostInfo.BootTime,
	}

	memInfo, _ := mem.VirtualMemory()
	res["mem"] = MemInfo{
		Total:     memInfo.Total,
		Available: memInfo.Available,
		Free:      memInfo.Free,
		Used:      memInfo.Used,
		Buffers:   memInfo.Buffers,
		Cached:    memInfo.Cached,
		Shared:    memInfo.Shared,
	}

	netInterface, _ := unet.Interfaces()
	res["netInterface"] = netInterface

	diskInfo, _ := disk.Usage("/")
	res["disk"] = []DiskInfo{
		{
			Path:        diskInfo.Path,
			Fstype:      diskInfo.Fstype,
			Total:       diskInfo.Total,
			Free:        diskInfo.Free,
			Used:        diskInfo.Used,
			UsedPercent: diskInfo.UsedPercent,
		},
	}

	cpuPercent, _ := cpu.Percent(2*time.Second, true)
	res["cpuPercent"] = cpuPercent

	return res
}

func GetSn() (string, error) {
	cmd := exec.Command("cat", "/boot/serialno")
	data, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil

}
func GetAddrs(name string) (net.HardwareAddr, []string, error) {
	inet, err := net.InterfaceByName(name)
	if err != nil {
		log.Println("[error] can't get net interface of the machine")
		return nil, nil, err
	}
	addrs, err := inet.Addrs()
	if err != nil {
		log.Println("[error] can't get ip address of the machine")
		return nil, nil, err
	}
	ips := make([]string, len(addrs))
	for i, addr := range addrs {
		ips[i] = addr.String()
	}
	return inet.HardwareAddr, ips, nil
}
