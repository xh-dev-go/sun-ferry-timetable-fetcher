package service

import (
	"fmt"
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/dataFetch"
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/dataFetch/ferry"
	"github.com/xh-dev-go/xhUtils/binaryFlag"
)

var centralMuiWoCacheETag = ETagCache[FerryRecordDto]{}
var centralMuiWoCache = SunFerryConfig{
	DecodeMode: DecodeMode1,
	routeName:  "Central - Mui Wo",
	url:        "https://www.sunferry.com.hk/eta/timetable/SunFerry_central_muiwo_timetable_eng.csv",
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

var muiwoConvert = ferry.Convert{
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
