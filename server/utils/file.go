package utils

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type FileStruct struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Size int64  `json:"size"`
}

func IsDirExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// 判断所给路径文件/文件夹是否存在(返回true是存在)
func isExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func CreateMutiDir(filePath string) error {
	if !isExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			fmt.Println("创建文件夹失败,error info:", err)
			return err
		}
		return err
	}
	return nil
}

func visit(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Println(err) // 可以根据需要处理错误
		return nil
	}
	if info.IsDir() {
		// 处理目录
		//fmt.Println("Directory:", path)
	} else {
		// 处理文件
		fmt.Println("File:", path)
	}
	return nil
}

func VisitDir(rootDir string) []FileStruct {
	var filearr []FileStruct
	err := filepath.Walk(rootDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err) // 可以根据需要处理错误
			return nil
		}
		if info.IsDir() {
			// 处理目录
			//fmt.Println("Directory:", path)
		} else {
			// 处理文件
			//fmt.Println("File:", path, info.Name())
			if path != "" {
				item := FileStruct{Name: info.Name(), Size: info.Size(), Path: path}
				//fmt.Println("File:", item)
				filearr = append(filearr, item)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return filearr
}
