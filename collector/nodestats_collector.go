package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

// NodeStatsCollector type
type NodeStatsCollector struct {
	endpoint string

	Up *prometheus.Desc

	JvmThreadsCount     *prometheus.Desc
	JvmThreadsPeakCount *prometheus.Desc

	MemHeapUsedInBytes         *prometheus.Desc
	MemHeapCommittedInBytes    *prometheus.Desc
	MemHeapMaxInBytes          *prometheus.Desc
	MemNonHeapUsedInBytes      *prometheus.Desc
	MemNonHeapCommittedInBytes *prometheus.Desc

	MemPoolPeakUsedInBytes  *prometheus.Desc
	MemPoolUsedInBytes      *prometheus.Desc
	MemPoolPeakMaxInBytes   *prometheus.Desc
	MemPoolMaxInBytes       *prometheus.Desc
	MemPoolCommittedInBytes *prometheus.Desc

	GCCollectionTimeInMillis *prometheus.Desc
	GCCollectionCount        *prometheus.Desc

	ProcessOpenFileDescriptors    *prometheus.Desc
	ProcessMaxFileDescriptors     *prometheus.Desc
	ProcessMemTotalVirtualInBytes *prometheus.Desc
	ProcessCPUTotalInMillis       *prometheus.Desc

	PipelinePluginFailures *prometheus.Desc
}

// NewNodeStatsCollector function
func NewNodeStatsCollector(logstashEndpoint string) (Collector, error) {
	const subsystem = "node"

	return &NodeStatsCollector{
		endpoint: logstashEndpoint,

		Up: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "up"),
			"up",
			nil,
			nil,
		),

		JvmThreadsCount: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "jvm_threads_count"),
			"jvm_threads_count",
			nil,
			nil,
		),

		JvmThreadsPeakCount: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "jvm_threads_peak_count"),
			"jvm_threads_peak_count",
			nil,
			nil,
		),

		MemHeapUsedInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "mem_heap_used_bytes"),
			"mem_heap_used_bytes",
			nil,
			nil,
		),

		MemHeapCommittedInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "mem_heap_committed_bytes"),
			"mem_heap_committed_bytes",
			nil,
			nil,
		),

		MemHeapMaxInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "mem_heap_max_bytes"),
			"mem_heap_max_bytes",
			nil,
			nil,
		),

		MemNonHeapUsedInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "mem_nonheap_used_bytes"),
			"mem_nonheap_used_bytes",
			nil,
			nil,
		),

		MemNonHeapCommittedInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "mem_nonheap_committed_bytes"),
			"mem_nonheap_committed_bytes",
			nil,
			nil,
		),

		MemPoolUsedInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "mem_pool_used_bytes"),
			"mem_pool_used_bytes",
			[]string{"pool"},
			nil,
		),

		MemPoolPeakUsedInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "mem_pool_peak_used_bytes"),
			"mem_pool_peak_used_bytes",
			[]string{"pool"},
			nil,
		),

		MemPoolPeakMaxInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "mem_pool_peak_max_bytes"),
			"mem_pool_peak_max_bytes",
			[]string{"pool"},
			nil,
		),

		MemPoolMaxInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "mem_pool_max_bytes"),
			"mem_pool_max_bytes",
			[]string{"pool"},
			nil,
		),

		MemPoolCommittedInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "mem_pool_committed_bytes"),
			"mem_pool_committed_bytes",
			[]string{"pool"},
			nil,
		),

		GCCollectionTimeInMillis: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "gc_collection_duration_seconds_total"),
			"gc_collection_duration_seconds_total",
			[]string{"collector"},
			nil,
		),

		GCCollectionCount: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "gc_collection_total"),
			"gc_collection_total",
			[]string{"collector"},
			nil,
		),

		ProcessOpenFileDescriptors: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "process_open_filedescriptors"),
			"process_open_filedescriptors",
			nil,
			nil,
		),

		ProcessMaxFileDescriptors: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "process_max_filedescriptors"),
			"process_max_filedescriptors",
			nil,
			nil,
		),

		ProcessMemTotalVirtualInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "process_mem_total_virtual_bytes"),
			"process_mem_total_virtual_bytes",
			nil,
			nil,
		),

		ProcessCPUTotalInMillis: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "process_cpu_total_seconds_total"),
			"process_cpu_total_seconds_total",
			nil,
			nil,
		),

		PipelinePluginFailures: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "pipelines_filters_failures"),
			"pipelines_filters_failures",
			[]string{"name", "id"},
			nil,
		),
	}, nil
}

