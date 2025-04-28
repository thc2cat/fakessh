package main

import (
	"log"
	"net"
)

// var (
// 	local = []string{ // + your networks
// 		"10.0.0.0/8",
// 		"172.16.0.0/12",
// 		"192.168.0.0/16",
// 	}
// )

func ipinrange(ip string) bool {
	// Check if the IP is in the local range
	for _, cidr := range local {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			log.Printf("Error parsing CIDR %s: %v", cidr, err)
			continue
		}
		if network.Contains(net.ParseIP(ip)) {
			return true
		}
	}
	return false
}
