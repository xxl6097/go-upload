package main

import (
	"fmt"
	"github.com/xxl6097/go-upload/server"
	"github.com/xxl6097/go-upload/server/utils"
	"net"
	"time"
)

var (
	Version string
)

func checkPort(host string, port string, second int) bool {
	conn, err := net.DialTimeout("tcp", host+":"+port, time.Duration(second)*time.Second)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	defer conn.Close()
	return true
}

func main() {
	//dir := "/Users/uuxia/Desktop/work"
	//fileServer := http.FileServer(http.Dir(dir))
	//http.Handle("/", http.StripPrefix("/", fileServer))
	//port := 8080
	//fmt.Printf("文件服务器已启动，访问地址：http://localhost:%d\n", port)
	//err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	//if err != nil {
	//	fmt.Println("服务器启动失败:", err)
	//}

	//is := checkPort("10.16.14.103", "31381", 300)
	//fmt.Println(is)
	utils.Version = Version
	fmt.Println("====>", Version)
	server.Bootstrap()
}
