package main

import "github.com/xxl6097/go-upload/server"

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

	server.Bootstrap()
}
