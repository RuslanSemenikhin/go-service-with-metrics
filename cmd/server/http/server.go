package http

import (
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/RuslanSemenikhin/go-service-with-metrics.git/internal/env"
	"github.com/RuslanSemenikhin/go-service-with-metrics.git/internal/handlefunc"
	mw "github.com/RuslanSemenikhin/go-service-with-metrics.git/internal/middleware"
	g "github.com/gin-gonic/gin"
)

type Srv struct {
	mtx sync.RWMutex
	srv *g.Engine
	box *env.Box
}

func NewSrv() *Srv {
	return &Srv{
		srv: g.New(),
	}
}

func (s *Srv) StartServer(port string) error {
	s.RegistrateRoutes()
	if err := s.srv.Run(port); err != nil {
		return err
	}
	return nil
}

func (s *Srv) WithBox(b *env.Box) *Srv {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.box = b
	return s
}

func (s *Srv) RegistrateRoutes() {
	s.srv.Use(g.CustomRecovery(mw.MiddlewareRecovery))
	s.srv.NoMethod(func(ctx *g.Context) {
		ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, g.H{
			"error": fmt.Sprintf("method not allowed, incoming method - '%s'", ctx.Request.Method),
		})
	})

	s.srv.LoadHTMLGlob(s.GetPathToTemplates())

	s.srv.GET(`/`, handlefunc.GetAllMetrics(s.box))
	s.srv.GET(`/value/:metricType/:metricName`, handlefunc.GetMetricValueByName(s.box))

	s.srv.POST(`/update/:metricType/:metricName/:metricValue`, handlefunc.Update(s.box))
}

func (s *Srv) GetPathToTemplates() string {
	_, filename, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(filename)
	tmplPath := filepath.Join(basePath, "..", "..", "..", "internal", "templates", "*.tmpl")
	return tmplPath
}

func (s *Srv) GetSrv() *Srv {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s
}
