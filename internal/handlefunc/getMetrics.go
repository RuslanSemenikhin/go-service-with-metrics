package handlefunc

import (
	"fmt"
	"net/http"

	cnst "github.com/RuslanSemenikhin/go-service-with-metrics.git/internal/consts"
	"github.com/RuslanSemenikhin/go-service-with-metrics.git/internal/env"
	g "github.com/gin-gonic/gin"
)

func GetAllMetrics(b *env.Box) g.HandlerFunc {
	return func(ctx *g.Context) {
		actualCounters := b.GetCounterManager().GetActualCountersMetrics()
		actualGauges := b.GetGaugeManager().GetActualGaugesMetrics()
		headerCounter := "Type metrics - 'COUNTER'"
		headerGauge := "Type metrics - 'GAUGE'"
		ctx.HTML(http.StatusOK, `allMetrics.tmpl`, g.H{
			`Title`:           `Main page`,
			`Header`:          `All metrics:`,
			`Counter`:         headerCounter,
			`Gauge`:           headerGauge,
			`CountersMetrics`: actualCounters,
			`GaugesMetrics`:   actualGauges,
		})
	}
}

func GetMetricValueByName(b *env.Box) g.HandlerFunc {
	return func(ctx *g.Context) {
		var (
			err error
		)
		metricType := ctx.Param(`metricType`)
		metricName := ctx.Param(`metricName`)

		err = CheckTypeMetric(metricType)
		if err != nil {
			ctx.String(http.StatusNotFound, err.Error())
			return
		}

		err = CheckNameMetric(metricType, metricName)
		if err != nil {
			ctx.String(http.StatusNotFound, err.Error())
			return
		}

		switch metricType {
		case cnst.COUNTER:
			val, err := b.GetCounterManager().GetMetricValueByName(metricName)
			if err != nil {
				ctx.String(http.StatusNotFound, err.Error())
				return
			}
			ctx.String(http.StatusOK, fmt.Sprintf("%d", val))
			return
		case cnst.GAUGE:
			val, err := b.GetGaugeManager().GetMetricValueByName(metricName)
			if err != nil {
				ctx.String(http.StatusNotFound, err.Error())
				return
			}
			ctx.String(http.StatusOK, fmt.Sprintf("%v", val))
			return
		}
	}
}
