package service

import (
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/dataFetch"
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/dataFetch/cachedResult"
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/dataFetch/ferry"
	"github.com/xh-dev-go/xhUtils/binaryFlag"
)

var hungHomNorthPointFerryETag = cachedResult.Cache[[]FerryRecordDto]{}
var hungHomNorthPointFerry = SunFerryConfig{
	DecodeMode: DecodeMode2,
	routeName:  "hung home - north point ferry",
	url:        "https://www.sunferry.com.hk/eta/timetable/SunFerry_northpoint_hunghom_timetable_eng.csv",
	dict: map[string]string{
		"North Point": "北角",
		"Hung Hom":    "紅磡",
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
var homeHomNorthPointConvert = ferry.Convert{
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
