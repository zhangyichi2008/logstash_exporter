# logstash_exporter 

Exports [Logstash](https://www.elastic.co/logstash/) metrics to [Prometheus](https://prometheus.io/) for monitoring.

**v0.0.2 / 2023-08-11**

**[ADD] 增加对logstash过滤状态pipelines_filters_failures的监控**

- HELP logstash_node_pipelines_filters_failures pipelines_filters_failures
- TYPE logstash_node_pipelines_filters_failures counter
- logstash_node_pipelines_filters_failures{id="apache_access_hr_drop",name="drop"} 0
- logstash_node_pipelines_filters_failures{id="apache_access_hr_mutate",name="mutate"} 0
- logstash_node_pipelines_filters_failures{id="apache_access_hr_useragent",name="useragent"} 0

**[DEL] 移除以go_开头的33个监控项,具体如下：**

- go_gc_duration_seconds{quantile="0"}
- go_gc_duration_seconds{quantile="0.25"}
- go_gc_duration_seconds{quantile="0.5"}
- go_gc_duration_seconds{quantile="0.75"}
- go_gc_duration_seconds{quantile="1"}
- go_gc_duration_seconds_sum
- go_gc_duration_seconds_count
- go_goroutines
- go_info{version="go1.19.1"}
- go_memstats_alloc_bytes
- go_memstats_alloc_bytes_total
- go_memstats_buck_hash_sys_bytes
- go_memstats_frees_total
- go_memstats_gc_sys_bytes
- go_memstats_heap_alloc_bytes
- go_memstats_heap_idle_bytes
- go_memstats_heap_inuse_bytes
- go_memstats_heap_objects
- go_memstats_heap_released_bytes
- go_memstats_heap_sys_bytes
- go_memstats_last_gc_time_seconds
- go_memstats_lookups_total
- go_memstats_mallocs_total
- go_memstats_mcache_inuse_bytes
- go_memstats_mcache_sys_bytes
- go_memstats_mspan_inuse_bytes
- go_memstats_mspan_sys_bytes
- go_memstats_next_gc_bytes
- go_memstats_other_sys_bytes
- go_memstats_stack_inuse_bytes
- go_memstats_stack_sys_bytes
- go_memstats_sys_bytes
- go_threads


###Build and CLI
```sh
cd logstash_exporter/
go build .

logstash_exporter <flags>
  -h, --help                    Show context-sensitive help (also try --help-long and --help-man).
      --logstash.endpoint="http://localhost:9600"  
                                The protocol, host and port on which logstash metrics API listens.
      --listen.address=":9198"  Address on which to expose metrics and web interface.
      --log.level="info"        The logging level to be defined.
      --version                 Show application version.
```

## Metrics

| Name | Type | Description |
| --- | --- | --- |
| `logstash_exporter_build_info` | gauge | Exporter build info |
| `logstash_exporter_scrape_duration_seconds` | summary | Duration of a scrape job. |
| `logstash_info_jvm` | counter |  A metric with a constant '1' value labeled by name, version and vendor of the JVM running Logstash.| 
| `logstash_info_node`| counter |  A metric with a constant '1' value labeled by Logstash version. |
| `logstash_info_os` | counter | A metric with a constant '1' value labeled by name, arch, version and available_processors to the OS running Logstash. |
| `logstash_node_gc_collection_duration_seconds_total` | counter | |
| `logstash_node_gc_collection_total` | gauge | | 
| `logstash_node_jvm_threads_count` | gauge | |
| `logstash_node_jvm_threads_peak_count` | gauge | | 
| `logstash_node_mem_heap_committed_bytes` | gauge | |
| `logstash_node_mem_heap_max_bytes` | gauge | |
| `logstash_node_mem_heap_used_bytes` | gauge | 
| `logstash_node_mem_nonheap_committed_bytes` | gauge | |
| `logstash_node_mem_nonheap_used_bytes` | gauge | |
| `logstash_node_mem_pool_committed_bytes` | gauge | | 
| `logstash_node_mem_pool_max_bytes` | gauge | |
| `logstash_node_mem_pool_peak_max_bytes` | gauge | |
| `logstash_node_mem_pool_peak_used_bytes` | gauge | |
| `logstash_node_mem_pool_used_bytes` | gauge | |
| `logstash_node_pipelines_filters_failures` | counter | |
| `logstash_node_process_cpu_total_seconds_total` | counter | |
| `logstash_node_process_max_filedescriptors` | gauge
| `logstash_node_process_mem_total_virtual_bytes` | gauge
| `logstash_node_process_open_filedescriptors` | gauge
| `logstash_node_queue_events` | counter | |
| `logstash_node_queue_size_bytes` | counter | |
| `logstash_node_queue_max_size_bytes` | counter | |
| `logstash_node_dead_letter_queue_size_bytes` | counter | |
| `logstash_node_up`: | gauge | whether logstash node is up (1) or not (0) |

