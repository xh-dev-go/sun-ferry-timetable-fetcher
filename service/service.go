package service

import (
	"errors"
	"fmt"
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/dataFetch"
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

type ETagCache struct {
	Value          *[]FerryRecordDto
	ETag           string ""
	dict           map[string]string
	url            string
	routeName      string
	speedSetup     binaryFlag.ValuePair[string]
	frequencySetup binaryFlag.ValuePair[string]
	remarkSetup    binaryFlag.ValuePair[string]
	zhRemarkSetup  binaryFlag.ValuePair[string]
	DecodeMode     string
}

func GetNorthPointKowloonCity() []FerryRecordDto {
	return get(&northPointKowloonCityFerry, northPointKowloonCityConvert)
}
func GetNorthPointHungHome() []FerryRecordDto {
	return get(&hungHomNorthPointFerry, homeHomNorthPointConvert)
}
func GetInterIsland() []FerryRecordDto {
	return get(&interIslandFerry, interIslandConvert)
}
func GetCentralToCheungChau() []FerryRecordDto {
	return get(&centralCheungChauCache, cheungChauConvert)
}
func GetCentralToMuiWo() []FerryRecordDto {
	return get(&centralMuiWoCache, muiwoConvert)
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

func get(cache *ETagCache, convert dataFetch.Convert) []FerryRecordDto {
	str, status, eTagValue := dataFetch.Extract(cache.url, cache.ETag)

	if status == 200 {
		fmt.Println("Load new records")
		cache.ETag = eTagValue

		if cache.DecodeMode == DecodeMode1 {
			cache.Value = convertToDto(*dataFetch.Decode(str, cache.routeName, cache.dict, convert),
				cache.speedSetup,
				cache.frequencySetup,
				cache.remarkSetup,
				cache.zhRemarkSetup,
			)
		} else if cache.DecodeMode == DecodeMode2 {
			cache.Value = convertToDto(*dataFetch.DecodeIsland(str, cache.routeName, cache.dict, convert),
				cache.speedSetup,
				cache.frequencySetup,
				cache.remarkSetup,
				cache.zhRemarkSetup,
			)
		}
		return *cache.Value
	} else if status == 304 {
		fmt.Println(fmt.Sprintf("[%s]Load records from cache", cache.routeName))
		return *cache.Value
	} else {
		panic(errors.New("should not occurred"))
	}
}
