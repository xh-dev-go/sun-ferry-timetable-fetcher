package main

import (
	"crypto/md5"
	"embed"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/dataFetch/holiday"
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/service"
	"mime"
	"time"
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

type IsPublicHolidayDto struct {
	IsPublicHoliday bool   `json:"isPublicHoliday"`
	Summary         string `json:"summary"`
}

var todayETag = service.ETagCache[IsPublicHolidayDto]{}

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
			result := holiday.GetHolidays()
			c.Header("Access-Control-Expose-Headers", "ETag")
			if result.HasError() {
				c.Status(result.Response.StatusCode)
			} else if result.Cached {
				c.Header("ETag", result.ETag)
				c.JSON(200, result.Value)
			} else {
				c.Status(result.Response.StatusCode)
			}
		})
		calendarGroup.GET("/check/today", func(c *gin.Context) {
			c.Header("Access-Control-Expose-Headers", "ETag")
			etagRequest := c.Request.Header.Get("If-None-Match")
			if etagRequest != "" && etagRequest == todayETag.ETag {
				data := *todayETag.Value
				c.Header("ETag", etagRequest)
				c.JSON(304, data[0])
				return
			}

			holidaysResult := holiday.GetHolidays()

			var setResponse = func(status bool, todayString string, context *gin.Context) {
				data := IsPublicHolidayDto{
					IsPublicHoliday: status,
					Summary:         "",
				}
				dBytes, err := json.Marshal(data)
				if err != nil {
					panic(err)
				}
				var arr []IsPublicHolidayDto
				arr = append(arr, data)
				md5Hash := md5.Sum(dBytes)

				todayETag.Value = &arr
				todayETag.ETag = hex.EncodeToString(md5Hash[:])

				context.Header("ETag", todayETag.ETag)
				context.JSON(200, data)
			}

			layout := "20060102"
			todayString := time.Now().Format(layout)
			for _, item := range holidaysResult.Value {
				if item.Date == todayString {
					setResponse(true, todayString, c)
					return
				}
			}
			setResponse(false, todayString, c)
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
