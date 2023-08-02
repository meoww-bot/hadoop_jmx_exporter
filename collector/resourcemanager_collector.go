package collector

import (
	"encoding/json"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
)

// https://hadoop.apache.org/docs/stable/hadoop-project-dist/hadoop-common/Metrics.html#ClusterMetrics
type ClusterMetrics struct {
	NodeManagerNums        *prometheus.GaugeVec
	AMLaunchDelayNumOps    prometheus.Gauge
	AMLaunchDelayAvgTime   prometheus.Gauge
	AMRegisterDelayNumOps  prometheus.Gauge
	AMRegisterDelayAvgTime prometheus.Gauge
}

// https://hadoop.apache.org/docs/stable/hadoop-project-dist/hadoop-common/Metrics.html#QueueMetrics
type QueueMetrics struct {
	Running_0    *prometheus.GaugeVec
	Running_60   *prometheus.GaugeVec
	Running_300  *prometheus.GaugeVec
	Running_1440 *prometheus.GaugeVec
	// AMResourceLimitMB                              *prometheus.GaugeVec
	// AMResourceLimitVCores                          *prometheus.GaugeVec
	// UsedAMResourceMB                               *prometheus.GaugeVec
	// UsedAMResourceVCores                           *prometheus.GaugeVec
	// UsedCapacity                                   *prometheus.GaugeVec
	// AbsoluteUsedCapacity                           *prometheus.GaugeVec
	AppsCount                                      *prometheus.GaugeVec
	AggregateNodeLocalContainersAllocated          *prometheus.GaugeVec
	AggregateRackLocalContainersAllocated          *prometheus.GaugeVec
	AggregateOffSwitchContainersAllocated          *prometheus.GaugeVec
	AggregateContainersPreempted                   *prometheus.GaugeVec
	AggregateMemoryMBSecondsPreempted              *prometheus.GaugeVec
	AggregateVcoreSecondsPreempted                 *prometheus.GaugeVec
	ActiveUsers                                    *prometheus.GaugeVec
	ActiveApplications                             *prometheus.GaugeVec
	AppAttemptFirstContainerAllocationDelayNumOps  *prometheus.GaugeVec
	AppAttemptFirstContainerAllocationDelayAvgTime *prometheus.GaugeVec
	AllocatedMB                                    *prometheus.GaugeVec
	AllocatedVCores                                *prometheus.GaugeVec
	AllocatedContainers                            *prometheus.GaugeVec
	AggregateContainersAllocated                   *prometheus.GaugeVec
	AggregateContainersReleased                    *prometheus.GaugeVec
	AvailableMB                                    *prometheus.GaugeVec
	AvailableVCores                                *prometheus.GaugeVec
	PendingMB                                      *prometheus.GaugeVec
	PendingVCores                                  *prometheus.GaugeVec
	PendingContainers                              *prometheus.GaugeVec
	ReservedMB                                     *prometheus.GaugeVec
	ReservedVCores                                 *prometheus.GaugeVec
	ReservedContainers                             *prometheus.GaugeVec
}

type ResourceManagerMetrics struct {
	BaseMetrics
	ClusterMetrics
	QueueMetrics
}

