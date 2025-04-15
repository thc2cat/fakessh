package main

import (
	"log"
	"net"
)

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
