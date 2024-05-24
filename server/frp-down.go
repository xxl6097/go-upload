package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	owner = "fatedier" // frp 项目的所有者
	repo  = "frp"      // frp 项目的仓库名称
)

// Release 表示 GitHub 仓库的发布版本
type Release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name        string `json:"name"`
		DownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func FrpcDown(path string) {
	for {
		check(path)
		time.Sleep(time.Hour * 1)
	}
}

func check(path string) {
	// 获取 frp 项目的发布版本信息
	releases, err := getReleases(owner, repo)
	if err != nil {
		fmt.Println("Error fetching releases:", err)
		return
	}

	frp_dir := filepath.Join(path, repo)

	// 下载每个发布版本的文件
	for _, release := range releases {
		fmt.Println("Downloading release:", release.TagName)
		CreateMutiDir(filepath.Join(frp_dir, release.TagName))
		err := downloadAssets(filepath.Join(frp_dir, release.TagName), release)
		if err != nil {
			fmt.Println("Error downloading assets for release", release.TagName, ":", err)
		}
	}
}

// 获取 frp 项目的发布版本信息
func getReleases(owner, repo string) ([]Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", owner, repo)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var releases []Release
	err = json.NewDecoder(resp.Body).Decode(&releases)
	if err != nil {
		return nil, err
	}

	return releases, nil
}

// 下载发布版本的文件
func downloadAssets(path string, release Release) error {
	for _, asset := range release.Assets {
		if IsExist(filepath.Join(path, asset.Name)) {
			fmt.Println("Asset already exists:", asset.Name)
			continue
		} else {
			fmt.Println("Downloading asset:", asset.Name, asset.DownloadURL)
		}
		resp, err := http.Get(asset.DownloadURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// 创建文件
		file, err := os.Create(filepath.Join(path, asset.Name))
		if err != nil {
			return err
		}
		defer file.Close()

		// 将响应体写入文件
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return err
		}
	}
	return nil
}

// 判断所给路径文件/文件夹是否存在(返回true是存在)
func IsExist(path string) bool {
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
	if !IsExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			fmt.Println("创建文件夹失败,error info:", err)
			return err
		}
		return err
	}
	return nil
}
