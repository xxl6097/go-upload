package utils

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

// 耗时统计，eg：defer TrackTime(time.Now())
func TrackTime(pre time.Time) time.Duration {
	elapsed := time.Since(pre)
	fmt.Println("elapsed:", elapsed)
	return elapsed
}

// 两阶段延时执行 eg:defer setTeardown()
func setupTeardown() func() {
	fmt.Println("init")
	return func() {
		fmt.Println("end..")
	}
}
func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
func GetHostIp() string {
	addrList, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("get current host ip err: ", err)
		return ""
	}
	var ip string
	for _, address := range addrList {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.IsPrivate() {
			if ipNet.IP.To4() != nil {
				ip = ipNet.IP.String()
				break
			}
		}
	}
	return ip
}
func GetPublicIP() (string, error) {
	resp, err := http.Get("https://ipinfo.io/ip")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}

func GetLocalIPs() ([]string, error) {
	var ips []string

	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	// 遍历每个网络接口
	for _, i := range interfaces {
		// 排除虚拟和回环接口
		if i.Flags&net.FlagUp == 0 || i.Flags&net.FlagLoopback != 0 {
			continue
		}

		// 获取接口的地址
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}

		// 遍历每个地址
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP.To4() != nil {
					// IPv4地址
					ips = append(ips, v.IP.String())
				}
			case *net.IPAddr:
				if v.IP.To4() != nil {
					// IPv4地址
					ips = append(ips, v.IP.String())
				}
			}
		}
	}

	return ips, nil
}

func FuzzySearch[T any](pattern string, data []T, f func(T) string) []T {
	var matches []T
	// 编译正则表达式
	re, err := regexp.Compile("(?i)" + pattern)
	if err != nil {
		fmt.Println("Invalid pattern")
		return matches
	}

	// 遍历数组，寻找匹配项
	for _, item := range data {
		if re.MatchString(f(item)) {
			matches = append(matches, item)
		}
	}
	return matches
}

// 模糊匹配搜索函数
func FuzzySearch1(pattern string, data []interface{}, f func(interface{}) string) []interface{} {
	var matches []interface{}
	// 编译正则表达式
	re, err := regexp.Compile("(?i)" + pattern)
	if err != nil {
		fmt.Println("Invalid pattern")
		return matches
	}

	// 遍历数组，寻找匹配项
	for _, item := range data {
		if re.MatchString(f(item)) {
			matches = append(matches, item)
		}
	}
	return matches
}
