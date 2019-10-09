package funcs

import (
	"github.com/GitHamburg/windows-agent/tools/disk"
	"github.com/open-falcon/common/model"
	"log"
	"github.com/GitHamburg/windows-agent/g"
)


func DiskIOPhysicalMetrics() (L []*model.MetricValue) {

	disk_iocounter, err := IOCounters()
	if err != nil {
		g.Logger().Println(err)
		return
	}

	for device, ds := range disk_iocounter {

		device := "device=" + device
		L = append(L, CounterValue("disk.io.msec_read", ds.Msec_Read, device))
		L = append(L, CounterValue("disk.io.msec_write", ds.Msec_Write, device))
		L = append(L, CounterValue("disk.io.read_bytes", ds.Read_Bytes, device))
		L = append(L, CounterValue("disk.io.read_requests", ds.Read_Requests, device))
		L = append(L, CounterValue("disk.io.write_bytes", ds.Write_Bytes, device))
		L = append(L, CounterValue("disk.io.write_requests", ds.Write_Requests, device))
		L = append(L, GaugeValue("disk.io.util", 100-ds.Util, device))
	}
	return
}

func DiskIOMetrics() (L []*model.MetricValue) {

	dsList, err := disk.DiskIOCounters()
	if err != nil {
		log.Println("Get devices io fail: ", err)
		return
	}

	for _, ds := range dsList {
		device := "device=" + ds.Name

		L = append(L, CounterValue("disk.io.read_requests", ds.ReadCount, device))
		L = append(L, CounterValue("disk.io.read_bytes", ds.ReadBytes, device))
		L = append(L, CounterValue("disk.io.write_requests", ds.WriteCount, device))
		L = append(L, CounterValue("disk.io.write_bytes", ds.WriteBytes, device))
		L = append(L, CounterValue("disk.io.read_time", ds.ReadTime, device))
		L = append(L, CounterValue("disk.io.write_time", ds.WriteTime, device))
		L = append(L, CounterValue("disk.io.iotime", ds.IoTime, device))
	}
	return
}