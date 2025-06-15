package managers

import (
	"sync"

	"github.com/RuslanSemenikhin/go-service-with-metrics.git/internal/interfaces"
)

type CaunterManage struct {
	mtx     sync.RWMutex
	storage interfaces.StorageInterface
}

func NewCaunterManager() *CaunterManage {
	return &CaunterManage{}
}

func (cm *CaunterManage) WithStorage(newStorage interfaces.StorageInterface) *CaunterManage {
	cm.mtx.Lock()
	defer cm.mtx.Unlock()
	cm.storage = newStorage
	return cm
}

func (cm *CaunterManage) UpdateCounter(name string, val int64) {
	cm.storage.UpdateCounter(name, val)
}

func (cm *CaunterManage) ShowAllCounters() map[string][]map[string]int64 {
	return cm.storage.GetCounterHistory()
}

func (cm *CaunterManage) GetActualCountersMetrics() map[string]int64 {
	return cm.storage.GetActualCounters()
}

func (cm *CaunterManage) GetMetricValueByName(name string) (int64, error) {
	return cm.storage.GetCountersByName(name)
}
