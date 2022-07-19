package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/service"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"sss": "ssss",
		})
	})
	r.GET("/mui-wo", func(c *gin.Context) {
		c.JSON(200, service.GetCentralToMuiWo())
	})
	r.GET("/cheung-chau", func(c *gin.Context) {
		c.JSON(200, service.GetCentralToCheungChau())
	})
	r.Run()
}
