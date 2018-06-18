package net

import (
	"net"
)

var privateIPv4NetsString = []string{
	"10.0.0.0/8",     // RFC1918
	"172.16.0.0/12",  // RFC1918
	"192.168.0.0/16", // RFC1918
	// "fe80::/10",      // IPv6 link-local
}

var privateIPv4Nets []*net.IPNet
var privateIPv6Net *net.IPNet

func init() {
	for _, IPNetString := range privateIPv4NetsString {
		_, ipNet, _ := net.ParseCIDR(IPNetString)
		privateIPv4Nets = append(privateIPv4Nets, ipNet)
	}

	_, privateIPv6Net, _ = net.ParseCIDR("fe80::/10")
}

func IsPublicIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
		return false
	}

	if v4 := ip.To4(); v4 != nil {
		for _, block := range privateIPv4Nets {
			if block.Contains(ip) {
				return false
			}
		}
		return true
	}

	if privateIPv6Net.Contains(ip) {
		return false
	}

	return true
}
