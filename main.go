package main

import (
	"fmt"
	"github.com/xxl6097/go-upload/server"
	"os"
	"strconv"
)

// https://blog.51cto.com/u_87634/7140335
func main() {
	//for idx, args := range os.Args {
	//	fmt.Println("参数"+strconv.Itoa(idx)+":", args)
	//}
	size := len(os.Args)
	fmt.Printf("--->%d", size)
	var port = 8080
	var token = "het002402"
	switch size {
	case 3:
		a := os.Args[1]
		token = os.Args[2]
		_port, err := strconv.Atoi(a)
		if err == nil {
			port = _port
		} else {
			fmt.Printf("输入错误,请输入数字")
		}

	default:
		fmt.Printf("请输入端口号:")
		fmt.Scanln(&port)
		fmt.Printf("请设置token:")
		fmt.Scanln(&token)
		fmt.Println(port, token)
		break
	}

	server.FileUploadWebServer(port, token)
}
