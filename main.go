package main

import (
	"flag"
	"github.com/xxl6097/go-upload/server"
	"github.com/xxl6097/go-upload/version"
)

func main() {
	// 构建信息，golang版本 commit id 时间
	var isVersion bool
	flag.BoolVar(&isVersion, "v", false, "version")
	flag.Parse()
	if isVersion {
		version.Version()
		return
	}
	server.Bootstrap()
}
