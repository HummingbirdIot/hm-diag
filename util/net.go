package util

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"strings"

	"github.com/kpango/glg"
)

var (
	getGatewayCmd = "route -n | awk '{if ($4 == \"UG\") {print $2}}'"
	GatewayAddr   string
)

func init() {
	cmd := exec.Command("bash", "-c", getGatewayCmd)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		glg.Error("get dns addr error : ", err.Error(), stderr.String())
	}

	outString := out.String()
	s := strings.ReplaceAll(outString, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	glg.Info("get dns addr:", s)
	GatewayAddr = s
}

func IsPrivateIp(ip net.IP) bool {
	return ip.IsPrivate() || ip.IsLoopback()
}

func IpByInterfaceName(name string) (string, error) {
	iface, err := net.InterfaceByName(name)
	glg.Infof("========== ip %#v", iface)
	if err != nil {
		return "", err
	}
	if iface.Flags&net.FlagUp == 0 {
		return "", errors.New("no ip")
	}
	if iface.Flags&net.FlagLoopback != 0 {
		return "", errors.New("no ip")
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if ip == nil || ip.IsLoopback() {
			continue
		}
		ip = ip.To4()
		if ip == nil {
			continue // not an ipv4 address
		}
		return ip.String(), nil
	}
	return "", errors.New("are you connected to the network?")
}

func PingTest(ip string) error {
	cmd := exec.Command("ping", ip, "-c", "1")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &out
	err := cmd.Run()
	glg.Debug("cmd out: ", out.String())
	if err != nil {
		glg.Error("network error : ", err.Error(), stderr.String())
		return fmt.Errorf(stderr.String())
	}
	glg.Info("network ok")
	return nil
}
