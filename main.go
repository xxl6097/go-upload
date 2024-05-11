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
	server.Bootstrap()
}
