package service

import (
	"errors"
	"fmt"
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/dataFetch"
)

type FerryRecordDto struct {
	Route     string
	From      string
	ZhFrom    string
	To        string
	ZhTo      string
	Frequency []string
	Time      int
	Speed     string
	Remark    []string
}

type ETagCache struct {
	Value     *[]FerryRecordDto
	ETag      string ""
	dict      map[string]string
	url       string
	routeName string
}

var centralMuiWoCache = ETagCache{
	routeName: "Central - Mui Wo",
	url:       "https://www.sunferry.com.hk/eta/timetable/SunFerry_central_muiwo_timetable_eng.csv",
	dict: map[string]string{
		"Mui Wo":  "梅窩",
		"Central": "中環",
	},
}
var centralCheungChauCache = ETagCache{
	routeName: "Central - Cheung Chau",
	url:       "https://www.sunferry.com.hk/eta/timetable/SunFerry_central_cheungchau_timetable_eng.csv",
	dict: map[string]string{
		"Cheung Chau": "長洲",
		"Central":     "中環",
	},
}

func GetCentralToCheungChau() []dataFetch.FerryRecord {
	return get(&centralCheungChauCache)
}
func GetCentralToMuiWo() []dataFetch.FerryRecord {
	return get(&centralMuiWoCache)
}

func convertToDto(arr []dataFetch.FerryRecord) {
	var list []FerryRecordDto
	for _, item := range arr {
		var speed string
		if item.Speed.IsSet(dataFetch.SpeedFast) {
			speed = "Fast"
		} else {
			speed = "Ordinary"
		}

		list = append(list, FerryRecordDto{
			Route:  item.Route,
			From:   item.From,
			To:     item.To,
			ZhFrom: item.ZhFrom,
			ZhTo:   item.ZhTo,
			Time:   item.Time,
			Speed:  speed,
		})
	}
	return list
}

func get(cache *ETagCache) []dataFetch.FerryRecord {
	str, status, eTagValue := dataFetch.Extract(cache.url, cache.ETag)

	if status == 200 {
		fmt.Println("Load new records")
		cache.ETag = eTagValue
		cache.Value = dataFetch.Decode(str, cache.routeName, cache.dict)
		return *cache.Value
	} else if status == 304 {
		fmt.Println(fmt.Sprintf("[%s]Load records from cache", cache.routeName))
		return *cache.Value
	} else {
		panic(errors.New("should not occurred"))
	}
}
