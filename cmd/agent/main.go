package main

import (
	"github.com/RuslanSemenikhin/go-service-with-metrics.git/internal/env"
	"github.com/RuslanSemenikhin/go-service-with-metrics.git/internal/managers"
)

func main() {
	agent := managers.NewAgent().
		WithTimeForMetricsTicker(2).
		WithTimeForServerTicker(10).WithMetricsName(env.GaugeMetricsNames, env.CounterMetricsNames)
	agent.Start()
}
