package service

import (
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/dataFetch"
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/dataFetch/cachedResult"
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/dataFetch/ferry"
	"github.com/xh-dev-go/xhUtils/binaryFlag"
)

var interIslandFerryETag = cachedResult.Cache[[]FerryRecordDto]{}
var interIslandFerry = SunFerryConfig{
	DecodeMode: DecodeMode2,
	routeName:  "inter island ferry",
	url:        "https://www.sunferry.com.hk/eta/timetable/SunFerry_interislands_timetable_eng.csv",
	dict: map[string]string{
		"Cheung Chau": "長洲",
		"Mui Wo":      "梅窩",
		"Peng Chau":   "坪洲",
		"Chi Ma Wan":  "芝麻灣",
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
		Values: map[int]string{
			dataFetch.TerminatePengChau:   "Terminate at Peng Chau",
			dataFetch.TerminateMuiWo:      "Terminate at Mui Wo",
			dataFetch.TerminateCheungChau: "Terminate at Cheung Chau",
		},
	},
	zhRemarkSetup: binaryFlag.ValuePair[string]{
		Values: map[int]string{
			dataFetch.TerminatePengChau:   "Terminate: 坪洲",
			dataFetch.TerminateMuiWo:      "Terminate: 梅窩",
			dataFetch.TerminateCheungChau: "Terminate: 長洲",
		},
	},
}
var interIslandConvert = ferry.Convert{
	ToSpeed: func(serviceDate string, remark string) binaryFlag.BinaryFlag {
		return *binaryFlag.New().SetBit(dataFetch.SpeedOrdinary)
	},
	ToFrequency: func(serviceDate string, remark string) binaryFlag.BinaryFlag {
		return *binaryFlag.New().SetBit(1).SetBit(2).SetBit(3).SetBit(4).SetBit(5).SetBit(6).SetBit(7)
	},
	ToRemark: func(serviceDate string, remark string) binaryFlag.BinaryFlag {
		if remark == "1" {
			return *binaryFlag.New().SetBit(dataFetch.TerminatePengChau)
		} else if remark == "2" {
			return *binaryFlag.New().SetBit(dataFetch.TerminateMuiWo)
		} else if remark == "3" {
			return *binaryFlag.New().SetBit(dataFetch.TerminateCheungChau)
		} else {
			return *binaryFlag.New()
		}
	},
}
