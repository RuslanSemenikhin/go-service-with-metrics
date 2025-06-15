package managers

import (
	"sync"

	"github.com/RuslanSemenikhin/go-service-with-metrics.git/internal/interfaces"
)

type GaugeManage struct {
	mtx     sync.RWMutex
	storage interfaces.StorageInterface
}

func NewGaugeManager() *GaugeManage {
	return &GaugeManage{}
}

func (gm *GaugeManage) WithStorage(newStorage interfaces.StorageInterface) *GaugeManage {
	gm.mtx.Lock()
	defer gm.mtx.Unlock()
	gm.storage = newStorage
	return gm
}

func (gm *GaugeManage) UpdateGauge(name string, val float64) {
	gm.storage.UpdateGauge(name, val)
}

func (gm *GaugeManage) ShowAllGauges() map[string][]float64 {
	return gm.storage.GetGaugeHistory()
}

func (gm *GaugeManage) GetActualGaugesMetrics() map[string]float64 {
	return gm.storage.GetActualGauges()
}

func (cm *GaugeManage) GetMetricValueByName(name string) (float64, error) {
	return cm.storage.GetGaugesByName(name)
}
