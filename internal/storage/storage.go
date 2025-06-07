package storage

import (
	"errors"
	"fmt"
	"sync"

	"github.com/RuslanSemenikhin/go-service-with-metrics.git/internal/interfaces"
)

var _ interfaces.StorageInterface = (*Storage)(nil)

type Storage struct {
	mtx            sync.RWMutex
	gaugeHistory   map[string][]float64          // key - name metric
	counterHistory map[string][]map[string]int64 // key1 - name metric; slice with map metrics, key2 - must be 2 key -> 'counter' and 'value'
}

func NewStorage() *Storage {
	return &Storage{
		gaugeHistory:   make(map[string][]float64),
		counterHistory: make(map[string][]map[string]int64),
	}
}

func (s *Storage) IsEmptyGauge() bool {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return len(s.gaugeHistory) == 0
}

func (s *Storage) IsEmptyCounter() bool {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return len(s.counterHistory) == 0
}

func (s *Storage) GetGaugesByName(name string) (float64, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	if len(s.gaugeHistory) == 0 {
		errString := "field 'gauge' from storage is empty"
		err := errors.New(errString)
		return 0, err
	}
	if slcVals, ok := s.gaugeHistory[name]; !ok {
		errString := fmt.Sprintf("metric - '%s' with type metric - 'gauge' does not exists from storage", name)
		err := errors.New(errString)
		return 0, err
	} else {
		return slcVals[0], nil
	}
}

func (s *Storage) GetCountersByName(name string) (int64, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	if len(s.counterHistory) == 0 {
		errString := "field 'counter' from storage is empty"
		err := errors.New(errString)
		return 0, err
	}
	if slcVals, ok := s.counterHistory[name]; !ok {
		errString := fmt.Sprintf("metric - '%s' with type metric - 'counter' does not exists from storage", name)
		err := errors.New(errString)
		return 0, err
	} else {
		return slcVals[0]["counter"], nil
	}
}

func (s *Storage) UpdateGauge(name string, val float64) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.gaugeHistory[name] = append([]float64{val}, s.gaugeHistory[name]...)
}

func (s *Storage) UpdateCounter(name string, val int64) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if slcMaps, ok := s.counterHistory[name]; !ok {
		newMap := map[string]int64{
			"counter": val,
			"value":   val,
		}
		s.counterHistory[name] = []map[string]int64{newMap}
	} else {
		currentCounter := slcMaps[0]["counter"]
		newCounter := currentCounter + val
		newMap := map[string]int64{
			"counter": newCounter,
			"value":   val,
		}
		slc := append([]map[string]int64{newMap}, slcMaps...)
		s.counterHistory[name] = slc
	}
}

func (s *Storage) GetGaugeHistory() map[string][]float64 {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	copyHistory := make(map[string][]float64, len(s.gaugeHistory))

	for key, slc := range s.gaugeHistory {
		copySlc := make([]float64, len(slc))
		copy(copySlc, slc)
		copyHistory[key] = copySlc
	}
	return copyHistory
}

func (s *Storage) GetCounterHistory() map[string][]map[string]int64 {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	copyHistory := make(map[string][]map[string]int64, len(s.counterHistory))

	for key, slc := range s.counterHistory {
		copySlc := make([]map[string]int64, 0, len(slc))
		for _, mp := range slc {
			copyMp := make(map[string]int64)
			for k, v := range mp {
				copyMp[k] = v
			}
			copySlc = append(copySlc, copyMp)
		}
		copyHistory[key] = copySlc
	}
	return copyHistory
}
