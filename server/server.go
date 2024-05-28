package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/xxl6097/go-upload/assets"
	"github.com/xxl6097/go-upload/server/utils"
	"github.com/xxl6097/go-upload/version"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

// curl -H "Authorization: 44" -F "file=@/Users/uuxia/Desktop/work/code/go/go-upload/main.go" http://localhost:4444/upload
var (
	fsys      http.FileSystem
	token     string
	origin    string
	files_dir string = "./files"
	//files_dir            = "/Users/uuxia/Desktop/work/code/go/go-upload"
	static_prefix string = "/files/"
	my                   = "/my"
	_port                = "8087"
)

func init() {
	fsys = assets.Load()
}

func Ok(data interface{}) map[string]interface{} {
	return map[string]interface{}{"code": 0, "msg": "sucess", "data": data}
}

func Result(code int, msg string, data interface{}) map[string]interface{} {
	return map[string]interface{}{"code": code, "msg": msg, "data": data}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	//fmt.Println(data)
	if json.NewEncoder(w).Encode(data) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Respond1(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	//fmt.Println(data)
	if json.NewEncoder(w).Encode(data) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetReqData[T any](w http.ResponseWriter, r *http.Request) *T {
	var t T
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	return &t
}
func getip(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(getPubIP()))
	if err != nil {
		fmt.Println("无法写入响应:", err)
	}
}

func auth(w http.ResponseWriter, r *http.Request) {
	_token := r.Header.Get("Authorization")
	if strings.ToLower(token) == strings.ToLower(_token) {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(502)
	}
	//w.Write([]byte(utils.Version))
}

func config(w http.ResponseWriter, r *http.Request) {
	if runtime.GOOS == "darwin" {
		Respond(w, Ok(map[string]interface{}{
			"AppName":      "goupload",
			"AppVersion":   "appversion",
			"BuildVersion": "BuildVersion",
			"BuildTime":    "BuildTime",
			"GitRevision":  "GitRevision",
			"GitBranch":    "GitBranch",
			"GoVersion":    "GoVersion",
		}))
	} else {
		Respond(w, Ok(map[string]interface{}{
			"AppName":      version.AppName,
			"AppVersion":   version.AppVersion,
			"BuildVersion": version.BuildVersion,
			"BuildTime":    version.BuildTime,
			"GitRevision":  version.GitRevision,
			"GitBranch":    version.GitBranch,
			"GoVersion":    version.GoVersion,
		}))
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	_authorization := r.Header.Get("Authorization")
	if strings.ToLower(token) != strings.ToLower(_authorization) {
		w.WriteHeader(502)
		Respond(w, Result(-1, "请检查Authorization是否正确!", _authorization))
		return
	}
	switch r.Method {
	case "GET": //获取目录或者子目录下的所有文件
		queryParams := r.URL.Query()
		origin = queryParams.Get("origin")
		fmt.Println("origin", origin)
		filearr := utils.VisitDir(files_dir, static_prefix)
		sort.Slice(filearr, func(i, j int) bool {
			return filearr[i].ModTime.Before(filearr[j].ModTime)
		})
		Respond(w, Ok(filearr))
	case "DELETE":
		req := GetReqData[map[string]interface{}](w, r)
		files := (*req)["files"]
		if files == nil {
			Respond(w, Result(-1, "path is nil", nil))
			return
		}
		var res = make([]interface{}, 0)
		for _, path := range files.([]interface{}) {
			realpath := files_dir + path.(string)[len(static_prefix)-1:]
			err := os.Remove(realpath)
			if err != nil {
				msg := fmt.Sprintf("[%s] 删除失败:%s", realpath, err.Error())
				var respone = struct {
					Path   string `json:"path"`
					Sucess bool   `json:"sucess"`
				}{realpath, false}
				res = append(res, respone)
				fmt.Println(msg)
			} else {
				msg := fmt.Sprintf("[%s] 删除成功", realpath)
				var respone = struct {
					Path   string `json:"path"`
					Sucess bool   `json:"sucess"`
				}{realpath, true}
				res = append(res, respone)
				fmt.Println(msg)
			}
		}
		Respond(w, Result(0, "", res))
	case "POST":
		source := r.Header.Get("source")
		//ParseMultipartForm将请求的主体作为multipart/form-data解析。请求的整个主体都会被解析，得到的文件记录最多 maxMemery字节保存在内存，其余部分保存在硬盘的temp文件里。如果必要，ParseMultipartForm会自行调用 ParseForm。重复调用本方法是无意义的
		//设置内存大小
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//files := r.MultipartForm.File["files"]
		m := r.MultipartForm
		files := m.File["file"]
		if files == nil {
			Respond(w, Result(-1, "请在MultipartForm字段中添加file字段和对应文件", nil))
			return
		}
		var filearrs []string
		var filearr []interface{}
		for i, _ := range files {
			file, err2 := files[i].Open()
			defer file.Close()
			if err2 != nil {
				http.Error(w, err2.Error(), http.StatusInternalServerError)
				return
			}
			dir := files_dir + "/" + utils.GetDirAtDay()
			//判断文件夹是否存在，不存在则创建文件夹
			if !utils.IsDirExists(dir) {
				err1 := utils.CreateMutiDir(dir)
				if err1 != nil {
					http.Error(w, err1.Error(), http.StatusInternalServerError)
					return
				}
			}
			filePath := dir + files[i].Filename
			_, err1 := os.Stat(filePath)
			if err1 == nil {
				fileName := filepath.Base(filePath)
				ext := filepath.Ext(fileName)
				fileName = fileName[:len(fileName)-len(filepath.Ext(fileName))]
				// 重命名文件
				newFile := dir + fileName + "_" + utils.GetFileNameWithTime() + ext
				err12 := os.Rename(filePath, newFile)
				if err12 != nil {
					fmt.Println("重命名文件时发生错误:", err)
					//文件存在，删除文件
					err2 := os.Remove(filePath)
					if err2 != nil {
						fmt.Println("删除文件时发生错误:", err)
					} else {
						fmt.Println("文件存在，已删除:", filePath)
					}
				}
			}
			dst, err3 := os.Create(filePath)
			defer dst.Close()
			if err3 != nil {
				http.Error(w, err3.Error(), http.StatusInternalServerError)
				return
			}
			//copy the uploaded file to the destination file
			if _, err4 := io.Copy(dst, file); err4 != nil {
				http.Error(w, err4.Error(), http.StatusInternalServerError)
				return
			}
			fileInfo, err5 := os.Stat(filePath)
			if err5 != nil {
				fmt.Println("Error:", err5, source)
				http.Error(w, err5.Error(), http.StatusInternalServerError)
				return
			}
			_path := static_prefix + filePath[len(files_dir)+1:]
			//item := utils.FileStruct{Name: fileInfo.Name(), Size: fileInfo.Size(), Path: _path, ModTime: fileInfo.ModTime().String()}
			item := utils.FileStruct{Name: fileInfo.Name(), Size: fileInfo.Size(), Path: _path, ModTime: fileInfo.ModTime()}
			//fmt.Println(item, _path)
			if source == "web" {
				filearr = append(filearr, item)
			} else {
				filearrs = append(filearrs, origin+_path)
			}
			fmt.Printf("文件上传成功:%s,%+v", filePath, item)
		}
		//curl -F "file=@/Users/uuxia/Desktop/work/code/go/go-upload/build.sh"  -H "Authorization: 88" http://localhost:8888/upload
		if source == "web" {
			Respond(w, Ok(filearr))
		} else {
			result := strings.Join(filearrs, "\r\n")
			fmt.Println(result)
			w.Write([]byte(result + "\r\n"))
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
func up(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "text/plain")
	// 编写要回复的数据
	responseText := "#!/bin/bash\ncmd=\"curl\"\nheader=\"-H \\\"Authorization: " + token + "\\\"\"\nhost=\"http://" + r.Host + "/upload\"\nfiles=\"\"\nfor arg in \"$@\"; do\n  if [[ $arg == /* ]]; then\n      files+=\"-F \\\"file=@$arg\\\" \"\n  else\n      absolute_path=$(realpath \"$arg\")\n      files+=\"-F \\\"file=@$absolute_path\\\" \"\n  fi\ndone\ncmd=\"curl $header $files$host\"\necho \"$cmd\"\neval $cmd"
	fmt.Println(responseText)
	// 将数据写入响应
	_, err := w.Write([]byte(responseText))
	if err != nil {
		fmt.Println("无法写入响应:", err)
	}
}
func fileserver(w http.ResponseWriter, r *http.Request) {

}

func initRouter(router *mux.Router) {
	//http server
	router.PathPrefix(static_prefix).Handler(http.StripPrefix(static_prefix, http.FileServer(http.Dir(files_dir))))
	//router.PathPrefix("/a/").Handler(http.StripPrefix("/a/", http.FileServer(http.Dir(dir))))
	if utils.IsDirExists(my) {
		sub := router.NewRoute().Subrouter()
		sub.Use(utils.NewHTTPAuthMiddleware("admin", "het002402").Middleware)
		sub.PathPrefix("/fileserver/").Handler(http.StripPrefix("/fileserver/", http.FileServer(http.Dir(my))))
	}

	router.Use(mux.CORSMethodMiddleware(router))
	router.HandleFunc("/config", config).Methods(http.MethodPost, http.MethodOptions) // view
	router.HandleFunc("/auth", auth).Methods(http.MethodPost, http.MethodOptions)     // view
	router.HandleFunc("/upload", upload).Methods(http.MethodPost, http.MethodOptions) // view
	router.HandleFunc("/upload", upload).Methods(http.MethodGet, http.MethodOptions)  // view
	router.HandleFunc("/getip", getip).Methods(http.MethodGet, http.MethodOptions)    // view
	router.HandleFunc("/up", up).Methods(http.MethodGet, http.MethodOptions)          // view
	//router.HandleFunc("/up", upload).Methods(http.MethodGet, http.MethodOptions)             // view
	router.HandleFunc("/upload", upload).Methods(http.MethodDelete, http.MethodOptions)      // view
	router.HandleFunc("/fileserver", fileserver).Methods(http.MethodGet, http.MethodOptions) // view

	router.Handle("/favicon.ico", http.FileServer(fsys)).Methods("GET")
	router.PathPrefix("/").Handler(utils.MakeHTTPGzipHandler(http.StripPrefix("/", http.FileServer(fsys)))).Methods("GET")
}

func FileUploadWebServer(port, _token string) {
	router := mux.NewRouter()
	initRouter(router)
	address := fmt.Sprintf(":%s", port)
	server := &http.Server{
		Addr:    address,
		Handler: router,
	}
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return
	}
	token = _token
	welcom(port, _token)
	_ = server.Serve(ln)
}

func preArgs() {
	_os := runtime.GOOS
	// 根据操作系统执行不同的逻辑
	switch _os {
	case "darwin":
		fmt.Println("当前运行在 macOS 操作系统上")
		os.Setenv("ENV_PORT", "8888")
		os.Setenv("ENV_TOKEN", "88")
		//os.Setenv("ENV_FILES", "/Users/uuxia/Desktop/work/doc")
	case "windows":
		fmt.Println("当前运行在 Windows 操作系统上")
	default:
		fmt.Println("无法识别的操作系统")
	}
}

// /Users/uuxia/Desktop/work
func Bootstrap() {
	preArgs()
	var port = os.Getenv("ENV_PORT")
	var token = os.Getenv("ENV_TOKEN")
	_dir := os.Getenv("ENV_FILES")
	if _dir != "" {
		files_dir = _dir
	}
	if port == "" && token == "" {
		switch len(os.Args) {
		case 3:
			port = os.Args[1]
			token = os.Args[2]
			if port == "" || token == "" {
				fmt.Printf("请正确输入端口和token等参数")
				return
			}
		default:
			for {
				fmt.Printf("请输入端口号：")
				_, err := fmt.Scanln(&port)
				if err != nil {
					fmt.Println("输入错误：", err)
					continue
				}
				if !utils.IsNumeric(port) {
					fmt.Println("请输入一个数字,谢谢!")
					continue
				}
				break
			}
			for {
				fmt.Printf("请设置token：")
				_, err := fmt.Scanln(&token)
				if err != nil {
					fmt.Println("输入错误：", err)
					continue
				}
				break
			}
		}
	}
	ip := utils.GetHostIp()
	if origin == "" {
		origin = fmt.Sprintf("http://%s:%s", ip, port)
	}
	_port = port

	if os.Getenv("FRP_DOWN") == "true" {
		go FrpcDown(files_dir)
	}
	FileUploadWebServer(port, token)
}

func welcom(port, token string) {
	fmt.Println("欢迎使用文件上传助手")
	fmt.Printf("文件路径：%s\n", files_dir)
	fmt.Printf("网页上传：http://localhost:%s\n", port)
	fmt.Printf("网页上传：http://localhost:%s%s\n", port, static_prefix)
	fmt.Printf("指令上传示例：curl -F \"file=@/root/xxx.log\" -F \"token=%s\" http://localhost:%s/upload\n", token, port)
}

type IPResponse struct {
	IP string `json:"ip"`
}

func getPubIP() string {
	// 发送GET请求到API
	resp, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		fmt.Println("请求失败:", err)
		return ""
	}
	defer resp.Body.Close()

	// 解析JSON响应
	var ipResponse IPResponse
	err = json.NewDecoder(resp.Body).Decode(&ipResponse)
	if err != nil {
		fmt.Println("解析JSON失败:", err)
		return ""
	}

	// 输出公网IP地址
	fmt.Println("公网IP地址:", ipResponse.IP)
	return fmt.Sprintf("http://%s:%s", ipResponse.IP, _port)
}
func initsh() {
	//sh := "#!/bin/bash\ncmd=\"curl \"\nfor arg in \"$@\"; do\n  if [[ $arg == /* ]]; then\n      cmd+=\"-F \\\"file=@$arg\\\" \"\n  else\n      absolute_path=$(realpath \"$arg\")\n      cmd+=\"-F \\\"file=@$absolute_path\\\" \"\n  fi\ndone\ncmd+=\"-F \\\"token=het002402\\\" http://uuxia.cn:8087/upload\"\necho \"运行命令：$cmd\"\neval $cmd\n\n"
}
