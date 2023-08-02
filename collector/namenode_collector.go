package collector

import (
	"encoding/json"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
)

type NameNodeMetrics struct {
	BaseMetrics
	OsMetrics
	MissingBlocks         prometheus.Gauge
	UnderReplicatedBlocks prometheus.Gauge
	Capacity              *prometheus.GaugeVec
	BlocksTotal           prometheus.Gauge
	FilesTotal            prometheus.Gauge
	CorruptBlocks         prometheus.Gauge
	ExcessBlocks          prometheus.Gauge
	StaleDataNodes        prometheus.Gauge
	LastHATransitionTime  prometheus.Gauge
	HAState               prometheus.Gauge
	RpcReceivedBytes      *prometheus.GaugeVec
	RpcSentBytes          *prometheus.GaugeVec
	RpcQueueTimeNumOps    *prometheus.GaugeVec // RpcProcessingTimeNumOps = RpcQueueTimeNumOps
	RpcAvgTime            *prometheus.GaugeVec
	RpcNumOpenConnections *prometheus.GaugeVec // current number of open connections
	RpcCallQueueLength    *prometheus.GaugeVec
}

func NewNameNodeMetrics(t Target) *NameNodeMetrics {

	const namespace = "hdfs_namenode"

	return &NameNodeMetrics{
		BaseMetrics: BuildBaseMetrics(t.BodyData, namespace),
		OsMetrics:   BuildOsMetrics(),
		MissingBlocks: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "fsname_system",
			Name:      "missing_blocks",
			Help:      "Current number of missing blocks",
		}),
		UnderReplicatedBlocks: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "fsname_system",
			Name:      "under_replicated_blocks",
			Help:      "Current number of blocks under replicated",
		}),
		Capacity: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "fsname_system",
			Name:      "capacity_bytes",
			Help:      "Current DataNodes capacity in each mode in bytes",
		}, []string{"mode"}),
		BlocksTotal: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "fsname_system",
			Name:      "blocks_total",
			Help:      "Current number of allocated blocks in the system",
		}),
		FilesTotal: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "fsname_system",
			Name:      "files_total",
			Help:      "Current number of files and directories",
		}),
		CorruptBlocks: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "fsname_system",
			Name:      "corrupt_blocks",
			Help:      "Current number of blocks with corrupt replicas",
		}),
		ExcessBlocks: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "fsname_system",
			Name:      "excess_blocks",
			Help:      "Current number of excess blocks",
		}),
		StaleDataNodes: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "fsname_system",
			Name:      "stale_datanodes",
			Help:      "Current number of DataNodes marked stale due to delayed heartbeat",
		}),
		LastHATransitionTime: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "namenode_status",
			Name:      "last_ha_transition_time",
			Help:      "last HA Transition Time",
		}),
		HAState: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "fsname_system",
			Name:      "hastate",
			Help:      "Current state of the NameNode: 0.0 (for initializing) or 1.0 (for active) or 2.0 (for standby) or 3.0 (for stopping) state",
		}),
		RpcReceivedBytes: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "rpc_activity",
			Name:      "received_bytes",
			Help:      "Total number of received bytes",
		}, []string{"port"}),
		RpcSentBytes: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "rpc_activity",
			Name:      "sent_bytes",
			Help:      "Total number of sent bytes",
		}, []string{"port"}),
		RpcQueueTimeNumOps: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "rpc_activity",
			Name:      "call_count",
			Help:      "Total number of RPC calls (same to RpcQueueTimeNumOps) ",
		}, []string{"port", "method"}),
		RpcAvgTime: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "rpc_activity",
			Name:      "avg_time_milliseconds",
			Help:      "current number of open connections",
		}, []string{"port", "method"}),
		RpcNumOpenConnections: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "rpc_activity",
			Name:      "open_connections_count",
			Help:      "current number of open connections",
		}, []string{"port"}),
		RpcCallQueueLength: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "rpc_activity",
			Name:      "call_queue_length",
			Help:      "Current length of the call queue",
		}, []string{"port"}),
	}
}

