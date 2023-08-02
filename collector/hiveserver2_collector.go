package collector

import (
	"encoding/json"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
)

type HiveServer2Metrics struct {
	BaseMetrics
	OpenConnectionsCount          prometheus.Gauge
	JvmPauseExtraSleepTime        prometheus.Gauge
	OpenOperationsCount           prometheus.Gauge
	CumulativeConnectionCount     prometheus.Gauge
	MetastoreHiveLocksCount       prometheus.Gauge
	ExecAsyncQueueSize            prometheus.Gauge
	ExecAsyncPoolSize             prometheus.Gauge
	WaitingCompileOps             prometheus.Gauge
	HiveTezTasks                  prometheus.Gauge
	ActiveCallsApiHs2Operation    *prometheus.GaugeVec
	ActiveCallsApiHs2SqlOperation *prometheus.GaugeVec
	ApiHs2Operation               *prometheus.GaugeVec
	ApiHs2SqlOperation            *prometheus.GaugeVec
	Hs2CompletedOperation         *prometheus.GaugeVec
	Hs2CompletedSqlOperation      *prometheus.GaugeVec
}

func NewHiveServer2Metrics(t Target) *HiveServer2Metrics {

	const namespace = "hive_hiveserver2"

	return &HiveServer2Metrics{
		BaseMetrics: BuildBaseMetrics(t.BodyData, namespace),
		OpenConnectionsCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "metrics",
			Name:      "open_connections_count",
			Help:      "open_connections",
		}),
		JvmPauseExtraSleepTime: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "jvm",
			Name:      "pause_extra_sleep_time_count_milliseconds",
			Help:      "GC 额外睡眠时间",
		}),
		OpenOperationsCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "metrics",
			Name:      "open_operations_count",
			Help:      "open_operations",
		}),
		CumulativeConnectionCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "metrics",
			Name:      "cumulative_connection_count",
			Help:      "累计连接数",
		}),
		MetastoreHiveLocksCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "metrics",
			Name:      "metastore_hive_locks",
			Help:      "metastore_hive_locks",
		}),
		ExecAsyncQueueSize: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "metrics",
			Name:      "exec_async_queue_size",
			Help:      "hs2 异步操作队列当前大小",
		}),
		ExecAsyncPoolSize: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "metrics",
			Name:      "exec_async_pool_size",
			Help:      "hs2 异步线程池当前大小",
		}),
		WaitingCompileOps: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "metrics",
			Name:      "waiting_compile_ops",
			Help:      "waiting_compile_ops",
		}),
		HiveTezTasks: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "metrics",
			Name:      "hive_tez_tasks",
			Help:      "提交的 Hive on Tez 作业总数",
		}),
		ActiveCallsApiHs2Operation: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "metrics",
			Name:      "active_calls_api_hs2_operation",
			Help:      "active_calls_api_hs2_operation",
		}, []string{"state"}),
		ActiveCallsApiHs2SqlOperation: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "metrics",
			Name:      "active_calls_api_hs2_sql_operation",
			Help:      "active_calls_api_hs2_sql_operation",
		}, []string{"state"}),
		ApiHs2Operation: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "metrics",
			Name:      "api_hs2_operation",
			Help:      "api_hs2_operation",
		}, []string{"state"}),
		Hs2CompletedOperation: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "metrics",
			Name:      "hs2_completed_operation",
			Help:      "hs2_completed_operation",
		}, []string{"state"}),
		Hs2CompletedSqlOperation: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "metrics",
			Name:      "hs2_completed_sql_operation",
			Help:      "hs2_completed_sql_operation",
		}, []string{"state"}),
	}
}

