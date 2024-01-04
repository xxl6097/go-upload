package utils

import "time"

const timezone = "Asia/Shanghai"

func TimeFormat(date time.Time, pattern string) string {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		location = time.FixedZone("CST", 8*3600) //替换上海时区方式
	}
	date.In(location)
	return date.Format(pattern)
}

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
	//loc, _ := time.LoadLocation("Asia/Shanghai")
	location, err := time.LoadLocation(timezone)
	if err != nil {
		location = time.FixedZone("CST", 8*3600) //替换上海时区方式
	}
	date := time.Now()
	date.In(location)
	return date.Format("2006/01/02/15/04/05/")

}
