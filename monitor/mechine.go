package monitor

import (
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	nett "github.com/shirou/gopsutil/v3/net"
	"net"
	"time"
)

func GetAllInfo() (ip net.IP, c []float64, mp, dp float64, m, mf, d, df, s, r uint64) {
	ip, _ = externalIP()
	c, _ = cpu.Percent(time.Second, false)
	memo, _ := mem.VirtualMemory()
	mp = memo.UsedPercent
	m = memo.Used
	mf = memo.Available
	usage, _ := disk.Usage("/")
	dp = usage.UsedPercent
	d = usage.Used
	df = usage.Free
	counters, err := nett.IOCounters(false)
	if err == nil {
		s, r = counters[0].BytesSent, counters[0].BytesRecv
	}
	return
}

// 获取ip
func externalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

// 获取ip
func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}
