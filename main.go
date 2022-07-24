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
	"regexp"
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
	calendarGroup := r.Group("/api/v1/calendar/public-holiday")
	{
		layout := "20060102"
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
			todayETag.ETag = todayString + ":" + hex.EncodeToString(md5Hash[:])

			context.Header("ETag", todayETag.ETag)
			context.JSON(200, data)
		}
		calendarGroup.GET("", func(c *gin.Context) {
			result := holiday.GetHolidays()
			c.Header("Access-Control-Expose-Headers", "ETag")
			if result.HasError() {
				c.Status(result.Response().StatusCode)
			} else if result.IsResultCached() {
				c.Header("ETag", result.Cache().Key())
				c.JSON(200, result.Cache().Value())
			} else {
				c.Status(result.Response().StatusCode)
			}
		})
		calendarGroup.GET("/date/:dateStr", func(c *gin.Context) {
			dateStr := c.Param("dateStr")
			if dateStr == "" {
				c.Status(400)
				return
			}
			matching, err := regexp.MatchString("^[0-9]{8}$", dateStr)

			if err != nil || !matching {
				c.Status(400)
				return
			}

			layout := "20060102"
			date, err := time.Parse(layout, dateStr)
			if err != nil || !matching {
				c.Status(400)
				return
			}
			todayString := date.Format(layout)
			c.Header("Access-Control-Expose-Headers", "ETag")
			etagRequest := c.Request.Header.Get("If-None-Match")
			if etagRequest != "" && todayString+":"+etagRequest == todayETag.ETag {
				data := *todayETag.Value
				c.Header("ETag", etagRequest)
				c.JSON(304, data[0])
				return
			}

			if holiday.IsPublicHoliday(todayString) == nil {
				setResponse(false, todayString, c)
			} else {
				setResponse(true, todayString, c)
			}
		})
		calendarGroup.GET("/today", func(c *gin.Context) {
			todayString := time.Now().Format(layout)
			c.Header("Access-Control-Expose-Headers", "ETag")
			etagRequest := c.Request.Header.Get("If-None-Match")
			if etagRequest != "" && todayString+":"+etagRequest == todayETag.ETag {
				data := *todayETag.Value
				c.Header("ETag", etagRequest)
				c.JSON(304, data[0])
				return
			}

			if holiday.IsPublicHoliday(todayString) == nil {
				setResponse(false, todayString, c)
			} else {
				setResponse(true, todayString, c)
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
