package holiday

import (
	"encoding/json"
	"errors"
	"fmt"
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

func ExtractJson(URL string, eTag string) (string, int, string) {
	client := http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		panic(err)
	}
	if eTag != "" {
		req.Header.Set("If-None-Match", eTag)
	}
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if response.StatusCode == 200 {
		msgByte, err := io.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		msg := string(msgByte)
		return msg, 200, response.Header.Get("ETag")

	} else if response.StatusCode == 304 {
		return "", 304, ""

	} else {
		panic(errors.New(fmt.Sprintf("Status[%d] not support", response.StatusCode)))
	}

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
