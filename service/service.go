package service

import (
	"fmt"
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/dataFetch"
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/dataFetch/cachedResult"
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/dataFetch/ferry"
	"github.com/xh-dev-go/xhUtils/binaryFlag"
	"io"
	"net/http"
)

// https://www.sunferry.com.hk/eta/SunFerry_time_table_and_fare_table_dataspec_eng.pdf

const DecodeMode1 = "mode1"
const DecodeMode2 = "mode2"

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
	ZhRemark  []string
}

//type ETagCache[T any] struct {
//	Value *[]T
//	ETag  string ""
//}

type SunFerryConfig struct {
	dict           map[string]string
	url            string
	routeName      string
	speedSetup     binaryFlag.ValuePair[string]
	frequencySetup binaryFlag.ValuePair[string]
	remarkSetup    binaryFlag.ValuePair[string]
	zhRemarkSetup  binaryFlag.ValuePair[string]
	DecodeMode     string
}

func GetNorthPointKowloonCity() ([]FerryRecordDto, string, int) {
	return get(&northPointKowloonCityFerry, northPointKowloonCityConvert, &northPointKowloonCityFerryETag)
}
func GetNorthPointHungHom() ([]FerryRecordDto, string, int) {
	return get(&hungHomNorthPointFerry, homeHomNorthPointConvert, &hungHomNorthPointFerryETag)
}
func GetInterIsland() ([]FerryRecordDto, string, int) {
	return get(&interIslandFerry, interIslandConvert, &interIslandFerryETag)
}
func GetCentralToCheungChau() ([]FerryRecordDto, string, int) {
	return get(&centralCheungChauCache, cheungChauConvert, &centralCheungChauCacheETag)
}
func GetCentralToMuiWo() ([]FerryRecordDto, string, int) {
	return get(&centralMuiWoCache, muiwoConvert, &centralMuiWoCacheETag)
}

func convertToDto(
	arr []dataFetch.FerryRecord,
	speedValuePair binaryFlag.ValuePair[string],
	frequencyValuePair binaryFlag.ValuePair[string],
	remarkValuePair binaryFlag.ValuePair[string],
	zhRemarkValuePair binaryFlag.ValuePair[string],
) *[]FerryRecordDto {
	var list []FerryRecordDto

	for _, item := range arr {
		speed, err := speedValuePair.ExtractAny(item.Speed)
		if err != nil {
			panic(err)
		}
		listOfFrequency := frequencyValuePair.ExtractAll(item.Frequency)
		listOfRemark := remarkValuePair.ExtractAll(item.Remark)
		listOfZhRemark := zhRemarkValuePair.ExtractAll(item.Remark)

		list = append(list, FerryRecordDto{
			Route:     item.Route,
			From:      item.From,
			To:        item.To,
			ZhFrom:    item.ZhFrom,
			ZhTo:      item.ZhTo,
			Time:      item.Time,
			Speed:     speed,
			Frequency: listOfFrequency,
			Remark:    listOfRemark,
			ZhRemark:  listOfZhRemark,
		})
	}
	return &list
}

func get(cache *SunFerryConfig, convert ferry.Convert, tagCache *cachedResult.Cache[[]FerryRecordDto]) ([]FerryRecordDto, string, int) {
	client := http.Client{}
	req, err := http.NewRequest("GET", cache.url, nil)
	if err != nil {
		panic(err)
	}
	result := tagCache.HttpCaching(req, &client, func(response http.Response) ([]FerryRecordDto, error) {
		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		if cache.DecodeMode == DecodeMode1 {
			return *convertToDto(*ferry.Decode(string(bytes), cache.routeName, cache.dict, convert),
				cache.speedSetup,
				cache.frequencySetup,
				cache.remarkSetup,
				cache.zhRemarkSetup,
			), nil
		} else if cache.DecodeMode == DecodeMode2 {
			return *convertToDto(*ferry.DecodeIsland(string(bytes), cache.routeName, cache.dict, convert),
				cache.speedSetup,
				cache.frequencySetup,
				cache.remarkSetup,
				cache.zhRemarkSetup,
			), nil
		} else {
			panic(fmt.Sprintf("Unmatch decode mode: %s", cache.DecodeMode))
		}
	})

	return result.Cache().Value(), result.Cache().Key(), result.Response().StatusCode

}
