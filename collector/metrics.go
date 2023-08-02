package collector

import "github.com/prometheus/client_golang/prometheus"

type BaseMetrics struct {
	BodyData        []byte
	GcCount         *prometheus.GaugeVec
	GcTime          *prometheus.GaugeVec
	HeapMemoryUsage *prometheus.GaugeVec
}

func (e *BaseMetrics) Describe(ch chan<- *prometheus.Desc) {

}

func BuildBaseMetrics(bodydata []byte, namespace string) BaseMetrics {
	return BaseMetrics{
		BodyData: bodydata,
		HeapMemoryUsage: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "memory",
			Name:      "heap_memory_usage_bytes",
			Help:      "Current heap memory of each mode in bytes",
		}, []string{"mode"}),
		GcCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "jvm_metrics",
			Name:      "gc_count",
			Help:      "GC count of each type",
		}, []string{"type"}),
		GcTime: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "jvm_metrics",
			Name:      "gc_time_milliseconds",
			Help:      "GC time of each type in milliseconds",
		}, []string{"type"}),
	}
}

type OsMetrics struct {
	OpenFileDescriptorCount    prometheus.Gauge
	MaxFileDescriptorCount     prometheus.Gauge
	CommittedVirtualMemorySize prometheus.Gauge
	TotalSwapSpaceSize         prometheus.Gauge
	FreeSwapSpaceSize          prometheus.Gauge
	ProcessCpuTime             prometheus.Gauge
	FreePhysicalMemorySize     prometheus.Gauge
	TotalPhysicalMemorySize    prometheus.Gauge
	SystemCpuLoad              prometheus.Gauge
	ProcessCpuLoad             prometheus.Gauge
	AvailableProcessors        prometheus.Gauge
	SystemLoadAverage          prometheus.Gauge
	OsUnameInfo                *prometheus.GaugeVec
}

func BuildOsMetrics() OsMetrics {
	const namespace = "hadoop"
	return OsMetrics{
		OpenFileDescriptorCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "os",
			Name:      "open_fds_count",
			Help:      "",
		}),
		MaxFileDescriptorCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "os",
			Name:      "max_fds_count",
			Help:      "",
		}),
		CommittedVirtualMemorySize: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "os",
			Name:      "committed_virtual_memory_size_bytes",
			Help:      "",
		}),
		TotalSwapSpaceSize: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "os",
			Name:      "total_swap_space_size_bytes",
			Help:      "",
		}),
		FreeSwapSpaceSize: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "os",
			Name:      "free_swap_space_size_bytes",
			Help:      "",
		}),
		ProcessCpuTime: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "os",
			Name:      "process_cpu_time",
			Help:      "",
		}),
		FreePhysicalMemorySize: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "os",
			Name:      "free_physical_memory_size_bytes",
			Help:      "",
		}),
		TotalPhysicalMemorySize: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "os",
			Name:      "total_physical_memory_size_bytes",
			Help:      "",
		}),
		SystemCpuLoad: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "os",
			Name:      "system_cpu_load",
			Help:      "",
		}),
		ProcessCpuLoad: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "os",
			Name:      "process_cpu_load",
			Help:      "",
		}),
		AvailableProcessors: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "os",
			Name:      "available_processors",
			Help:      "",
		}),
		SystemLoadAverage: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "os",
			Name:      "system_load_average",
			Help:      "",
		}),
		OsUnameInfo: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "os",
			Name:      "info",
			Help:      "",
		}, []string{"arch", "name", "version"}),
	}

}
