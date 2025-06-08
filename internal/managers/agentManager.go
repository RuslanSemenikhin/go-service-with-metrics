package managers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"sync"
	"syscall"
	"time"
)

type AgentManager struct {
	mtx               sync.RWMutex
	metricsTickerTime time.Duration
	serverTickerTime  time.Duration
	metricsTicker     *time.Ticker
	serverTicker      *time.Ticker
	gauge             map[string][]float64
	counter           map[string]int64
	stopChan          chan os.Signal
}

func NewAgent() *AgentManager {
	return &AgentManager{
		gauge:    make(map[string][]float64),
		counter:  make(map[string]int64),
		stopChan: make(chan os.Signal, 1),
	}
}

func (am *AgentManager) newMetricsTicker(t time.Duration) {
	ticker := time.NewTicker(t)
	am.metricsTicker = ticker
}

func (am *AgentManager) newServerTicker(t time.Duration) {
	ticker := time.NewTicker(t)
	am.serverTicker = ticker
}

func (am *AgentManager) Start() {
	timeMetrics, timeServer := am.GetMetricsTime(), am.GetServerTime()
	am.newMetricsTicker(timeMetrics)
	defer am.metricsTicker.Stop()

	am.newServerTicker(timeServer)
	defer am.serverTicker.Stop()

	signal.Notify(am.stopChan, syscall.SIGINT, syscall.SIGTERM)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case mt := <-am.metricsTicker.C:
				log.Printf(`metrics Ticker passed - '%v'.`, mt)
				am.collectMetrics()
			case st := <-am.serverTicker.C:
				log.Printf(`server Ticker passed - '%v'.`, st)
				am.sendMetricsOnServer()
			case <-am.stopChan:
				log.Println(`received stop signal...`)
				am.sendMetricsOnServer()
				return
			}
		}
	}()

	wg.Wait()
}

func (am *AgentManager) SetTimeForMetricsTicker(t time.Duration) *AgentManager {
	am.mtx.Lock()
	defer am.mtx.Unlock()
	am.metricsTickerTime = t
	return am
}

func (am *AgentManager) SetTimeForServerTicker(t time.Duration) *AgentManager {
	am.mtx.Lock()
	defer am.mtx.Unlock()
	am.serverTickerTime = t
	return am
}

func (am *AgentManager) SetGauge(newMetricsName []string) {
	am.mtx.Lock()
	defer am.mtx.Unlock()
	for _, name := range newMetricsName {
		am.gauge[name] = make([]float64, 0)
	}
}

func (am *AgentManager) SetCounter(newMetricsName []string) {
	am.mtx.Lock()
	defer am.mtx.Unlock()
	for _, name := range newMetricsName {
		am.counter[name] = 0
	}
}

func (am *AgentManager) WithMetricsName(realNames []string, otherName []string) *AgentManager {
	am.SetGauge(realNames)
	am.SetCounter(otherName)
	return am
}

func (am *AgentManager) WithTimeForMetricsTicker(seconds int) *AgentManager {
	sec := time.Duration(seconds) * time.Second
	am.SetTimeForMetricsTicker(sec)
	return am
}

func (am *AgentManager) WithTimeForServerTicker(seconds int) *AgentManager {
	sec := time.Duration(seconds) * time.Second
	am.SetTimeForServerTicker(sec)
	return am
}

func (am *AgentManager) GetMetricsTime() time.Duration {
	am.mtx.Lock()
	defer am.mtx.Unlock()
	return am.metricsTickerTime
}

func (am *AgentManager) GetServerTime() time.Duration {
	am.mtx.Lock()
	defer am.mtx.Unlock()
	return am.serverTickerTime
}

func (am *AgentManager) collectMetrics() {
	am.mtx.Lock()
	defer am.mtx.Unlock()

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	v := reflect.ValueOf(mem)
	t := v.Type()

	for key, slc := range am.gauge {
		currentVal := v.FieldByName(key)
		if currentType, ok := t.FieldByName(key); ok {
			switch currentType.Type.Kind() {
			case reflect.Float64:
				slc = append(slc, currentVal.Float())
				am.gauge[key] = slc
			case reflect.Uint64:
				slc = append(slc, float64(currentVal.Uint()))
				am.gauge[key] = slc
			case reflect.Uint32:
				slc = append(slc, float64(currentVal.Uint()))
				am.gauge[key] = slc
			default:
				log.Printf("bad dataType - '%s', nameMetric - '%s', value - '%v'", currentType.Name, key, currentVal)
			}
		} else {
			randVal := rand.Float64()
			slc = append(slc, randVal)
			am.gauge[key] = slc
		}
	}

	for key := range am.counter {
		am.counter[key] += 1
	}
}

func (am *AgentManager) clearGaugeMetric(name string) {
	am.gauge[name] = make([]float64, 0)
}

func (am *AgentManager) sendGaugeMetric(name string, vals []float64) {
	am.mtx.Lock()
	defer am.mtx.Unlock()
	for _, val := range vals {
		url := fmt.Sprintf("http://localhost:8080/update/gauge/%s/%v", name, val)
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(nil))
		if err != nil {
			log.Printf("can`t create request: name metric - '%s', value - '%v', url - '%s'", name, val, url)
			continue
		}

		req.Header.Set("Content-Type", "text/plain")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("can`t send request: name metric - '%s', value - '%v', url - '%s'", name, val, url)
			continue
		}

		respCode := resp.StatusCode
		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Printf("can`t read response body: name metric - '%s', value - '%v', url - '%s', respCode - '%d'", name, val, url, respCode)
			continue
		}
		log.Printf("code - '%d', body - '%s', content-type - '%s'", respCode, respBody, resp.Header.Get("content-type"))
	}
	am.clearGaugeMetric(name)
}

func (am *AgentManager) sendCounterMetric(name string, val int64) {
	am.mtx.RLock()
	defer am.mtx.RUnlock()
	url := fmt.Sprintf("http://localhost:8080/update/counter/%s/%d", name, val)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(nil))
	if err != nil {
		log.Printf("can`t create request: name metric - '%s', value - '%v', url - '%s'", name, val, url)
		return
	}

	req.Header.Set("Content-Type", "text/plain")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("can`t send request: name metric - '%s', value - '%v', url - '%s'", name, val, url)
		return
	}

	respCode := resp.StatusCode
	respBody, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Printf("can`t read response body: name metric - '%s', value - '%v', url - '%s', respCode - '%d'", name, val, url, respCode)
		return
	}
	log.Printf("code - '%d', body - '%s', content-type - '%s'", respCode, respBody, resp.Header.Get("content-type"))
}

func (am *AgentManager) sendMetricsOnServer() {
	for key, vals := range am.gauge {
		am.sendGaugeMetric(key, vals)
	}

	for key, vals := range am.counter {
		am.sendCounterMetric(key, vals)
	}
}
