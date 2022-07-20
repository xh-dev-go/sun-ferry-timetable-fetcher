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

type Convert struct {
	ToSpeed     func(string, string) binaryFlag.BinaryFlag
	ToFrequency func(string, string) binaryFlag.BinaryFlag
	ToRemark    func(string, string) binaryFlag.BinaryFlag
}

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

func DecodeIsland(msg, routeName string, dict map[string]string, convert Convert) *[]FerryRecord {
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
		direction := strings.Split(line[0], "to")
		record.From = strings.TrimSpace(direction[0])

		// TODO Bug
		record.From = strings.ReplaceAll(record.From, "Chueung", "Cheung")
		record.To = strings.TrimSpace(direction[1])

		record.ZhFrom = dict[record.From]
		record.ZhTo = dict[record.To]

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
		serviceDate := line[1]
		record.Frequency = convert.ToFrequency(serviceDate, remark)
		record.Speed = convert.ToSpeed(serviceDate, remark)
		record.Remark = convert.ToRemark(serviceDate, remark)

		records = append(records, record)
		fmt.Println(fmt.Sprintf("%+v", record))

	}
	return &records
}
func Decode(msg, routeName string, dict map[string]string, convert Convert) *[]FerryRecord {
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
		serviceDate := line[1]
		record.Frequency = convert.ToFrequency(serviceDate, remark)
		record.Speed = convert.ToSpeed(serviceDate, remark)
		record.Remark = convert.ToRemark(serviceDate, remark)

		records = append(records, record)
		fmt.Println(fmt.Sprintf("%+v", record))

	}
	return &records
}
