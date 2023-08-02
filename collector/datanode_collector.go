package collector

import (
	"encoding/json"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
)

type DataNodeMetrics struct {
	BaseMetrics
	Capacity              *prometheus.GaugeVec
	CacheCapacity         prometheus.Gauge
	CacheUsed             prometheus.Gauge
	FailedVolumes         prometheus.Gauge
	EstimatedCapacityLost prometheus.Gauge
	BlocksCached          prometheus.Gauge
	BlocksFailedToCache   prometheus.Gauge
	BlocksFailedToUncache prometheus.Gauge
}

func NewDataNodeMetrics(t Target) *DataNodeMetrics {

	const namespace = "hdfs_datanode"

	return &DataNodeMetrics{
		BaseMetrics: BuildBaseMetrics(t.BodyData, namespace),
		Capacity: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "fsname_system",
			Name:      "capacity_bytes",
			Help:      "Current DataNodes capacity in each mode in bytes",
		}, []string{"mode"}),
		CacheCapacity: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "CacheCapacity",
			Help:      "CacheCapacity",
		}),
		CacheUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "CacheUsed",
			Help:      "CacheUsed",
		}),

		FailedVolumes: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "FailedVolumes",
			Help:      "FailedVolumes",
		}),
		EstimatedCapacityLost: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "EstimatedCapacityLost",
			Help:      "EstimatedCapacityLost",
		}),

		BlocksCached: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "BlocksCached",
			Help:      "BlocksCached",
		}),
		BlocksFailedToCache: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "BlocksFailedToCache",
			Help:      "BlocksFailedToCache",
		}),
		BlocksFailedToUncache: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "BlocksFailedToUncache",
			Help:      "BlocksFailedToUncache",
		}),
	}
}

// Collect implements the prometheus.Collector interface.
func (e *DataNodeMetrics) Collect(ch chan<- prometheus.Metric) {
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
		if DataMap["name"] == "Hadoop:service=DataNode,name=FSDatasetState" {

			e.Capacity.WithLabelValues("Total").Set(DataMap["Capacity"].(float64))
			e.Capacity.WithLabelValues("DfsUsed").Set(DataMap["DfsUsed"].(float64))
			e.Capacity.WithLabelValues("Remaining").Set(DataMap["Remaining"].(float64))

			e.CacheCapacity.Set(DataMap["CacheCapacity"].(float64))
			e.CacheUsed.Set(DataMap["CacheUsed"].(float64))

			e.FailedVolumes.Set(DataMap["NumFailedVolumes"].(float64))
			e.EstimatedCapacityLost.Set(DataMap["EstimatedCapacityLostTotal"].(float64))

			e.BlocksCached.Set(DataMap["NumBlocksCached"].(float64))
			e.BlocksFailedToCache.Set(DataMap["NumBlocksFailedToCache"].(float64))
			e.BlocksFailedToUncache.Set(DataMap["NumBlocksFailedToUncache"].(float64))
		}

		if DataMap["name"] == "Hadoop:service=DataNode,name=JvmMetrics" {
			e.GcCount.WithLabelValues("ParNew").Set(DataMap["GcCountParNew"].(float64))
			e.GcCount.WithLabelValues("ConcurrentMarkSweep").Set(DataMap["GcCountConcurrentMarkSweep"].(float64))

			e.GcTime.WithLabelValues("ParNew").Set(DataMap["GcTimeMillisParNew"].(float64))
			e.GcTime.WithLabelValues("ConcurrentMarkSweep").Set(DataMap["GcTimeMillisConcurrentMarkSweep"].(float64))

		}
		if DataMap["name"] == "java.lang:type=Memory" {
			heapMemoryUsage := DataMap["HeapMemoryUsage"].(map[string]interface{})
			e.HeapMemoryUsage.WithLabelValues("committed").Set(heapMemoryUsage["committed"].(float64))
			e.HeapMemoryUsage.WithLabelValues("init").Set(heapMemoryUsage["init"].(float64))
			e.HeapMemoryUsage.WithLabelValues("max").Set(heapMemoryUsage["max"].(float64))
			e.HeapMemoryUsage.WithLabelValues("used").Set(heapMemoryUsage["used"].(float64))
			e.HeapMemoryUsage.Collect(ch)

		}
	}
	e.Capacity.Collect(ch)
	e.CacheCapacity.Collect(ch)
	e.CacheUsed.Collect(ch)
	e.FailedVolumes.Collect(ch)
	e.EstimatedCapacityLost.Collect(ch)
	e.BlocksCached.Collect(ch)
	e.BlocksFailedToCache.Collect(ch)
	e.BlocksFailedToUncache.Collect(ch)

}

func DataNodeCollector(target Target, registry *prometheus.Registry) (success bool) {

	metrics := NewDataNodeMetrics(target)
	registry.MustRegister(metrics)
	return
}
