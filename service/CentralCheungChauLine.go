package service

import (
	"errors"
	"fmt"
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/dataFetch"
	"github.com/xh-dev-go/xhUtils/binaryFlag"
)

var centralCheungChauCache = ETagCache{
	DecodeMode: DecodeMode1,
	routeName:  "Central - Cheung Chau",
	url:        "https://www.sunferry.com.hk/eta/timetable/SunFerry_central_cheungchau_timetable_eng.csv",
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
