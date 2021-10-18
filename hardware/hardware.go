package hardware

import (
	"log"
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

func GetWifiInfo() map[string]interface{} {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	obj := conn.Object("net.connman", "/net/connman/technology/wifi")

	call := obj.Call("net.connman.Technology.GetProperties", 0)
	if call.Err != nil {
		panic(call.Err)
	}

	r := call.Body[0].(map[string]dbus.Variant)
	resMap := make(map[string]interface{})
	for k, v := range r {
		a := v.Value()
		ka := util.FisrtLower(k)
		resMap[ka] = a
	}

	return resMap
}

func GetCpuFreq() (interface{}, error) {
	cmd := exec.Command("vcgencmd", "measure_clock", "arm")
	data, err := cmd.Output()
	if err != nil {
		return "", err
	}
	str := strings.Split(strings.ReplaceAll(string(data), "\n", ""), "=")[1]
	log.Println("cpu current frequence", str)
	return str, nil

}

func GetCpuTemp() (string, error) {
	cmd := exec.Command("vcgencmd", "measure_temp")
	data, err := cmd.Output()
	if err != nil {
		return "", err
	}

	str := strings.Split(strings.ReplaceAll(string(data), "\n", ""), "=")[1]
	log.Println("cpu temp", str)
	return str, nil
}

func GetInfo() map[string]interface{} {

	res := make(map[string]interface{})

	hostInfo, _ := host.Info()
	res["host"] = hostInfo

	memInfo, _ := mem.VirtualMemory()
	res["mem"] = memInfo

	netInterface, _ := unet.Interfaces()
	res["net_interface"] = netInterface

	diskInfo, _ := disk.Usage("/")
	res["disk"] = diskInfo

	cpuPercent, _ := cpu.Percent(1*time.Second, true)
	res["cpu_percent"] = cpuPercent

	cpuInfo, _ := cpu.Info()
	res["cpu_info"] = cpuInfo

	return res
}
