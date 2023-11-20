package mdns

import (
	_ "embed"
	"encoding/json"
	"github.com/pion/mdns"
	"golang.org/x/net/ipv4"
	"net"
)

type mDNSConf struct {
	LocalNames       []string `json:"local_names"`
	BroadCastAddress string   `json:"broad_cast_address"`
}

//go:embed  conf.json
var file []byte

// DNS 启动内网 mDNS
func DNS() {
	var conf mDNSConf
	err := json.Unmarshal(file, &conf)
	if err != nil {
		panic(err)
	}
	if conf.BroadCastAddress == "" {
		conf.BroadCastAddress = mdns.DefaultAddress
	}
	addr, err := net.ResolveUDPAddr("udp", conf.BroadCastAddress)
	if err != nil {
		panic(err)
	}

	l, err := net.ListenUDP("udp4", addr)
	if err != nil {
		panic(err)
	}

	_, err = mdns.Server(ipv4.NewPacketConn(l), &mdns.Config{
		LocalNames: conf.LocalNames,
	})
	if err != nil {
		panic(err)
	}
}
