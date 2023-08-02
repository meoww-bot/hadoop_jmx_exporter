# hadoop_jmx_exporter

一个开箱即用、支持 [multi-target](https://prometheus.io/docs/guides/multi-target-exporter/) 的 hadoop jmx 导出工具，通过 url GET 传参解析对应的 jmx 地址，返回 prometheus 风格的指标内容。

因为 jmx 地址可以通过网络访问，所以不必在机器上运行程序来导出指标，减少对集群节点的侵入。

已经在 HDP 3.1.5 上测试。

## Feature

1. 预请求 jmx 内容，自动识别 jmx 类型，使用对应的 collector 解析
2. 支持 kerberos，密码认证和 keytab 认证，取决于配置的参数是 `ktpath` 还是 `password`

❌ 暂不支持自定义指标  

因为强迫症，有的多个指标实际属于同一指标的不同维度，这不太好配置

### Support Service

|service|role|support|
|-|-|-|
|HDFS|NameNode|✅|
|HDFS|DataNode|✅|
|HDFS|JournalNode|✅|
|HBASE|HbaseMaster|✅|
|HBASE|RegionServer|✅|
|YARN|ResourceManager|✅|
|YARN|NodeManager|✅|
|HIVE|HiveServer2|✅|


## Build

make build

## Run


1. 添加一个低权限用户
```
useradd -rs /bin/false nodeusr
```

2. 将二进制文件放到 /usr/local/bin/hadoop_jmx_exporter

3. 配置 hadoop_jmx_exporter.service

```
[Unit]
Description=Hadoop Jmx Exporter
After=network-online.target

[Service]
Type=simple
User=nodeusr
Group=nodeusr
ExecStart=/usr/local/bin/hadoop_jmx_exporter
KillMode=process
RemainAfterExit=no
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target

```

4. service hadoop_jmx_exporter start

## Prometheus Configuration

without kerberos （HDP 2.6.4 访问 jmx 无需 kerberos 认证）
```
  - job_name: 'hadoop_jmx_exporter'
    scrape_interval: 30s
    metrics_path: /scrape
    static_configs:
      - targets:
        - http://yarn-rm.example.com:8088/jmx
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        regex: "http://([^/:]+):\\d+/jmx"
        target_label: instance
      - target_label: __address__
        replacement: 127.0.0.1:9070 # hadoop_jmx_exporter 服务所在的机器和端口
```

kerberos keytab auth （HDP 3.1.5）
```
  - job_name: 'hadoop_jmx_exporter'
    scrape_interval: 30s
    metrics_path: /scrape
    params:
      ktpath:  
      - /etc/xxxxx.keytab
      principal:  
      - xxxxx@EXAMPLE.COM
    static_configs:
      - targets:
        - http://yarn-rm.example.com:8088/jmx
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        regex: "http://([^/:]+):\\d+/jmx"
        target_label: instance
      - target_label: __address__
        replacement: 127.0.0.1:9070 # hadoop_jmx_exporter 服务所在的机器和端口
```

kerberos password auth （HDP 3.1.5）
```
  - job_name: 'hadoop_jmx_exporter'
    scrape_interval: 30s
    metrics_path: /scrape
    params:
      principal:  
      - xxxxx@EXAMPLE.COM
      password:  
      - "yourpassword"

    static_configs:
      - targets:
        - http://yarn-rm.example.com:8088/jmx
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        regex: "http://([^/:]+):\\d+/jmx"
        target_label: instance
      - target_label: __address__
        replacement: 127.0.0.1:9070 # hadoop_jmx_exporter 服务所在的机器和端口
```

如果你有多个集群，而  jmx 没有集群名，如何区分不同集群的指标？

比如 ResourceManager 的 jmx 没有代表集群名的指标，暂时考虑在 prometheus 人为新增一个标签，比如
```
  - job_name: 'hadoop_jmx_exporter'
    scrape_interval: 30s
    metrics_path: /scrape
    params:
      principal:  
      - xxxxx@EXAMPLE.COM
      password:  
      - "yourpassword"

    static_configs:
      - targets:
        - http://yarn-rm.example.com:8088/jmx
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        regex: "http://([^/:]+):\\d+/jmx"
        target_label: instance
      - target_label: __address__
        replacement: 127.0.0.1:9070 # hadoop_jmx_exporter 服务所在的机器和端口
      - source_labels: []
        target_label: cluster
        replacement: "hadoop1"

```

也可以在 job name 写 集群名，怎么方便怎么来

## Metrics Map

指标定义准则

1. 将同一指标的不同维度放到标签里面，降低基数
2. 指标定义： `<hadoop service>_<component>_<jmx beans modelerType>_<metrics>`
     

    比如：BlocksTotal 对应的 prometheus 指标 `hdfs_namenode_fsname_system_blocks_total` 
    
        - hadoop service: hdfs  
        - component: namenode  
        - jmx beans modelerType: FSNamesystem -> fsname_system  
        - metrics: BlocksTotal -> blocks_total  
3. prometheus 指标全部是小写字母，使用 `_` 下划线分隔
4. 如果指标有单位，尽量带单位，比如 count，milliseconds，bytes


### NameNode

#### Hadoop:service=NameNode,name=FSNamesystem

|Jmx Metric|Prometheus Metric|Description|Chinese Description|
|-|-|-|-|
|MissingBlocks|hdfs_namenode_fsname_system_missing_blocks|Current number of missing blocks|
|UnderReplicatedBlocks|hdfs_namenode_fsname_system_under_replicated_blocks|Current number of blocks under replicated 
|CapacityTotal|hdfs_namenode_fsname_system_capacity_bytes{mode="Total"}|Current raw capacity of DataNodes in bytes
|CapacityUsed|hdfs_namenode_fsname_system_capacity_bytes{mode="Used"}|Current used capacity across all DataNodes in bytes
|CapacityRemaining|hdfs_namenode_fsname_system_capacity_bytes{mode="Remaining"}|Current remaining capacity in bytes
|CapacityUsedNonDFS|hdfs_namenode_fsname_system_capacity_bytes{mode="UsedNonDFS"}|Current space used by DataNodes for non DFS purposes in bytes
|BlocksTotal|hdfs_namenode_fsname_system_blocks_total|Current number of allocated blocks in the system
|FilesTotal|hdfs_namenode_fsname_system_files_total|Current number of files and directories
|CorruptBlocks|hdfs_namenode_fsname_system_corrupt_blocks|Current number of blocks with corrupt replicas
|ExcessBlocks|hdfs_namenode_fsname_system_excess_blocks|Current number of excess blocks
|StaleDataNodes|hdfs_namenode_fsname_system_stale_datanodes|Current number of DataNodes marked stale due to delayed heartbeat
|tag.HAState|hdfs_namenode_fsname_system_hastate|(HA-only) Current state of the NameNode: initializing or active or standby or stopping state |


#### Hadoop:service=NameNode,name=JvmMetrics

|Jmx Metric|Prometheus Metric|Description|Chinese Description|
|-|-|-|-|
|GcCountParNew|hdfs_namenode_jvm_metrics_gc_count{type="ParNew"}|ParNew GC count
|GcCountConcurrentMarkSweep|hdfs_namenode_jvm_metrics_gc_count{type="ConcurrentMarkSweep"}|ConcurrentMarkSweep GC count
|GcTimeMillisParNew|hdfs_namenode_jvm_metrics_gc_time_milliseconds{type="ParNew"}|ParNew GC time in milliseconds
|GcTimeMillisConcurrentMarkSweep|hdfs_namenode_jvm_metrics_gc_time_milliseconds{type="ConcurrentMarkSweep"}|ConcurrentMarkSweep GC time in milliseconds


#### java.lang:type=Memory

|Jmx Metric|Prometheus Metric|Description|Chinese Description|
|-|-|-|-|
|HeapMemoryUsage{committed}|hdfs_namenode_memory_heap_memory_usage_bytes{mode="committed"}|
|HeapMemoryUsage{init}|hdfs_namenode_memory_heap_memory_usage_bytes{mode="init"}|
|HeapMemoryUsage{max}|hdfs_namenode_memory_heap_memory_usage_bytes{mode="max"}|
|HeapMemoryUsage{used}|hdfs_namenode_memory_heap_memory_usage_bytes{mode="used"}|

#### Hadoop:service=NameNode,name=NameNodeStatus

|Jmx Metric|Prometheus Metric|Description|Chinese Description|
|-|-|-|-|
|LastHATransitionTime|hdfs_namenode_namenode_status_last_ha_transition_time|


####  Hadoop:service=NameNode,name=RpcActivityForPort8020/8060

|Jmx Metric|Prometheus Metric|Description|Chinese Description|
|-|-|-|-|
|ReceivedBytes|hdfs_namenode_rpc_activity_received_bytes|Total number of received bytes
|SentBytes|hdfs_namenode_rpc_activity_sent_bytes|Total number of sent bytes
|RpcQueueTimeNumOps|hdfs_namenode_rpc_activity_call_count{method="QueueTime"}|Total number of RPC calls 
|RpcQueueTimeAvgTime|hdfs_namenode_rpc_activity_avg_time_milliseconds{method="RpcQueueTime"}|Average queue time in milliseconds 
|RpcProcessingTimeAvgTime|hdfs_namenode_rpc_activity_avg_time_milliseconds{method="RpcProcessingTime"}|Average Processing time in milliseconds
|NumOpenConnections|hdfs_namenode_rpc_activity_open_connections_count|Current number of open connections
|CallQueueLength|hdfs_namenode_rpc_activity_call_queue_length|Current length of the call queue

