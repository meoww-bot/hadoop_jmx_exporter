package collector

import (
	"encoding/json"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
)

type HbaseMasterMetrics struct {
	BaseMetrics
	OsMetrics
}

func NewHbaseMasterMetrics(t Target) *HbaseMasterMetrics {

	const namespace = "hbase_master"
	return &HbaseMasterMetrics{
		BaseMetrics: BuildBaseMetrics(t.BodyData, namespace),
		OsMetrics:   BuildOsMetrics(),
	}
}

// Collect implements the prometheus.Collector interface.
func (e *HbaseMasterMetrics) Collect(ch chan<- prometheus.Metric) {
	var err error

	var f interface{}
	err = json.Unmarshal(e.BodyData, &f)
	if err != nil {
		log.Error(err)
	}
	m := f.(map[string]interface{})
	var List = m["beans"].([]interface{})
	for _, Data := range List {
		DataMap := Data.(map[string]interface{})

		if DataMap["name"] == "java.lang:type=GarbageCollector,name=ParNew" {
			e.GcTime.WithLabelValues("ParNew").Set(DataMap["CollectionTime"].(float64))
			e.GcCount.WithLabelValues("ParNew").Set(DataMap["CollectionCount"].(float64))
		}
		if DataMap["name"] == "java.lang:type=GarbageCollector,name=ConcurrentMarkSweep" {
			e.GcTime.WithLabelValues("ConcurrentMarkSweep").Set(DataMap["CollectionTime"].(float64))
			e.GcCount.WithLabelValues("ConcurrentMarkSweep").Set(DataMap["CollectionCount"].(float64))

		}
		if DataMap["name"] == "java.lang:type=Memory" {
			heapMemoryUsage := DataMap["HeapMemoryUsage"].(map[string]interface{})
			e.HeapMemoryUsage.WithLabelValues("committed").Set(heapMemoryUsage["committed"].(float64))
			e.HeapMemoryUsage.WithLabelValues("init").Set(heapMemoryUsage["init"].(float64))
			e.HeapMemoryUsage.WithLabelValues("max").Set(heapMemoryUsage["max"].(float64))
			e.HeapMemoryUsage.WithLabelValues("used").Set(heapMemoryUsage["used"].(float64))
			e.HeapMemoryUsage.Collect(ch)
		}

		if DataMap["name"] == "java.lang:type=OperatingSystem" {
			e.OsMetrics.MaxFileDescriptorCount.Set(DataMap["MaxFileDescriptorCount"].(float64))
			e.OsMetrics.OpenFileDescriptorCount.Set(DataMap["OpenFileDescriptorCount"].(float64))
			e.OsMetrics.CommittedVirtualMemorySize.Set(DataMap["CommittedVirtualMemorySize"].(float64))
			e.OsMetrics.TotalSwapSpaceSize.Set(DataMap["TotalSwapSpaceSize"].(float64))
			e.OsMetrics.FreeSwapSpaceSize.Set(DataMap["FreeSwapSpaceSize"].(float64))
			e.OsMetrics.ProcessCpuTime.Set(DataMap["ProcessCpuTime"].(float64))
			e.OsMetrics.TotalPhysicalMemorySize.Set(DataMap["TotalPhysicalMemorySize"].(float64))
			e.OsMetrics.SystemCpuLoad.Set(DataMap["SystemCpuLoad"].(float64))
			e.OsMetrics.ProcessCpuLoad.Set(DataMap["ProcessCpuLoad"].(float64))
			e.OsMetrics.FreePhysicalMemorySize.Set(DataMap["FreePhysicalMemorySize"].(float64))
			e.OsMetrics.AvailableProcessors.Set(DataMap["AvailableProcessors"].(float64))
			e.OsMetrics.SystemLoadAverage.Set(DataMap["SystemLoadAverage"].(float64))

			e.OsMetrics.OsUnameInfo.With(
				prometheus.Labels{
					"arch":    DataMap["Arch"].(string),
					"name":    DataMap["Name"].(string),
					"version": DataMap["Version"].(string),
				}).Set(1)

			e.OsMetrics.MaxFileDescriptorCount.Collect(ch)
			e.OsMetrics.OpenFileDescriptorCount.Collect(ch)
			e.OsMetrics.CommittedVirtualMemorySize.Collect(ch)
			e.OsMetrics.TotalSwapSpaceSize.Collect(ch)
			e.OsMetrics.FreeSwapSpaceSize.Collect(ch)
			e.OsMetrics.ProcessCpuTime.Collect(ch)
			e.OsMetrics.TotalPhysicalMemorySize.Collect(ch)
			e.OsMetrics.SystemCpuLoad.Collect(ch)
			e.OsMetrics.ProcessCpuLoad.Collect(ch)
			e.OsMetrics.FreePhysicalMemorySize.Collect(ch)
			e.OsMetrics.AvailableProcessors.Collect(ch)
			e.OsMetrics.SystemLoadAverage.Collect(ch)

			e.OsMetrics.OsUnameInfo.Collect(ch)

		}
	}
	e.GcCount.Collect(ch)
	e.GcTime.Collect(ch)

}

func HbaseMasterCollector(target Target, registry *prometheus.Registry) (success bool) {

	metrics := NewHbaseMasterMetrics(target)
	registry.MustRegister(metrics)

	return true
}
