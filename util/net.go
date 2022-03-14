package util

import "net"

func IsPrivateIp(ip net.IP) bool {
	return ip.IsPrivate() || ip.IsLoopback()
}
