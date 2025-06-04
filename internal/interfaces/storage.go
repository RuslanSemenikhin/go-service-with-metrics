package interfaces

type StorageInterface interface {
	IsEmptyGauge() bool
	IsEmptyCounter() bool
	GetGaugesByName(string) (float64, error)
	GetCountersByName(string) (int64, error)
	UpdateGauge(string, float64)
	UpdateCounter(string, int64)
	GetCounterHistory() map[string][]map[string]int64
	GetGaugeHistory() map[string][]float64
}