// Collect implements the prometheus.Collector interface.
func (e *NameNodeMetrics) Collect(ch chan<- prometheus.Metric) {

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
		if DataMap["name"] == "Hadoop:service=NameNode,name=FSNamesystem" {
			e.MissingBlocks.Set(DataMap["MissingBlocks"].(float64))
			e.UnderReplicatedBlocks.Set(DataMap["UnderReplicatedBlocks"].(float64))
			e.Capacity.WithLabelValues("Total").Set(DataMap["CapacityTotal"].(float64))
			e.Capacity.WithLabelValues("Used").Set(DataMap["CapacityUsed"].(float64))
			e.Capacity.WithLabelValues("Remaining").Set(DataMap["CapacityRemaining"].(float64))
			e.Capacity.WithLabelValues("UsedNonDFS").Set(DataMap["CapacityUsedNonDFS"].(float64))
			e.BlocksTotal.Set(DataMap["BlocksTotal"].(float64))
			e.FilesTotal.Set(DataMap["FilesTotal"].(float64))
			e.CorruptBlocks.Set(DataMap["CorruptBlocks"].(float64))
			e.ExcessBlocks.Set(DataMap["ExcessBlocks"].(float64))
			e.StaleDataNodes.Set(DataMap["StaleDataNodes"].(float64))

			switch DataMap["tag.HAState"] {

			case "initializing":
				e.HAState.Set(0)
			case "active":
				e.HAState.Set(1)
			case "standby":
				e.HAState.Set(2)
			case "stopping":
				e.HAState.Set(3)

			}
		}

		if DataMap["name"] == "Hadoop:service=NameNode,name=NameNodeStatus" {

			e.LastHATransitionTime.Set(DataMap["LastHATransitionTime"].(float64))
		}

		if DataMap["name"] == "Hadoop:service=NameNode,name=JvmMetrics" {
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
		}

		if strings.HasPrefix(DataMap["modelerType"].(string), "RpcActivityForPort") {

			port := DataMap["tag.port"].(string)

			e.RpcReceivedBytes.WithLabelValues(port).Set(DataMap["ReceivedBytes"].(float64))
			e.RpcSentBytes.WithLabelValues(port).Set(DataMap["SentBytes"].(float64))
			e.RpcQueueTimeNumOps.WithLabelValues(port, "QueueTime").Set(DataMap["RpcQueueTimeNumOps"].(float64))
			e.RpcAvgTime.WithLabelValues(port, "RpcQueueTime").Set(DataMap["RpcQueueTimeAvgTime"].(float64))
			e.RpcAvgTime.WithLabelValues(port, "RpcProcessingTime").Set(DataMap["RpcProcessingTimeAvgTime"].(float64))
			e.RpcNumOpenConnections.WithLabelValues(port).Set(DataMap["NumOpenConnections"].(float64))
			e.RpcCallQueueLength.WithLabelValues(port).Set(DataMap["CallQueueLength"].(float64))
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

	e.MissingBlocks.Collect(ch)
	e.UnderReplicatedBlocks.Collect(ch)
	e.Capacity.Collect(ch)
	e.BlocksTotal.Collect(ch)
	e.FilesTotal.Collect(ch)
	e.CorruptBlocks.Collect(ch)
	e.ExcessBlocks.Collect(ch)
	e.StaleDataNodes.Collect(ch)
	e.GcCount.Collect(ch)
	e.GcTime.Collect(ch)
	e.HeapMemoryUsage.Collect(ch)
	e.LastHATransitionTime.Collect(ch)
	e.HAState.Collect(ch)
	e.RpcReceivedBytes.Collect(ch)
	e.RpcSentBytes.Collect(ch)
	e.RpcQueueTimeNumOps.Collect(ch)
	e.RpcAvgTime.Collect(ch)
	e.RpcNumOpenConnections.Collect(ch)
	e.RpcCallQueueLength.Collect(ch)
}

func NameNodeCollector(target Target, registry *prometheus.Registry) (success bool) {

	metrics := NewNameNodeMetrics(target)
	registry.MustRegister(metrics)

	return true
}
