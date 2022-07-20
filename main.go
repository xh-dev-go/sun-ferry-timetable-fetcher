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
	r.GET("/ferry/sun-ferry/mui-wo", func(c *gin.Context) {
		c.JSON(200, service.GetCentralToMuiWo())
	})
	r.GET("/ferry/sun-ferry/cheung-chau", func(c *gin.Context) {
		c.JSON(200, service.GetCentralToCheungChau())
	})
	r.GET("/ferry/sun-ferry/inter-island", func(c *gin.Context) {
		c.JSON(200, service.GetInterIsland())
	})
	r.GET("/ferry/sun-ferry/hung-hom-north-point", func(c *gin.Context) {
		c.JSON(200, service.GetNorthPointHungHome())
	})
	r.GET("/ferry/sun-ferry/north-point-kowloon-city", func(c *gin.Context) {
		c.JSON(200, service.GetNorthPointKowloonCity())
	})
	r.Run()
}
