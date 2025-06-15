package handlefunc

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	cnst "github.com/RuslanSemenikhin/go-service-with-metrics.git/internal/consts"
)

func CheckTypeMetric(typeMetric string) error {
	if typeMetric != cnst.GAUGE && typeMetric != cnst.COUNTER {
		errString := fmt.Sprintf("Incorrect type metric - '%s'.", typeMetric)
		return errors.New(errString)
	}
	return nil
}

func CheckNameMetric(typeMetric, nameMetric string) error {
	// switch typeMetric {
	// case cnst.COUNTER:
	// 	if ok := slices.Contains(env.CounterMetricsNames, nameMetric); ok {
	// 		return nil
	// 	}
	// case cnst.GAUGE:
	// 	if ok := slices.Contains(env.GaugeMetricsNames, nameMetric); ok {
	// 		return nil
	// 	}
	// default:
	// 	errString := fmt.Sprintf("Incorrect metric type - '%s'", typeMetric)
	// 	return errors.New(errString)
	// }
	if strings.TrimSpace(nameMetric) == "" {
		errString := fmt.Sprintf("Incorrect url-param (metricName) - '%s'. Name metric is missing.", nameMetric)
		// errString := fmt.Sprintf("Metric name - '%s' with metric type - '%s' does not exists", nameMetric, typeMetric)
		return errors.New(errString)
	}
	return nil
}

func CheckCounterValue(val string) (int64, error) {
	intVal, err := strconv.Atoi(val)
	if err != nil {
		errString := fmt.Sprintf("Incorrect value metric - '%s' must have type 'int64'", val)
		return 0, errors.New(errString)
	}
	return int64(intVal), nil
}

func CheckGaugeValue(val string) (float64, error) {
	floatVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		errString := fmt.Sprintf("Incorrect value metric - '%s' must have type 'float64'", val)
		return 0, errors.New(errString)
	}
	return floatVal, nil
}
