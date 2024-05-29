package main

import (
	"flag"
	"github.com/xxl6097/go-upload/server"
	"github.com/xxl6097/go-upload/version"
)

func main() {
	var isVersion bool
	flag.BoolVar(&isVersion, "v", false, "version")
	flag.Parse()
	if isVersion {
		version.Version()
		return
	}
	serve()
}

func serve() {
	//path := "/Users/uuxia/Desktop/work/code/go/go-upload/files"
	//os.Setenv("ENV_FILES", path)
	server.Bootstrap()
}
