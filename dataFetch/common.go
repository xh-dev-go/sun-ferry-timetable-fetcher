package dataFetch

import "github.com/xh-dev-go/xhUtils/binaryFlag"

type FerryRecord struct {
	Route     string
	From      string
	ZhFrom    string
	To        string
	ZhTo      string
	Frequency binaryFlag.BinaryFlag
	Time      int
	Speed     binaryFlag.BinaryFlag
	Remark    binaryFlag.BinaryFlag
}

const (
	SpeedFast     int = 1
	SpeedOrdinary     = 2
)

const (
	ViaPengChau int = 1
)
