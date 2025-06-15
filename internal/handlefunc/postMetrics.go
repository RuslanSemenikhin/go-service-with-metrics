package handlefunc

import (
	"fmt"
	"net/http"

	cnst "github.com/RuslanSemenikhin/go-service-with-metrics.git/internal/consts"
	"github.com/RuslanSemenikhin/go-service-with-metrics.git/internal/env"
	g "github.com/gin-gonic/gin"
)

func Update(box *env.Box) g.HandlerFunc {
	return func(ctx *g.Context) {
		// switch ctx.Request.Header.Get("content-type") {
		// case "text/plain":

		metricType := ctx.Param(`metricType`)
		metricName := ctx.Param(`metricName`)
		metricVal := ctx.Param(`metricValue`)

		err := CheckTypeMetric(metricType)
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		err = CheckNameMetric(metricType, metricName)
		if err != nil {
			ctx.String(http.StatusNotFound, err.Error())
			return
		}

		switch metricType {
		case cnst.COUNTER:
			correctVal, err := CheckCounterValue(metricVal)
			if err != nil {
				ctx.String(http.StatusBadRequest, err.Error())
			}
			box.GetCounterManager().UpdateCounter(metricName, correctVal)
		case cnst.GAUGE:
			correctVal, err := CheckGaugeValue(metricVal)
			if err != nil {
				ctx.String(http.StatusBadRequest, err.Error())
			}
			box.GetGaugeManager().UpdateGauge(metricName, correctVal)
		}

		respText := fmt.Sprintf("Metric with type - '%s' with name - '%s' with value - '%s' added successfuly.", metricType, metricName, metricVal)
		ctx.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(respText))
		// return

		// default:
		// 	ctx.String(http.StatusBadRequest, "content-type - '%s' not supported.", ctx.Request.Header.Get("content-type"))
		// 	return
		// }
	}
}
