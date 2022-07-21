package main

import (
	"embed"
	_ "embed"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/dataFetch/holiday"
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/service"
	"mime"
)

//go:embed static
var f embed.FS

func extractGet[T any](function func() ([]T, string, int), c *gin.Context) {
	arr, eTag, status := function()
	if status == 200 {
		c.Header("ETag", eTag)
		c.Header("Access-Control-Expose-Headers", "ETag")
		c.JSON(200, arr)
	} else if status == 304 {
		c.Status(304)
	} else {
		panic("error while return")
	}
}

func main() {
	err := mime.AddExtensionType(".js", "application/javascript")
	if err != nil {
		return
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))
	r.Use(static.Serve("/", static.LocalFile("./static", true)))
	calendarGroup := r.Group("/api/v1/calendar")
	{
		calendarGroup.GET("", func(c *gin.Context) {
			data, status, eTag := holiday.ExtractJson("https://www.1823.gov.hk/common/ical/en.json", "")
			holidays := holiday.DecodeHoliday(data)
			c.Header("Access-Control-Expose-Headers", "ETag")
			if status == 200 {
				c.Header("ETag", eTag)
				c.JSON(200, holidays)
			} else if status == 304 {
				c.Status(304)
			} else {
				panic("error while return")
			}

		})

	}
	sunFerryGroup := r.Group("/api/v1/ferry/sun-ferry")
	{
		sunFerryGroup.GET("mui-wo", func(c *gin.Context) {
			extractGet(service.GetCentralToMuiWo, c)
		})
		sunFerryGroup.GET("cheung-chau", func(c *gin.Context) {
			extractGet(service.GetCentralToCheungChau, c)
		})
		sunFerryGroup.GET("inter-island", func(c *gin.Context) {
			extractGet(service.GetInterIsland, c)
		})
		sunFerryGroup.GET("hung-hom-north-point", func(c *gin.Context) {
			extractGet(service.GetNorthPointHungHom, c)
		})
		sunFerryGroup.GET("north-point-kowloon-city", func(c *gin.Context) {
			extractGet(service.GetNorthPointKowloonCity, c)
		})
	}
	r.StaticFile("/", "./static")
	r.Run()
}
