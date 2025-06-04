package handlefunc

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	cnst "github.com/RuslanSemenikhin/go-service-with-metrics.git/internal/consts"
	"github.com/RuslanSemenikhin/go-service-with-metrics.git/internal/env"
)

func sharePath(req *http.Request) []string {
	pathString := req.URL.Path
	trimPath := strings.Trim(pathString, "/")
	slcPath := strings.Split(trimPath, "/")
	return slcPath
}

func checkValueAndQuantityParams(partsUrlPath []string) error {
	if len(partsUrlPath) < 4 {
		errString := "value metric was not received, must have value metric with type - 'float64'"
		return errors.New(errString)
	} else if len(partsUrlPath) > 4 {
		errString := "many parameters passed, request must have - 4 params"
		return errors.New(errString)
	} else {
		return nil
	}
}

func withMetricType(req *http.Request, partsUrlPath []string) (string, int, error) {
	if len(partsUrlPath) < 2 || strings.TrimSpace(partsUrlPath[1]) == "" {
		errString := fmt.Sprintf("Incorrect url-path - '%s'. Type metric is missing.", req.URL.Path)
		return "", http.StatusBadRequest, errors.New(errString)
	}

	metricType := partsUrlPath[1]
	if metricType != cnst.GAUGE && metricType != cnst.COUNTER {
		errString := fmt.Sprintf("Incorrect type metric - '%s'.", metricType)
		return "", http.StatusBadRequest, errors.New(errString)
	}

	return metricType, 0, nil
}

func withMetricName(req *http.Request, partsUrlPath []string) (string, int, error) {
	if len(partsUrlPath) < 3 || strings.TrimSpace(partsUrlPath[2]) == "" {
		errString := fmt.Sprintf("Incorrect url-path - '%s'. Name metric is missing.", req.URL.Path)
		return "", http.StatusNotFound, errors.New(errString)
	}

	// TODO: maybe checking pattern name (Regexp)...
	metricName := strings.TrimSpace(partsUrlPath[2])
	return metricName, 0, nil
}

func withGaugeValue(partsUrlPath []string) (float64, int, error) {
	err := checkValueAndQuantityParams(partsUrlPath)
	if err != nil {
		return 0, http.StatusBadRequest, err
	}
	floatVal, err := strconv.ParseFloat(partsUrlPath[3], 64)
	if err != nil {
		errString := fmt.Sprintf("Incorrect value metric - '%s' must have type 'float64'", partsUrlPath[3])
		return 0, http.StatusBadRequest, errors.New(errString)
	}
	return floatVal, 0, nil
}

func withCounterValue(partsUrlPath []string) (int64, int, error) {
	err := checkValueAndQuantityParams(partsUrlPath)
	if err != nil {
		return 0, http.StatusBadRequest, err
	}
	intVal, err := strconv.Atoi(partsUrlPath[3])
	if err != nil {
		errString := fmt.Sprintf("Incorrect value metric - '%s' must have type 'int64'", partsUrlPath[3])
		return 0, http.StatusBadRequest, errors.New(errString)
	}
	return int64(intVal), 0, nil
}

func Update(box *env.Box) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			switch req.Header.Get("content-type") {
			case "text/plain":

				var (
					partsUrlPath []string
					metricType   string
					metricName   string
					metricVal    string
				)

				partsUrlPath = sharePath(req)
				metricType, statusHttp, err := withMetricType(req, partsUrlPath)
				if err != nil {
					http.Error(resp, err.Error(), statusHttp)
					return
				}

				metricName, statusHttp, err = withMetricName(req, partsUrlPath)
				if err != nil {
					http.Error(resp, err.Error(), statusHttp)
					return
				}

				switch metricType {
				case cnst.COUNTER:
					metricValue, statusHttp, err := withCounterValue(partsUrlPath)
					if err != nil {
						http.Error(resp, err.Error(), statusHttp)
						return
					}
					box.GetCounterManager().UpdateCounter(metricName, metricValue)
					metricVal = fmt.Sprintf("%v", metricValue)
				case cnst.GAUGE:
					metricValue, statusHttp, err := withGaugeValue(partsUrlPath)
					if err != nil {
						http.Error(resp, err.Error(), statusHttp)
						return
					}
					box.GetGaugeManager().UpdateGauge(metricName, metricValue)
					metricVal = fmt.Sprintf("%v", metricValue)
				default:
					errString := fmt.Sprintf("Incorrect type metric - '%s'.", metricType)
					http.Error(resp, errString, http.StatusBadRequest)
					return
				}
				resp.Header().Set("content-type", "text/plain; charset=utf-8")
				resp.WriteHeader(http.StatusOK)
				respText := fmt.Sprintf("Metric with type - '%s' with name - '%s' with value - '%s' added successfuly.", metricType, metricName, metricVal)
				resp.Write([]byte(respText))
				return
			default:
				errString := fmt.Sprintf("content-type - '%s' not supported.", req.Header.Get("content-type"))
				http.Error(resp, errString, http.StatusBadRequest)
				return
			}
		default:
			errString := fmt.Sprintf(`method - '%s' does not supported`, req.Method)
			http.Error(resp, errString, http.StatusMethodNotAllowed)
			return
		}
	}
}
