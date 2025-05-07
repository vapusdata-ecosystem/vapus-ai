package dmutils

import (
	"fmt"
	"time"
)

type DatetimeOpts struct {
	minutes         int64
	hours           int64
	days            int64
	exactStringDate int64
	weeks           int64
	unix            bool
	format          string
}

type DateLog struct {
	Date int64
}

func GetEpochTime() int64 {
	return time.Now().Unix()
}

func GetMilliEpochTime() int64 {
	return time.Now().UnixMilli()
}

func GetCurrent(opt ...DatetimeOpts) int64 {
	return time.Now().Unix()
}

func GetTime(opt ...DatetimeOpts) int64 {
	return time.Now().Add(5).Unix()
}

func GetFormattedTime(t int64, format string) string {
	return time.Unix(t, 0).Format(format)
}

func GetStrEpochTime() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}
