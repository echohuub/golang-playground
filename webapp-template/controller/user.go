package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"webapp.demo/model"
	"webapp.demo/service"
)

func SignUpHandler(c *gin.Context) {
	var p model.ParamSignUp
	if err := c.ShouldBind(&p); err != nil {
		zap.L().Error("", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"msg": "参数错误",
		})
		return
	}

	service.SignUp(p)
	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}
