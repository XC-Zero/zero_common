package main

import (
	_ "embed"
	"github.com/davecgh/go-spew/spew"
	"github.com/mitchellh/mapstructure"
	"github.com/pion/mdns"
	"github.com/spf13/viper"
	"golang.org/x/net/ipv4"

	"net"
	"strings"
)

type mDNSConf struct {
	LocalNames       []string `toml:"local_names"`
	BroadcastAddress string   `toml:"broadcast_address"`
}

var globalViper = viper.New()
var conf mDNSConf

func init() {
	globalViper.SetConfigName("conf")
	globalViper.AutomaticEnv()
	globalViper.AddConfigPath("./conf")
	globalViper.SetConfigType("toml")
	if err := globalViper.ReadInConfig(); err != nil {
		panic(err)
	}

	err := globalViper.Unmarshal(&conf, setTagName)
	if err != nil {
		panic(err)
	}
}

// DNS 启动内网 mDNS
func DNS() {

	spew.Dump(conf)
	if strings.TrimSpace(conf.BroadcastAddress) == "" {
		conf.BroadcastAddress = mdns.DefaultAddress
	}
	spew.Dump("Broad cast address is " + conf.BroadcastAddress)
	addr, err := net.ResolveUDPAddr("udp", conf.BroadcastAddress)
	if err != nil {
		panic(err)
	}
	spew.Dump("Listen UDP4 address is " + addr.String())

	l, err := net.ListenUDP("udp4", addr)
	if err != nil {
		panic(err)
	}
	spew.Dump("Server domain is [" + strings.Join(conf.LocalNames, ",") + "]")
	_, err = mdns.Server(ipv4.NewPacketConn(l), &mdns.Config{
		LocalNames: conf.LocalNames,
	})
	if err != nil {
		panic(err)
	}
	select {}

}

// 设置 config 对应的结构体的 tag
func setTagName(d *mapstructure.DecoderConfig) {
	d.TagName = "toml"
}

func main() {
	DNS()
}
