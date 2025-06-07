package main

import (
	"github.com/RuslanSemenikhin/go-service-with-metrics.git/internal/managers"
)

var AGENT *managers.AgentManager

func init() {
	realMetricsName := []string{
		"Alloc",
		"BuckHashSys",
		"Frees",
		"GCCPUFraction",
		"GCSys",
		"HeapAlloc",
		"HeapIdle",
		"HeapInuse",
		"HeapObjects",
		"HeapReleased",
		"HeapSys",
		"LastGC",
		"Lookups",
		"MCacheInuse",
		"MCacheSys",
		"MSpanInuse",
		"MSpanSys",
		"Mallocs",
		"NextGC",
		"NumForcedGC",
		"NumGC",
		"OtherSys",
		"PauseTotalNs",
		"StackInuse",
		"StackSys",
		"Sys",
		"TotalAlloc",
		"RandomValue",
	}

	counterMetricsNames := []string{"PollCount"}
	AGENT = managers.NewAgent().
		WithTimeForMetricsTicker(2).
		WithTimeForServerTicker(10).WithMetricsName(realMetricsName, counterMetricsNames)
}

func main() {
	AGENT.Start()
}
