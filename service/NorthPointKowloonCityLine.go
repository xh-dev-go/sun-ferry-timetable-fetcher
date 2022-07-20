package service

import (
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/dataFetch"
	"github.com/xh-dev-go/xhUtils/binaryFlag"
)

var northPointKowloonCityFerry = ETagCache{
	DecodeMode: DecodeMode2,
	routeName:  "north point - kowloon city ferry",
	url:        "https://www.sunferry.com.hk/eta/timetable/SunFerry_northpoint_kowlooncity_timetable_eng.csv",
	dict: map[string]string{
		"North Point":  "北角",
		"Kowloon City": "九龍城",
	},
	speedSetup: binaryFlag.ValuePair[string]{
		Values: map[int]string{
			dataFetch.SpeedOrdinary: "Ordinary",
		},
	},
	frequencySetup: binaryFlag.ValuePair[string]{
		Values: map[int]string{
			1: "Monday",
			2: "Tuesday",
			3: "Wednesday",
			4: "Thursday",
			5: "Friday",
			6: "Saturday",
			7: "Sunday",
		},
	},
	remarkSetup: binaryFlag.ValuePair[string]{
		Values: map[int]string{},
	},
	zhRemarkSetup: binaryFlag.ValuePair[string]{
		Values: map[int]string{},
	},
}
var northPointKowloonCityConvert = dataFetch.Convert{
	ToSpeed: func(serviceDate string, remark string) binaryFlag.BinaryFlag {
		return *binaryFlag.New().SetBit(dataFetch.SpeedOrdinary)
	},
	ToFrequency: func(serviceDate string, remark string) binaryFlag.BinaryFlag {
		if remark == "1" {
			return *binaryFlag.New().SetBit(1).SetBit(2).SetBit(3).SetBit(4).SetBit(5)
		} else {
			return *binaryFlag.New().SetBit(1).SetBit(2).SetBit(3).SetBit(4).SetBit(5).SetBit(6).SetBit(7)
		}
	},
	ToRemark: func(serviceDate string, remark string) binaryFlag.BinaryFlag {
		return *binaryFlag.New()
	},
}
