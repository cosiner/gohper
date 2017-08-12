package net2

import "net"

var (
	IPv4 = []string{"192.168.0.0/16", "172.16.0.0/12", "10.0.0.0/8"}
)

func ParseCIDRs(networks ...string) ([]*net.IPNet, error) {
	nets := make([]*net.IPNet, 0, len(networks))
	for _, network := range networks {
		_, ipnet, err := net.ParseCIDR(network)
		if err != nil {
			return nil, err
		}
		nets = append(nets, ipnet)
	}
	return nets, nil
}

func IsIPContained(ip net.IP, networks []*net.IPNet) bool {
	for _, network := range networks {
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

const LOCALHOST = "127.0.0.1"

func Localhost(networks ...string) string {
	if len(networks) == 0 {
		networks = IPv4
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return LOCALHOST
	}
	ipnets, err := ParseCIDRs(networks...)
	if err != nil {
		return LOCALHOST
	}
	for _, addr := range addrs {
		if ip, ok := addr.(*net.IPNet); ok {
			if IsIPContained(ip.IP, ipnets) {
				return ip.IP.String()
			}
		}
	}
	return LOCALHOST
}

func ReplaceHost(addr, host string) (string, error) {
	_, port, err := net.SplitHostPort(addr)
	return host + ":" + port, err
}
