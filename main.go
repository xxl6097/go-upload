package main

import (
	"flag"
	"fmt"
	"github.com/xxl6097/go-upload/server"
	"github.com/xxl6097/go-upload/version"
	"golang.org/x/crypto/bcrypt"
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
	//SetPassword("het002402")
	CheckPassword("het002402")
	server.Bootstrap()
}

// go get -u golang.org/x/crypto@v0.23.0
func SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	c := string(hash)
	fmt.Println(c)
	//bcrypt.CompareHashAndPassword()
	return nil
}

func CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte("$2a$10$3OzuBHMUG22tQHQVubcV1.uchgt420yVF5LS4QuZgGtHd3ZFjBaH6"), []byte(password))
	fmt.Println(err)
	return err
}
