package collector

import (
	"encoding/json"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
)

type JournalNodeMetrics struct {
	BaseMetrics
}

func NewJournalNodeMetrics(t Target) *JournalNodeMetrics {

	const namespace = "hdfs_journalnode"
	return &JournalNodeMetrics{
		BaseMetrics: BuildBaseMetrics(t.BodyData, namespace),
	}
}

// Collect implements the prometheus.Collector interface.
func (e *JournalNodeMetrics) Collect(ch chan<- prometheus.Metric) {
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
		/*
			"name" : "java.lang:type=Memory",
			"modelerType" : "sun.management.MemoryImpl",
			"HeapMemoryUsage" : {
				"committed" : 1060372480,
				"init" : 1073741824,
				"max" : 1060372480,
				"used" : 124571464
			},
		*/
		if DataMap["name"] == "java.lang:type=Memory" {
			heapMemoryUsage := DataMap["HeapMemoryUsage"].(map[string]interface{})
			e.HeapMemoryUsage.WithLabelValues("committed").Set(heapMemoryUsage["committed"].(float64))
			e.HeapMemoryUsage.WithLabelValues("init").Set(heapMemoryUsage["init"].(float64))
			e.HeapMemoryUsage.WithLabelValues("max").Set(heapMemoryUsage["max"].(float64))
			e.HeapMemoryUsage.WithLabelValues("used").Set(heapMemoryUsage["used"].(float64))
		}
	}
	e.GcCount.Collect(ch)
	e.GcTime.Collect(ch)
	e.HeapMemoryUsage.Collect(ch)
}

func JournalNodeCollector(target Target, registry *prometheus.Registry) (success bool) {

	metrics := NewJournalNodeMetrics(target)
	registry.MustRegister(metrics)

	return true
}
