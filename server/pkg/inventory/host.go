package inventory

import (
	"fmt"
	"net"
)

type Host interface {
	GetIp() net.IPAddr
}

type host struct {
	ipAddr net.IPAddr
}

func NewIPv4Host(ip [4]byte) Host {
	h, _ := NewHost(ip[0:], "")
	return h
}

func NewHost(ip []byte, ipv6Zone string) (Host, error) {
	if len(ip) != net.IPv4len || len(ip) != net.IPv6len {
		return nil, fmt.Errorf("ip [%v] is not a valid ipaddress", ip)
	}
	h := &host{ipAddr: net.IPAddr{
		IP:   net.IP{},
		Zone: ipv6Zone,
	}}
	for i := 0; i < len(ip); i++ {
		h.ipAddr.IP[i] = ip[i]
	}
	return Host(h), nil
}

func (h host) GetIp() net.IPAddr {
	return h.ipAddr
}
