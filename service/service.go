package service

import (
	"errors"
	"fmt"
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/dataFetch"
	"github.com/xh-dev-go/xhUtils/binaryFlag"
)

// https://www.sunferry.com.hk/eta/SunFerry_time_table_and_fare_table_dataspec_eng.pdf

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
}

var centralMuiWoCache = ETagCache{
	routeName: "Central - Mui Wo",
	url:       "https://www.sunferry.com.hk/eta/timetable/SunFerry_central_muiwo_timetable_eng.csv",
	dict: map[string]string{
		"Mui Wo":  "梅窩",
		"Central": "中環",
	},
	speedSetup: binaryFlag.ValuePair[string]{
		Values: map[int]string{
			dataFetch.SpeedFast:     "Fast",
			dataFetch.SpeedOrdinary: "Ordinary",
		},
	},
	frequencySetup: binaryFlag.ValuePair[string]{
		Values: map[int]string{
			1:  "Monday",
			2:  "Tuesday",
			3:  "Wednesday",
			4:  "Thursday",
			5:  "Friday",
			6:  "Saturday",
			7:  "Sunday",
			10: "Public Holiday",
		},
	},
	remarkSetup: binaryFlag.ValuePair[string]{
		Values: map[int]string{
			dataFetch.ViaPengChau: "Via Peng Chau",
		},
	},
	zhRemarkSetup: binaryFlag.ValuePair[string]{
		Values: map[int]string{
			dataFetch.ViaPengChau: "經坪洲",
		},
	},
}

var muiwoConvert = dataFetch.Convert{
	ToSpeed: func(serviceDate string, remark string) binaryFlag.BinaryFlag {
		if remark == "1" || remark == "2" {
			return *binaryFlag.New().SetBit(dataFetch.SpeedOrdinary)
		} else {
			return *binaryFlag.New().SetBit(dataFetch.SpeedFast)
		}
	},
	ToFrequency: func(serviceDate string, remark string) binaryFlag.BinaryFlag {
		mtf := *binaryFlag.New().SetBit(1).SetBit(2).SetBit(3).SetBit(4).SetBit(5)
		sat := *binaryFlag.New().SetBit(6)
		sunPub := *binaryFlag.New().SetBit(7).SetBit(10)

		if remark == "3" {
			return sat
		} else if serviceDate == "Mondays to Fridays except public holidays" {
			return mtf
		} else if serviceDate == "Saturdays except public holidays" {
			return sat
		} else if serviceDate == "Sundays and public holidays" {
			return sunPub
		} else {
			panic(fmt.Sprintf("service date[%s] and remark[%s] not working", serviceDate, remark))
		}
	},
	ToRemark: func(serviceDate string, remark string) binaryFlag.BinaryFlag {
		fast := *binaryFlag.New().SetBit(dataFetch.SpeedFast)
		ordinary := *binaryFlag.New().SetBit(dataFetch.SpeedOrdinary)
		if remark == "1" || remark == "2" {
			return ordinary
		} else {
			return fast
		}
	},
}

var centralCheungChauCache = ETagCache{
	routeName: "Central - Cheung Chau",
	url:       "https://www.sunferry.com.hk/eta/timetable/SunFerry_central_cheungchau_timetable_eng.csv",
	dict: map[string]string{
		"Cheung Chau": "長洲",
		"Central":     "中環",
	},
	speedSetup: binaryFlag.ValuePair[string]{
		Values: map[int]string{
			dataFetch.SpeedFast:     "Fast",
			dataFetch.SpeedOrdinary: "Ordinary",
		},
	},
	frequencySetup: binaryFlag.ValuePair[string]{
		Values: map[int]string{
			1:  "Monday",
			2:  "Tuesday",
			3:  "Wednesday",
			4:  "Thursday",
			5:  "Friday",
			6:  "Saturday",
			7:  "Sunday",
			10: "Public Holiday",
		},
	},
	remarkSetup: binaryFlag.ValuePair[string]{
		Values: map[int]string{},
	},
	zhRemarkSetup: binaryFlag.ValuePair[string]{
		Values: map[int]string{},
	},
}

var cheungChauConvert = dataFetch.Convert{
	ToSpeed: func(serviceDate string, remark string) binaryFlag.BinaryFlag {
		if remark == "1" || remark == "3" || remark == "4" {
			return *binaryFlag.New().SetBit(dataFetch.SpeedOrdinary)
		} else {
			return *binaryFlag.New().SetBit(dataFetch.SpeedFast)
		}
	},
	ToFrequency: func(serviceDate string, remark string) binaryFlag.BinaryFlag {
		mtf := *binaryFlag.New().SetBit(1).SetBit(2).SetBit(3).SetBit(4).SetBit(5)
		sunPub := *binaryFlag.New().SetBit(7).SetBit(10)

		mts := *binaryFlag.New().SetBit(1).SetBit(2).SetBit(3).SetBit(4).SetBit(5).SetBit(6).SetBit(7)
		sat := *binaryFlag.New().SetBit(6)
		if serviceDate == "Mondays to Fridays except public holidays" && remark == "2" {
			return *binaryFlag.New().SetBinary(mtf)
		} else if serviceDate == "Sundays and public holidays" {
			return *binaryFlag.New().SetBinary(sunPub)
		} else if serviceDate == "Saturdays except public holidays" {
			return *binaryFlag.New().SetBinary(sat)
		} else if serviceDate == "Mondays to Saturdays except public holidays" {
			return *binaryFlag.New().SetBinary(mts)
		} else {
			panic(errors.New(fmt.Sprintf("Line[%s] not support", serviceDate)))
		}
	},
	ToRemark: func(serviceDate string, remark string) binaryFlag.BinaryFlag {
		if remark == "3" {
			return *binaryFlag.New().SetBit(dataFetch.ViaPengChau)
		} else {
			return *binaryFlag.New()
		}
	},
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
		cache.Value = convertToDto(*dataFetch.Decode(str, cache.routeName, cache.dict, convert),
			cache.speedSetup,
			cache.frequencySetup,
			cache.remarkSetup,
			cache.zhRemarkSetup,
		)
		return *cache.Value
	} else if status == 304 {
		fmt.Println(fmt.Sprintf("[%s]Load records from cache", cache.routeName))
		return *cache.Value
	} else {
		panic(errors.New("should not occurred"))
	}
}
