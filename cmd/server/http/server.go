package http

import (
	"net/http"
	"sync"

	"github.com/RuslanSemenikhin/go-service-with-metrics.git/internal/env"
	"github.com/RuslanSemenikhin/go-service-with-metrics.git/internal/handlefunc"
)

type Srv struct {
	mtx    sync.RWMutex
	srvMax *http.ServeMux
	box    *env.Box
}

func NewSrv() *Srv {
	newSrv := http.NewServeMux()
	return &Srv{
		srvMax: newSrv,
	}
}

func StartServer(port string, srv *Srv) error {
	s := srv.GetSrv()
	srv.InitializeRoutes()
	err := http.ListenAndServe(port, s.srvMax)
	return err
}

func (s *Srv) WithBox(b *env.Box) *Srv {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.box = b
	return s
}

func (s *Srv) InitializeRoutes() {
	s.srvMax.HandleFunc(`/update/`, handlefunc.Update(s.box))

}

func (s *Srv) GetSrv() *Srv {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s
}
