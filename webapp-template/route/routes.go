package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webapp.demo/config"
	"webapp.demo/logger"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, config.Conf.Version)
	})
	return r
}
