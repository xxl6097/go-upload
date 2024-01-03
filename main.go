package main

import (
	"github.com/xxl6097/go-upload/server"
)

func main() {
	server.TestFileServer()

	//dir := "./files"
	//fileServer := http.FileServer(http.Dir(dir))
	//// 将文件服务器注册到路径 "/static/"
	//http.Handle("/files/", http.StripPrefix("/files/", fileServer))
	//// 启动服务器
	//port := "8080"
	//println("Server is running on http://localhost:" + port + "/files/")
	//err := http.ListenAndServe(":"+port, nil)
	//if err != nil {
	//	println("Error:", err)
	//}
}
