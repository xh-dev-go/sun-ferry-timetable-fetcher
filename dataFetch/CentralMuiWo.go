package dataFetch

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/xh-dev-go/xhUtils/binaryFlag"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func Extract(URL string, eTag string) (string, int, string) {
	client := http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		panic(err)
	}
	if eTag != "" {
		req.Header.Set("If-None-Match", eTag)
	}
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if response.StatusCode == 200 {
		msgByte, err := io.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		msg := string(msgByte)
		return msg, 200, response.Header.Get("ETag")

	} else if response.StatusCode == 304 {
		return "", 304, ""

	} else {
		panic(errors.New(fmt.Sprintf("Status[%d] not support", response.StatusCode)))

	}

}

func Decode(msg, routeName string, dict map[string]string) *[]FerryRecord {
	var records []FerryRecord
	csvReader := csv.NewReader(strings.NewReader(msg))
	csvReader.Read()
	lines, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}
	for _, line := range lines {
		record := FerryRecord{}
		record.Route = routeName
		locationArr := strings.Split(line[0], "to")
		record.From = strings.TrimSpace(locationArr[0])
		record.ZhFrom = dict[record.From]
		record.To = strings.TrimSpace(locationArr[1])
		record.ZhTo = dict[record.To]

		mtf := *binaryFlag.New().SetBit(1).SetBit(2).SetBit(3).SetBit(4).SetBit(5)
		mts := *binaryFlag.New().SetBit(1).SetBit(2).SetBit(3).SetBit(4).SetBit(5).SetBit(6).SetBit(7)
		sat := *binaryFlag.New().SetBit(6)
		sunPub := *binaryFlag.New().SetBit(7).SetBit(10)
		if line[1] == "Mondays to Fridays except public holidays" {
			record.Frequency = *binaryFlag.New().SetBinary(mtf)
		} else if line[1] == "Saturdays except public holidays" {
			record.Frequency = *binaryFlag.New().SetBinary(sat)
		} else if line[1] == "Sundays and public holidays" {
			record.Frequency = *binaryFlag.New().SetBinary(sunPub)
		} else if line[1] == "Mondays to Saturdays except public holidays" {
			record.Frequency = *binaryFlag.New().SetBinary(mts)
		} else {
			panic(errors.New(fmt.Sprintf("Line[%s] not support", line[1])))
		}
		times := strings.Split(strings.TrimSpace(line[2]), " ")
		time, err := strconv.Atoi(strings.ReplaceAll(times[0], ":", ""))
		if err != nil {
			panic(errors.New("fail convert time to number"))
		}
		if times[1] == "p.m." {
			time += 1200
		}
		record.Time = time

		remark := line[3]
		if remark == "" {
			record.Speed = *binaryFlag.New().SetBit(SpeedFast)
			record.Remark = *binaryFlag.New()
		} else if remark == "1" {
			record.Speed = *binaryFlag.New().SetBit(SpeedOrdinary)
			record.Remark = *binaryFlag.New()
		} else if remark == "2" {
			record.Speed = *binaryFlag.New().SetBit(SpeedOrdinary)
			record.Remark = *binaryFlag.New().SetBit(ViaPengChau)
		}

		records = append(records, record)
		fmt.Println(fmt.Sprintf("%+v", record))

	}
	return &records
}
