package funcs

import (
	"strings"

	"github.com/GitHamburg/windows-agent/g"
	"github.com/open-falcon/common/model"
	"github.com/shirou/gopsutil/net"
	"github.com/Lofanmi/pinyin-golang/pinyin"
)

func net_status(ifacePrefix []string) ([]net.IOCountersStat, error) {
	net_iocounter, err := net.IOCounters(true)
	netIfs := []net.IOCountersStat{}
	for _, iface := range ifacePrefix {
		for _, netIf := range net_iocounter {
			if strings.Contains(netIf.Name, iface) {
				netIfs = append(netIfs, netIf)
			}
		}
	}
	return netIfs, err
}

func NetMetrics() []*model.MetricValue {
	return CoreNetMetrics(g.Config().Collector.IfacePrefix)
}

func CoreNetMetrics(ifacePrefix []string) (L []*model.MetricValue) {

	netIfs, err := net_status(ifacePrefix)
	if err != nil {
		g.Logger().Println(err)
		return []*model.MetricValue{}
	}

	for _, netIf := range netIfs {
		netIfName := pinyin.NewDict().Convert(netIf.Name, "_").None()
		iface := "iface=" + netIfName
		L = append(L, CounterValue("net.if.in.bytes", netIf.BytesRecv, iface)) //此处乘以8即为bit的流量
		L = append(L, CounterValue("net.if.in.packets", netIf.PacketsRecv, iface))
		L = append(L, CounterValue("net.if.in.errors", netIf.Errin, iface))
		L = append(L, CounterValue("net.if.in.dropped", netIf.Dropin, iface))
		L = append(L, CounterValue("net.if.in.fifo.errs", netIf.Fifoin, iface))
		L = append(L, CounterValue("net.if.out.bytes", netIf.BytesSent, iface)) //此处乘以8即为bit的流量
		L = append(L, CounterValue("net.if.out.packets", netIf.PacketsSent, iface))
		L = append(L, CounterValue("net.if.out.errors", netIf.Errout, iface))
		L = append(L, CounterValue("net.if.out.dropped", netIf.Dropout, iface))
		L = append(L, CounterValue("net.if.out.fifo.errs", netIf.Fifoout, iface))
	}
	return
}
