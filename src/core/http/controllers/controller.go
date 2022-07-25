package controllers

import (
	"log"
	"main/src/core/config"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Config *config.Config
}

func (controller *Controller) SendError(c *gin.Context, e error, code int) {
	if code == 500 && !controller.Config.IsDebug() {
		log.Printf("controller.error: %s", e.Error())
		c.AbortWithStatusJSON(code, gin.H{
			"message": "Server fatal error",
		})
	} else {
		c.AbortWithStatusJSON(code, gin.H{
			"message": e.Error(),
		})
	}
}

func (controller *Controller) GetUserIP(c *gin.Context) string {
	ip := c.GetHeader("CF-Connecting-IP")
	if len(ip) == 0 {
		ip = c.ClientIP()
	}

	return ip
}