func NewResourceManagerMetrics(t Target) *ResourceManagerMetrics {

	const namespace = "yarn_resourcemanager"

	return &ResourceManagerMetrics{
		BaseMetrics: BuildBaseMetrics(t.BodyData, namespace),
		ClusterMetrics: ClusterMetrics{
			NodeManagerNums: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "cluster_metrics",
				Name:      "nodemanager_nums",
				Help:      "Current NodeManagers numbers of each state",
			}, []string{"state"}),
			AMLaunchDelayNumOps: prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "cluster_metrics",
				Name:      "am_launch_delay_num_ops",
				Help:      "Total number of AMs launched",
			}),
			AMLaunchDelayAvgTime: prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "cluster_metrics",
				Name:      "am_launch_delay_avg_time_milliseconds",
				Help:      "Average time in milliseconds RM spends to launch AM containers after the AM container is allocated",
			}),
			AMRegisterDelayNumOps: prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "cluster_metrics",
				Name:      "am_register_delay_num_ops",
				Help:      "Total number of AMs registered",
			}),
			AMRegisterDelayAvgTime: prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "cluster_metrics",
				Name:      "am_register_delay_avg_time_milliseconds",
				Help:      "Average time in milliseconds AM spends to register with RM after the AM container gets launched",
			}),
		},
		QueueMetrics: QueueMetrics{
			Running_0: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "running_0",
				Help:      "Current number of running applications whose elapsed time are less than 60 minutes",
			}, []string{"queue"}),
			Running_60: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "running_60",
				Help:      "Current number of running applications whose elapsed time are between 60 and 300 minutes",
			}, []string{"queue"}),
			Running_300: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "running_300",
				Help:      "Current number of running applications whose elapsed time are between 300 and 1440 minutes",
			}, []string{"queue"}),
			Running_1440: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "Running_1440",
				Help:      "Current number of running applications elapsed time are more than 1440 minutes",
			}, []string{"queue"}),
			// AMResourceLimitMB: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			// 	Namespace: namespace,
			// 	Subsystem: "queue_metrics",
			// 	Name:      "nodemanager_nums",
			// 	Help:      "Current NodeManagers numbers of each state",
			// }, []string{"queue"}),
			// AMResourceLimitVCores: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			// 	Namespace: namespace,
			// 	Subsystem: "queue_metrics",
			// 	Name:      "nodemanager_nums",
			// 	Help:      "Current NodeManagers numbers of each state",
			// }, []string{"queue"}),
			// UsedAMResourceMB: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			// 	Namespace: namespace,
			// 	Subsystem: "queue_metrics",
			// 	Name:      "nodemanager_nums",
			// 	Help:      "Current NodeManagers numbers of each state",
			// }, []string{"queue"}),
			// UsedAMResourceVCores: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			// 	Namespace: namespace,
			// 	Subsystem: "queue_metrics",
			// 	Name:      "nodemanager_nums",
			// 	Help:      "Current NodeManagers numbers of each state",
			// }, []string{"queue"}),
			// UsedCapacity: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			// 	Namespace: namespace,
			// 	Subsystem: "queue_metrics",
			// 	Name:      "nodemanager_nums",
			// 	Help:      "Current NodeManagers numbers of each state",
			// }, []string{"queue"}),
			// AbsoluteUsedCapacity: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			// 	Namespace: namespace,
			// 	Subsystem: "queue_metrics",
			// 	Name:      "nodemanager_nums",
			// 	Help:      "Current NodeManagers numbers of each state",
			// }, []string{"queue"}),
			AppsCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "apps_count",
				Help:      "Applications count of each state",
			}, []string{"queue", "state"}),
			AggregateNodeLocalContainersAllocated: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "nodemanager_nums",
				Help:      "Current NodeManagers numbers of each state",
			}, []string{"queue"}),
			AggregateRackLocalContainersAllocated: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "nodemanager_nums",
				Help:      "Current NodeManagers numbers of each state",
			}, []string{"queue"}),
			AggregateOffSwitchContainersAllocated: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "nodemanager_nums",
				Help:      "Current NodeManagers numbers of each state",
			}, []string{"queue"}),
			AggregateContainersPreempted: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "nodemanager_nums",
				Help:      "Current NodeManagers numbers of each state",
			}, []string{"queue"}),
			AggregateMemoryMBSecondsPreempted: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "nodemanager_nums",
				Help:      "Current NodeManagers numbers of each state",
			}, []string{"queue"}),
			AggregateVcoreSecondsPreempted: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "nodemanager_nums",
				Help:      "Current NodeManagers numbers of each state",
			}, []string{"queue"}),
			ActiveUsers: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "nodemanager_nums",
				Help:      "Current NodeManagers numbers of each state",
			}, []string{"queue"}),
			ActiveApplications: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "nodemanager_nums",
				Help:      "Current NodeManagers numbers of each state",
			}, []string{"queue"}),
			AppAttemptFirstContainerAllocationDelayNumOps: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "nodemanager_nums",
				Help:      "Current NodeManagers numbers of each state",
			}, []string{"queue"}),
			AppAttemptFirstContainerAllocationDelayAvgTime: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "nodemanager_nums",
				Help:      "Current NodeManagers numbers of each state",
			}, []string{"queue"}),
			AllocatedMB: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "nodemanager_nums",
				Help:      "Current NodeManagers numbers of each state",
			}, []string{"queue"}),
			AllocatedVCores: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "nodemanager_nums",
				Help:      "Current NodeManagers numbers of each state",
			}, []string{"queue"}),
			AllocatedContainers: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "nodemanager_nums",
				Help:      "Current NodeManagers numbers of each state",
			}, []string{"queue"}),
			AggregateContainersAllocated: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "nodemanager_nums",
				Help:      "Current NodeManagers numbers of each state",
			}, []string{"queue"}),
			AggregateContainersReleased: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "nodemanager_nums",
				Help:      "Current NodeManagers numbers of each state",
			}, []string{"queue"}),
			AvailableMB: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "nodemanager_nums",
				Help:      "Current NodeManagers numbers of each state",
			}, []string{"queue"}),
			AvailableVCores: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "nodemanager_nums",
				Help:      "Current NodeManagers numbers of each state",
			}, []string{"queue"}),
			PendingMB: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "nodemanager_nums",
				Help:      "Current NodeManagers numbers of each state",
			}, []string{"queue"}),
			PendingVCores: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "nodemanager_nums",
				Help:      "Current NodeManagers numbers of each state",
			}, []string{"queue"}),
			PendingContainers: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "nodemanager_nums",
				Help:      "Current NodeManagers numbers of each state",
			}, []string{"queue"}),
			ReservedMB: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "nodemanager_nums",
				Help:      "Current NodeManagers numbers of each state",
			}, []string{"queue"}),
			ReservedVCores: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "nodemanager_nums",
				Help:      "Current NodeManagers numbers of each state",
			}, []string{"queue"}),
			ReservedContainers: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "queue_metrics",
				Name:      "nodemanager_nums",
				Help:      "Current NodeManagers numbers of each state",
			}, []string{"queue"}),
		},
	}
}

