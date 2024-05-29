package utils

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileStruct struct {
	Name    string    `json:"name"`
	Path    string    `json:"path"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
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

func GetTree(path, prefix, defaultDir string) []FileStruct {
	root := Node{Title: filepath.Base(path)}
	finfo, err := os.Stat(path)
	if err != nil {
		return nil
	}
	root.Path = path
	root.Spread = true

	if !finfo.IsDir() {
		return nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil
	}
	var files []FileStruct
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		var _path string
		info, _ := entry.Info()
		if strings.HasPrefix(path, defaultDir) {
			// 去掉字符串的前缀
			_path = strings.TrimPrefix(path, defaultDir)
			//path = strings.TrimPrefix(path, rootDir)
			if strings.HasPrefix(_path, "/") {
				index := strings.Index(_path, "/")
				if index != -1 {
					_path = _path[index+1:]
				}
				if !strings.HasSuffix(_path, "/") {
					_path += "/"
				}
			}
		} else {
			if !strings.HasPrefix(path, "/") {
				index := strings.Index(path, "/")
				// 如果找到了 '/'，则截取字符串
				if index != -1 {
					_path = path[index+1:]
				}
				if !strings.HasSuffix(_path, "/") {
					_path += "/"
				}
			}
		}
		item := FileStruct{Name: info.Name(), Size: info.Size(), Path: prefix + _path + entry.Name(), ModTime: info.ModTime()}
		files = append(files, item)
	}

	return files
}

func VisitDir(rootDir, prefix string) []FileStruct {
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
			// 对路径进行编码
			//fmt.Println("File:", path, info.ModTime().String())
			if path != "" {
				// 判断字符串是否以指定的前缀开头
				if strings.HasPrefix(path, rootDir) {
					// 去掉字符串的前缀
					path = strings.TrimPrefix(path, rootDir)
					//path = strings.TrimPrefix(path, rootDir)
					if strings.HasPrefix(path, "/") {
						index := strings.Index(path, "/")
						if index != -1 {
							path = path[index+1:]
						}
					}
				} else {
					if !strings.HasPrefix(path, "/") {
						index := strings.Index(path, "/")
						// 如果找到了 '/'，则截取字符串
						if index != -1 {
							path = path[index+1:]
						}
					}
				}
				//item := FileStruct{Name: info.Name(), Size: info.Size(), Path: prefix + path, ModTime: info.ModTime().String()}
				item := FileStruct{Name: info.Name(), Size: info.Size(), Path: prefix + path, ModTime: info.ModTime()}
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

func IsPath(path string) bool {
	// Clean the path to handle any relative or redundant elements
	cleanedPath := filepath.Clean(path)

	// Check if the path exists
	_, err := os.Stat(cleanedPath)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil || os.IsExist(err)
}
