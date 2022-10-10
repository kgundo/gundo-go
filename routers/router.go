package routers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup(ctx context.Context) *gin.Engine {
	r := gin.Default()

	r.GET("ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong") //ping-pong lol
	})

	v1 := r.Group("/api/v1")
	{
		UserRouter(v1)
	}
	return r
}
