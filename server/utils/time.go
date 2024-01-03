package utils

import "time"

// "2023-05-29 15:10:41"
func GetNowStr() string {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	return time.Now().In(loc).Format("2006-01-02 15:04:05")
}

func GetFileNameWithTime() string {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	return time.Now().In(loc).Format("20060102150405")
}

func GetTimeDir() string {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	return time.Now().In(loc).Format("2006/01/02/15/04/05/")
}
