package middleware

import (
	"log"
	"net/http"

	g "github.com/gin-gonic/gin"
)

func MiddlewareRecovery(ctx *g.Context, recovered interface{}) {
	log.Printf("===>>> panic ocured wile processing - \n\t==>> URL path: %s\n\t==>> request method: %s\n\t==>> trouble: %v",
		ctx.Request.URL.Path,
		ctx.Request.Method,
		recovered,
	)
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, g.H{
		"error": "server error",
	})
}
