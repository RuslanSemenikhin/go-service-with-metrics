package env

import (
	"github.com/RuslanSemenikhin/go-service-with-metrics.git/internal/managers"
)

var (
	BOX *Box
)

type Box struct {
	gaugeManager   *managers.GaugeManage
	counterManager *managers.CaunterManage
}

func NewBox() *Box {
	return &Box{}
}

func (b *Box) WithGagugeManager(gm *managers.GaugeManage) *Box {
	b.gaugeManager = gm
	return b
}

func (b *Box) WithCaunterManager(cm *managers.CaunterManage) *Box {
	b.counterManager = cm
	return b
}

func (b *Box) GetGaugeManager() *managers.GaugeManage {
	return b.gaugeManager
}

func (b *Box) GetCounterManager() *managers.CaunterManage {
	return b.counterManager
}
