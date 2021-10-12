package hardware

import (
	"log"
	"os/exec"
	"strings"

	"github.com/godbus/dbus/v5"
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
		resMap[k] = a
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