// Collect implements the prometheus.Collector interface.
func (e *HiveServer2Metrics) Collect(ch chan<- prometheus.Metric) {

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

		if DataMap["name"] == "java.lang:type=Memory" {
			heapMemoryUsage := DataMap["HeapMemoryUsage"].(map[string]interface{})
			e.HeapMemoryUsage.WithLabelValues("committed").Set(heapMemoryUsage["committed"].(float64))
			e.HeapMemoryUsage.WithLabelValues("init").Set(heapMemoryUsage["init"].(float64))
			e.HeapMemoryUsage.WithLabelValues("max").Set(heapMemoryUsage["max"].(float64))
			e.HeapMemoryUsage.WithLabelValues("used").Set(heapMemoryUsage["used"].(float64))

			e.HeapMemoryUsage.Collect(ch)

		}

		if DataMap["name"] == "metrics:name=jvm.pause.extraSleepTime" {
			e.JvmPauseExtraSleepTime.Set(DataMap["Count"].(float64))
			e.JvmPauseExtraSleepTime.Collect(ch)
		}

		if DataMap["name"] == "metrics:name=open_connections" {
			e.OpenConnectionsCount.Set(DataMap["Count"].(float64))
			e.OpenConnectionsCount.Collect(ch)
		}

		if DataMap["name"] == "metrics:name=open_operations" {
			e.OpenOperationsCount.Set(DataMap["Count"].(float64))
			e.OpenOperationsCount.Collect(ch)
		}

		if DataMap["name"] == "cumulative_connection_count" {
			e.CumulativeConnectionCount.Set(DataMap["Count"].(float64))
			e.CumulativeConnectionCount.Collect(ch)
		}

		if DataMap["name"] == "metrics:name=metastore_hive_locks" {
			e.MetastoreHiveLocksCount.Set(DataMap["Count"].(float64))
			e.MetastoreHiveLocksCount.Collect(ch)
		}

		if DataMap["name"] == "metrics:name=exec_async_queue_size" {
			e.ExecAsyncQueueSize.Set(DataMap["Value"].(float64))
			e.ExecAsyncQueueSize.Collect(ch)
		}

		if DataMap["name"] == "metrics:name=exec_async_pool_size" {
			e.ExecAsyncPoolSize.Set(DataMap["Value"].(float64))
			e.ExecAsyncPoolSize.Collect(ch)
		}
		if DataMap["name"] == "metrics:name=waiting_compile_ops" {
			e.WaitingCompileOps.Set(DataMap["Count"].(float64))
			e.WaitingCompileOps.Collect(ch)
		}
		if DataMap["name"] == "metrics:name=hive_tez_tasks" {
			e.HiveTezTasks.Set(DataMap["Count"].(float64))
			e.HiveTezTasks.Collect(ch)
		}

		if strings.HasPrefix(DataMap["name"].(string), "metrics:name=active_calls_api_hs2_operation_") {
			state := strings.ToLower(getLastUpperWithDelimiter(DataMap["name"].(string), "_"))
			e.ActiveCallsApiHs2Operation.WithLabelValues(state).Set(DataMap["Count"].(float64))
		}
		if strings.HasPrefix(DataMap["name"].(string), "metrics:name=active_calls_api_hs2_sql_operation_") {
			state := strings.ToLower(getLastUpperWithDelimiter(DataMap["name"].(string), "_"))
			e.ActiveCallsApiHs2SqlOperation.WithLabelValues(state).Set(DataMap["Count"].(float64))
		}
		if strings.HasPrefix(DataMap["name"].(string), "metrics:name=api_hs2_sql_operation_") {
			state := strings.ToLower(getLastUpperWithDelimiter(DataMap["name"].(string), "_"))
			e.ApiHs2SqlOperation.WithLabelValues(state).Set(DataMap["Count"].(float64))
		}
		if strings.HasPrefix(DataMap["name"].(string), "metrics:name=hs2_completed_operation_") {
			state := strings.ToLower(getLastUpperWithDelimiter(DataMap["name"].(string), "_"))
			e.Hs2CompletedOperation.WithLabelValues(state).Set(DataMap["Count"].(float64))
		}
		if strings.HasPrefix(DataMap["name"].(string), "metrics:name=hs2_completed_sql_operation_") {
			state := strings.ToLower(getLastUpperWithDelimiter(DataMap["name"].(string), "_"))
			e.Hs2CompletedSqlOperation.WithLabelValues(state).Set(DataMap["Count"].(float64))
		}

	}
	e.ActiveCallsApiHs2Operation.Collect(ch)
	e.ActiveCallsApiHs2SqlOperation.Collect(ch)
	e.ApiHs2SqlOperation.Collect(ch)
	e.Hs2CompletedOperation.Collect(ch)
	e.Hs2CompletedSqlOperation.Collect(ch)

}

func getLastUpperWithDelimiter(s string, delimiter string) string {
	// 使用分隔符将字符串拆分成多个子字符串
	parts := strings.Split(s, delimiter)

	// 取最后一个子字符串作为最后的大写部分
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	// 如果字符串没有分隔符，则返回空字符串
	return ""
}

func HiveServer2Collector(target Target, registry *prometheus.Registry) (success bool) {

	metrics := NewHiveServer2Metrics(target)
	registry.MustRegister(metrics)

	return true
}
