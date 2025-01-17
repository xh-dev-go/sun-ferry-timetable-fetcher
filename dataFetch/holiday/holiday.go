package holiday

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/dataFetch/cachedResult"
	"github.com/xh-dev-go/xhUtils/binaryFlag"
	"io"
	"net/http"
	"strings"
	"time"
)

type Holiday struct {
	Uid  string
	Date string
	Name string
}

type Calendar struct {
	VCalendar []VCalendar `json:"vcalendar"`
}
type VCalendar struct {
	Prodid      string   `json:"prodid"`
	Version     string   `json:"version"`
	CalScale    string   `json:"calscale"`
	XwrTimezone string   `json:"x-wr-timezone"`
	XWrCalname  string   `json:"x-wr-calname"`
	XWrCalDes   string   `json:"x-wr-caldesc"`
	VEvent      []VEvent `json:"vevent"`
}

type VEvent struct {
	DtStart [2]interface{} `json:"dtstart"`
	DtEnd   [2]interface{} `json:"dtend"`
	Transp  string         `json:"transp"`
	Uid     string         `json:"uid"`
	Summary string         `json:"summary"`
}

var CachedHolidayApi = cachedResult.Cache[[]Holiday]{}

const layout = "20060102"

func TodayHolidayFlag(todayDate time.Time) *binaryFlag.BinaryFlag {
	dayStr := todayDate.Format(layout)
	holidays := IsPublicHoliday(dayStr)
	bFlag := binaryFlag.New()

	switch todayDate.Weekday() {
	case time.Monday:
		bFlag.SetBit(1)
	case time.Tuesday:
		bFlag.SetBit(2)
	case time.Wednesday:
		bFlag.SetBit(3)
	case time.Thursday:
		bFlag.SetBit(4)
	case time.Friday:
		bFlag.SetBit(5)
	case time.Saturday:
		bFlag.SetBit(6)
	case time.Sunday:
		bFlag.SetBit(7)
	}

	if len(holidays) > 0 {
		bFlag.SetBit(10)
	}

	return bFlag
}

func IsPublicHoliday(dateString string) []Holiday {
	var arr []Holiday
	date, err := time.Parse(layout, dateString)
	if err != nil {
		panic(err)
	}
	if date.Weekday() == time.Sunday {
		arr = append(arr, Holiday{
			Uid:  "",
			Date: dateString,
			Name: time.Sunday.String(),
		})
	}

	holidaysResult := GetHolidays()
	if holidaysResult.HasError() {
		panic(holidaysResult.Error())
	} else if !holidaysResult.IsResultCached() {
		panic(fmt.Sprintf("Not expected result [%s]", holidaysResult.Response().Status))
	}
	for _, item := range holidaysResult.Cache().Value() {
		if item.Date == dateString {
			arr = append(arr, item)
		}
	}
	return arr
}

func GetHolidays() cachedResult.CacheResult[[]Holiday] {
	url := "https://www.1823.gov.hk/common/ical/en.json"
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	return CachedHolidayApi.HttpCaching(req, &client, func(response http.Response) ([]Holiday, error) {
		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		} else {
			return *DecodeHoliday(string(bytes)), nil
		}
	})

}
func DecodeHoliday(msg string) *[]Holiday {
	data := strings.TrimPrefix(msg, "\xef\xbb\xbf")
	var cal = Calendar{}
	err := json.Unmarshal([]byte(data), &cal)
	if err != nil {
		panic(err)
	}

	var holidays []Holiday

	layout := "20060102"

	for _, v := range cal.VCalendar[0].VEvent {
		first := v.DtStart[1].(map[string]interface{})
		second := v.DtEnd[1].(map[string]interface{})
		if first["value"] != "DATE" || second["value"] != "DATE" {
			panic(errors.New("value should be Date"))
		}

		firstStr := v.DtStart[0].(string)
		secStr := v.DtEnd[0].(string)

		var curDateStr = firstStr
		for curDateStr != secStr {
			holidays = append(holidays, Holiday{
				Date: curDateStr,
				Uid:  v.Uid,
				Name: v.Summary,
			})

			tempDate, err := time.Parse(layout, firstStr)
			if err != nil {
				panic(err)
			}
			curDateStr = tempDate.AddDate(0, 0, 1).Format(layout)
		}

	}

	return &holidays
}