// Collect implements the prometheus.Collector interface.
func (e *ResourceManagerMetrics) Collect(ch chan<- prometheus.Metric) {
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

		if DataMap["name"] == "Hadoop:service=ResourceManager,name=ClusterMetrics" {
			e.NodeManagerNums.WithLabelValues("active").Set(DataMap["NumActiveNMs"].(float64))
			e.NodeManagerNums.WithLabelValues("decommissioning").Set(DataMap["NumActiveNMs"].(float64))
			e.NodeManagerNums.WithLabelValues("decommissioned").Set(DataMap["NumDecommissionedNMs"].(float64))
			e.NodeManagerNums.WithLabelValues("lost").Set(DataMap["NumLostNMs"].(float64))
			e.NodeManagerNums.WithLabelValues("unhealthy").Set(DataMap["NumUnhealthyNMs"].(float64))
			e.NodeManagerNums.WithLabelValues("rebooted").Set(DataMap["NumRebootedNMs"].(float64))
			e.NodeManagerNums.WithLabelValues("shutdown").Set(DataMap["NumShutdownNMs"].(float64))

		}

		if strings.HasPrefix(DataMap["name"].(string), "Hadoop:service=ResourceManager,name=QueueMetrics") {
			queue := DataMap["tag.Queue"].(string)
			e.AppsCount.WithLabelValues(queue, "submitted").Set(DataMap["AppsSubmitted"].(float64))
			e.AppsCount.WithLabelValues(queue, "running").Set(DataMap["AppsRunning"].(float64))
			e.AppsCount.WithLabelValues(queue, "pending").Set(DataMap["AppsPending"].(float64))
			e.AppsCount.WithLabelValues(queue, "completed").Set(DataMap["AppsCompleted"].(float64))
			e.AppsCount.WithLabelValues(queue, "killed").Set(DataMap["AppsKilled"].(float64))
			e.AppsCount.WithLabelValues(queue, "failed").Set(DataMap["AppsFailed"].(float64))

		}

	}

	e.NodeManagerNums.Collect(ch)
	e.AppsCount.Collect(ch)
}

func ResourceManagerCollector(target Target, registry *prometheus.Registry) (success bool) {

	metrics := NewResourceManagerMetrics(target)
	registry.MustRegister(metrics)

	return true
}
