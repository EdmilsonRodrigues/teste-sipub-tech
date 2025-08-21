package middlewares

import (
	"github.com/gin-gonic/gin"
	
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/core/ports"
)

func AddMovieQueryService(service ports.MovieQueryService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(ports.ServiceKey, service)
		ctx.Next()
	}
}

func AddMovieExecutorService(service ports.MovieExecutorService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(ports.ServiceKey, service)
		ctx.Next()
	}
}
