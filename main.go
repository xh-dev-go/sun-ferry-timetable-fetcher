package main

import (
	"context"
	"crypto/md5"
	"embed"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/eko/gocache/v3/store"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/dataFetch/cachedResult"
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/dataFetch/holiday"
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/service"
	"golang.org/x/exp/slices"
	"mime"
	"regexp"
	"sort"
	"strings"
	"time"
)

//go:embed static
var f embed.FS

var CacheManager = cachedResult.CacheC[[]service.FerryRecordDto]()

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

var todayETag = cachedResult.Cache[IsPublicHolidayDto]{}

const LAYOUT = "20060102"

const (
	DestinationCheungChau string = "cheung-chau"
	DestinationMuiWo             = "mui-wo"
)
const (
	DestExchangeCheungChau string = "Cheung Chau"
	DestExchangeMuiWo             = "Mui Wo"
)
const (
	directionFrom string = "from"
	directionTo          = "to"
)

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

			todayETag.Update(todayString+":"+hex.EncodeToString(md5Hash[:]), data)

			context.Header("ETag", todayETag.Key())
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
			if todayETag.Match(etagRequest) {
				data := todayETag.Value()
				c.Header("ETag", etagRequest)
				c.JSON(304, data)
				return
			}

			if len(holiday.IsPublicHoliday(todayString)) == 0 {
				setResponse(false, todayString, c)
			} else {
				setResponse(true, todayString, c)
			}
		})
		calendarGroup.GET("/today", func(c *gin.Context) {
			todayString := time.Now().Format(LAYOUT)
			c.Header("Access-Control-Expose-Headers", "ETag")
			etagRequest := c.Request.Header.Get("If-None-Match")

			if todayETag.Match(etagRequest) {
				data := todayETag.Value()
				c.Header("ETag", etagRequest)
				c.JSON(304, data)
				return
			}

			if len(holiday.IsPublicHoliday(todayString)) == 0 {
				setResponse(false, todayString, c)
			} else {
				setResponse(true, todayString, c)
			}
		})
	}
	sunFerryGroup := r.Group("/api/v1/ferry/sun-ferry")
	{
		sunFerryGroup.GET("service-map", func(c *gin.Context) {
			c.JSON(200, service.MapDict)
		})
		sunFerryGroup.GET("mui-wo", func(c *gin.Context) {
			extractGet(service.GetCentralToMuiWo, c)
		})
		sunFerryGroup.GET("/:direction/:from/:to/:time", func(c *gin.Context) {
			direction := c.Param("direction")
			from := c.Param("from")
			to := c.Param("to")

			timeParam := c.Param("time")
			if timeParam != "today" {
				match, err := regexp.MatchString("[0-9]{8}", timeParam)
				if err != nil || !match {
					panic("time not match")
				}
			}

			var dayString string
			if timeParam == "today" {
				dayString = time.Now().Format(LAYOUT)
			} else {
				date, err := time.Parse(LAYOUT, timeParam)
				if err != nil {
					panic("time not match")
				}
				dayString = date.Format(LAYOUT)
			}

			var ctx = context.Background()
			var cacheKey = fmt.Sprintf("/%s/%s/%s/%s", direction, from, to, dayString)

			MW := "mui-wo"
			CC := "cheung-chau"
			IS := "island"
			KC := "kowloon-city"
			HH := "hung-hom"
			var list []string
			list = append(list, MW)
			list = append(list, CC)
			list = append(list, IS)
			list = append(list, KC)
			list = append(list, HH)

			if !slices.Contains(list, direction) {
				panic("direction not match")
			}

			dict := service.MapDict[direction]

			var allLocation []string
			for k, _ := range dict {
				allLocation = append(allLocation, strings.ReplaceAll(strings.ToLower(k), " ", "-"))
			}

			if !slices.Contains(allLocation, from) || !slices.Contains(allLocation, to) {
				panic(fmt.Sprintf("From[%s] or To[%s] not match", from, to))
			}

			var dtos []service.FerryRecordDto
			var tag string
			if direction == MW {
				dtos, tag, _ = service.GetCentralToMuiWo()
			} else if direction == CC {
				dtos, tag, _ = service.GetCentralToCheungChau()
			} else if direction == IS {
				dtos, tag, _ = service.GetInterIsland()
			} else if direction == KC {
				dtos, tag, _ = service.GetNorthPointKowloonCity()
			} else if direction == HH {
				dtos, tag, _ = service.GetNorthPointHungHom()
			} else {
				panic(fmt.Sprintf("No match route found: %s", direction))
			}

			v, err := CacheManager.Get(ctx, cacheKey)

			etag := c.GetHeader("If-None-Match")

			if err == nil {
				if etag != "" && etag == tag && etag == v.Key() {
					c.Status(304)
					return
				}

			} else if err.Error() != store.NOT_FOUND_ERR {
				panic(err)
			}

			todayDate, err := time.Parse(LAYOUT, dayString)
			if err != nil {
				panic(err)
			}

			bFlag := holiday.TodayHolidayFlag(todayDate)

			var filtered []service.FerryRecordDto

			for _, dto := range dtos {
				if strings.ReplaceAll(strings.ToLower(dto.From), " ", "-") != from {
					continue
				}
				if strings.ReplaceAll(strings.ToLower(dto.To), " ", "-") != to {
					continue
				}
				flag := dto.GetFlag()
				if flag.AnyMatch(*bFlag) {
					filtered = append(filtered, dto)
				}
			}

			sort.SliceStable(filtered, func(i, j int) bool {
				return filtered[i].Time < filtered[j].Time
			})

			var cc cachedResult.Cache[[]service.FerryRecordDto]
			cc.Update(tag, filtered)

			CacheManager.Set(context.Background(), cacheKey, cc)
			c.Header("ETag", tag)
			c.JSON(200, filtered)
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