// Collect function implements nodestats_collector collector
func (c *NodeStatsCollector) Collect(ch chan<- prometheus.Metric) error {
	stats, err := NodeStats(c.endpoint)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(
			c.Up,
			prometheus.GaugeValue,
			0,
		)

		return err
	}

	ch <- prometheus.MustNewConstMetric(
		c.Up,
		prometheus.GaugeValue,
		1,
	)

	c.collectJVM(stats, ch)
	c.collectProcess(stats, ch)

	// For backwards compatibility with Logstash 5
	pipelines := make(map[string]Pipeline)
	if len(stats.Pipelines) == 0 {
		pipelines["main"] = stats.Pipeline
	} else {
		pipelines = stats.Pipelines
	}
	c.collectPipelines(pipelines, ch)

	return nil
}

func (c *NodeStatsCollector) collectJVM(stats NodeStatsResponse, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		c.JvmThreadsCount,
		prometheus.GaugeValue,
		float64(stats.Jvm.Threads.Count),
	)

	ch <- prometheus.MustNewConstMetric(
		c.JvmThreadsPeakCount,
		prometheus.GaugeValue,
		float64(stats.Jvm.Threads.PeakCount),
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemHeapUsedInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.HeapUsedInBytes),
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemHeapMaxInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.HeapMaxInBytes),
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemHeapCommittedInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.HeapCommittedInBytes),
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemNonHeapUsedInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.NonHeapUsedInBytes),
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemNonHeapCommittedInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.NonHeapCommittedInBytes),
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolPeakUsedInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Old.PeakUsedInBytes),
		"old",
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolUsedInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Old.UsedInBytes),
		"old",
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolPeakMaxInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Old.PeakMaxInBytes),
		"old",
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolMaxInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Old.MaxInBytes),
		"old",
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolCommittedInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Old.CommittedInBytes),
		"old",
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolPeakUsedInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Old.PeakUsedInBytes),
		"young",
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolUsedInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Young.UsedInBytes),
		"young",
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolPeakMaxInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Old.PeakMaxInBytes),
		"young",
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolMaxInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Young.MaxInBytes),
		"young",
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolCommittedInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Young.CommittedInBytes),
		"young",
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolPeakUsedInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Old.PeakUsedInBytes),
		"survivor",
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolUsedInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Survivor.UsedInBytes),
		"survivor",
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolPeakMaxInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Old.PeakMaxInBytes),
		"survivor",
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolMaxInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Survivor.MaxInBytes),
		"survivor",
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolCommittedInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Survivor.CommittedInBytes),
		"survivor",
	)

	ch <- prometheus.MustNewConstMetric(
		c.GCCollectionTimeInMillis,
		prometheus.CounterValue,
		float64(stats.Jvm.Gc.Collectors.Old.CollectionTimeInMillis),
		"old",
	)

	ch <- prometheus.MustNewConstMetric(
		c.GCCollectionCount,
		prometheus.GaugeValue,
		float64(stats.Jvm.Gc.Collectors.Old.CollectionCount),
		"old",
	)

	ch <- prometheus.MustNewConstMetric(
		c.GCCollectionTimeInMillis,
		prometheus.CounterValue,
		float64(stats.Jvm.Gc.Collectors.Young.CollectionTimeInMillis),
		"young",
	)

	ch <- prometheus.MustNewConstMetric(
		c.GCCollectionCount,
		prometheus.GaugeValue,
		float64(stats.Jvm.Gc.Collectors.Young.CollectionCount),
		"young",
	)
}

func (c *NodeStatsCollector) collectProcess(stats NodeStatsResponse, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		c.ProcessOpenFileDescriptors,
		prometheus.GaugeValue,
		float64(stats.Process.OpenFileDescriptors),
	)

	ch <- prometheus.MustNewConstMetric(
		c.ProcessMaxFileDescriptors,
		prometheus.GaugeValue,
		float64(stats.Process.MaxFileDescriptors),
	)

	ch <- prometheus.MustNewConstMetric(
		c.ProcessMemTotalVirtualInBytes,
		prometheus.GaugeValue,
		float64(stats.Process.Mem.TotalVirtualInBytes),
	)

	ch <- prometheus.MustNewConstMetric(
		c.ProcessCPUTotalInMillis,
		prometheus.CounterValue,
		float64(stats.Process.CPU.TotalInMillis)/1000,
	)
}

func (c *NodeStatsCollector) collectPipelines(pipelines map[string]Pipeline, ch chan<- prometheus.Metric) {
	for _, pipeline := range pipelines {
		for _, plugin := range pipeline.Plugins.Filters {
			ch <- prometheus.MustNewConstMetric(
				c.PipelinePluginFailures,
				prometheus.CounterValue,
				float64(plugin.Failures),
				plugin.Name,
				plugin.ID,
			)
		}
	}
}
