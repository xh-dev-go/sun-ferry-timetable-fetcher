package service

import (
	"errors"
	"fmt"
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/dataFetch"
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/dataFetch/ferry"
	"github.com/xh-dev-go/xhUtils/binaryFlag"
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

type ETagCache[T any] struct {
	Value *[]T
	ETag  string ""
}

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
	return get(&northPointKowloonCityFerry, northPointKowloonCityConvert, northPointKowloonCityFerryETag)
}
func GetNorthPointHungHom() ([]FerryRecordDto, string, int) {
	return get(&hungHomNorthPointFerry, homeHomNorthPointConvert, hungHomNorthPointFerryETag)
}
func GetInterIsland() ([]FerryRecordDto, string, int) {
	return get(&interIslandFerry, interIslandConvert, interIslandFerryETag)
}
func GetCentralToCheungChau() ([]FerryRecordDto, string, int) {
	return get(&centralCheungChauCache, cheungChauConvert, centralCheungChauCacheETag)
}
func GetCentralToMuiWo() ([]FerryRecordDto, string, int) {
	return get(&centralMuiWoCache, muiwoConvert, centralMuiWoCacheETag)
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

func get(cache *SunFerryConfig, convert ferry.Convert, tagCache ETagCache[FerryRecordDto]) ([]FerryRecordDto, string, int) {
	str, status, eTagValue := ferry.Extract(cache.url, tagCache.ETag)

	if status == 200 {
		fmt.Println("Load new records")
		tagCache.ETag = eTagValue

		if cache.DecodeMode == DecodeMode1 {
			tagCache.Value = convertToDto(*ferry.Decode(str, cache.routeName, cache.dict, convert),
				cache.speedSetup,
				cache.frequencySetup,
				cache.remarkSetup,
				cache.zhRemarkSetup,
			)
		} else if cache.DecodeMode == DecodeMode2 {
			tagCache.Value = convertToDto(*ferry.DecodeIsland(str, cache.routeName, cache.dict, convert),
				cache.speedSetup,
				cache.frequencySetup,
				cache.remarkSetup,
				cache.zhRemarkSetup,
			)
		}
		return *tagCache.Value, tagCache.ETag, 200
	} else if status == 304 {
		fmt.Println(fmt.Sprintf("[%s]Load records from cache", cache.routeName))
		return nil, tagCache.ETag, 304
	} else {
		panic(errors.New("should not occurred"))
	}
}
