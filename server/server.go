package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/xxl6097/go-upload/assets"
	"github.com/xxl6097/go-upload/server/utils"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	fsys          http.FileSystem
	token         string
	files_dir     string = "./files"
	static_prefix string = "/files/"
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
	if json.NewEncoder(w).Encode(data) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		filearr := utils.VisitDir(files_dir, static_prefix)
		Respond(w, Ok(filearr))
	case "POST":
		//ParseMultipartForm将请求的主体作为multipart/form-data解析。请求的整个主体都会被解析，得到的文件记录最多 maxMemery字节保存在内存，其余部分保存在硬盘的temp文件里。如果必要，ParseMultipartForm会自行调用 ParseForm。重复调用本方法是无意义的
		//设置内存大小
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		m := r.MultipartForm
		files := m.File["file"]
		_token := m.Value["token"]
		if files == nil {
			Respond(w, Result(-1, "请在MultipartForm字段中添加file字段和对应文件", nil))
			return
		}
		if _token == nil {
			Respond(w, Result(-1, "你带上token字段", nil))
			return
		}
		if strings.ToLower(token) == strings.ToLower(_token[0]) {
		} else {
			Respond(w, Result(-1, "请检查token是否正确!", _token))
			return
		}
		var filearr []utils.FileStruct
		for i, _ := range files {
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			dir := files_dir + "/" + utils.GetTimeDir()
			if !utils.IsDirExists(dir) {
				err := utils.CreateMutiDir(dir)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			filePath := dir + files[i].Filename
			dst, err := os.Create(filePath)
			defer dst.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//copy the uploaded file to the destination file
			if _, err := io.Copy(dst, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fileName := filepath.Base(filePath)
			item := utils.FileStruct{Name: fileName, Size: files[i].Size, Path: strings.TrimPrefix(dst.Name(), ".")}
			fmt.Println(item)
			filearr = append(filearr, item)
		}

		Respond(w, Ok(filearr))

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func initRouter(router *mux.Router) {
	//http server
	router.PathPrefix(static_prefix).Handler(http.StripPrefix(static_prefix, http.FileServer(http.Dir(files_dir))))
	//router.PathPrefix("/a/").Handler(http.StripPrefix("/a/", http.FileServer(http.Dir(dir))))

	router.Use(mux.CORSMethodMiddleware(router))
	router.HandleFunc("/upload", upload).Methods(http.MethodPost, http.MethodOptions) // view
	router.HandleFunc("/upload", upload).Methods(http.MethodGet, http.MethodOptions)  // view
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
		os.Setenv("ENV_PORT", "4444")
		os.Setenv("ENV_TOKEN", "44")
		os.Setenv("ENV_FILES", "/Users/uuxia/Desktop/work/doc")
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
	FileUploadWebServer(port, token)
}

func welcom(port, token string) {
	fmt.Println("欢迎使用文件上传助手")
	fmt.Printf("文件路径：%s\n", files_dir)
	fmt.Printf("网页上传：http://localhost:%s\n", port)
	fmt.Printf("网页上传：http://localhost:%s%s\n", port, files_dir)
	fmt.Printf("指令上传示例：curl -F \"file=@/root/xxx.log\" -F \"token=%s\" http://localhost:%s/upload\n", token, port)
}
