package mdns

import (
	"github.com/pion/mdns"
	"golang.org/x/net/ipv4"
	"net"
)

// DNS 启动内网 mDNS
func DNS() {
	addr, err := net.ResolveUDPAddr("udp", mdns.DefaultAddress)
	if err != nil {
		panic(err)
	}

	l, err := net.ListenUDP("udp4", addr)
	if err != nil {
		panic(err)
	}

	_, err = mdns.Server(ipv4.NewPacketConn(l), &mdns.Config{
		LocalNames: []string{"tessan_sit.local", "tessan_sit"},
	})
	if err != nil {
		panic(err)
	}
}
