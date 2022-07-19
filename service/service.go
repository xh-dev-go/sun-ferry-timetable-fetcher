package service

import "github.com/xh-dev-go/sun-ferry-timetable-fetcher/dataFetch"

func get() []dataFetch.FerryRecord {
	str := dataFetch.Extract("https://www.sunferry.com.hk/eta/timetable/SunFerry_central_muiwo_timetable_eng.csv")
	dict := map[string]string{
		"Mui Wo":  "梅窩",
		"Central": "中環",
	}
	return dataFetch.Decode(str, "Central - Mui Wo", dict)
}
