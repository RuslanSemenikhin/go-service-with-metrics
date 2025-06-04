package server

import (
	serv_http "github.com/RuslanSemenikhin/go-service-with-metrics.git/cmd/server/http"
	"github.com/RuslanSemenikhin/go-service-with-metrics.git/internal/env"
	"github.com/RuslanSemenikhin/go-service-with-metrics.git/internal/managers"
	"github.com/RuslanSemenikhin/go-service-with-metrics.git/internal/storage"
)

var SRV *serv_http.Srv

func init() {
	db := storage.NewStorage()

	gaugeManager := managers.NewGaugeManager().WithStorage(db)
	counterManager := managers.NewCaunterManager().WithStorage(db)
	box := env.NewBox().
		WithGagugeManager(gaugeManager).
		WithCaunterManager(counterManager)
	SRV = serv_http.NewSrv().WithBox(box)
}

func Start() {
	if err := serv_http.StartServer(`:8080`, SRV); err != nil {
		panic(err)
	}
}
