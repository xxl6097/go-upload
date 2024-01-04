package utils

import (
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
)

func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
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
