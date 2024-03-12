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
	"runtime"
	"strings"
)

var (
	fsys      http.FileSystem
	token     string
	origin    string
	files_dir string = "./files"
	//files_dir            = "/Users/uuxia/Desktop/work/code/go/go-upload"
	static_prefix string = "/files/"
	my                   = "/my"
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

func GetReqData[T any](w http.ResponseWriter, r *http.Request) *T {
	var t T
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	return &t
}

func upload(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET": //获取目录或者子目录下的所有文件
		queryParams := r.URL.Query()
		origin = queryParams.Get("origin")
		fmt.Println("origin", origin)
		filearr := utils.VisitDir(files_dir, static_prefix)
		Respond(w, Ok(filearr))
	case "DELETE":
		req := GetReqData[map[string]interface{}](w, r)
		_token := (*req)["token"]
		if _token == nil || _token.(string) == "" {
			Respond(w, Result(-1, "不好意思，没有token，请滚蛋～", nil))
			return
		}
		if strings.ToLower(token) != strings.ToLower(_token.(string)) {
			Respond(w, Result(-1, "不好意思，token无效，请滚蛋～", _token))
			return
		}
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
		var filearr []string
		//var filesstr []string
		for i, _ := range files {
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			dir := files_dir + "/" + utils.GetDirAtDay()
			//判断文件夹是否存在，不存在则创建文件夹
			if !utils.IsDirExists(dir) {
				err := utils.CreateMutiDir(dir)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			filePath := dir + files[i].Filename
			_, err1 := os.Stat(filePath)
			if err1 == nil {
				// 删除文件
				err2 := os.Remove(filePath)
				if err2 != nil {
					fmt.Println("删除文件时发生错误:", err)
				} else {
					fmt.Println("文件存在，已删除:", filePath)
				}
			}
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
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				fmt.Println("Error:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_path := static_prefix + filePath[len(files_dir)+1:]
			item := utils.FileStruct{Name: fileInfo.Name(), Size: fileInfo.Size(), Path: _path, ModTime: fileInfo.ModTime().String()}
			//fmt.Println(item, _path)
			filearr = append(filearr, origin+_path)
			fmt.Printf("文件上传成功:%s,%+v", filePath, item)
		}

		Respond(w, Ok(filearr))

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
func up(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "text/plain")
	// 编写要回复的数据
	responseText := "#!/bin/bash\n"
	responseText += "cmd=\"curl \"\n"
	responseText += "for arg in \"$@\"; do\n"
	responseText += "  if [[ $arg == /* ]]; then\n"
	responseText += "      cmd+=\"-F \\\"file=@$arg\\\" \"\n"
	responseText += "  else\n"
	responseText += "      absolute_path=$(realpath \"$arg\")\n"
	responseText += "      cmd+=\"-F \\\"file=@$absolute_path\\\" \"\n"
	responseText += "  fi\n"
	responseText += "done\n"
	responseText += "echo \"please input token\"\n"
	responseText += "read token\n"
	responseText += "cmd+=\"-F \\\"token=$token\\\" " + origin + "/upload\"\n"
	responseText += "echo \"run cmd: $cmd\"\n"
	responseText += "eval $cmd"
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
	router.HandleFunc("/upload", upload).Methods(http.MethodPost, http.MethodOptions) // view
	router.HandleFunc("/upload", upload).Methods(http.MethodGet, http.MethodOptions)  // view
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
		os.Setenv("ENV_PORT", "4444")
		os.Setenv("ENV_TOKEN", "44")
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
	FileUploadWebServer(port, token)
}

func welcom(port, token string) {
	fmt.Println("欢迎使用文件上传助手")
	fmt.Printf("文件路径：%s\n", files_dir)
	fmt.Printf("网页上传：http://localhost:%s\n", port)
	fmt.Printf("网页上传：http://localhost:%s%s\n", port, static_prefix)
	fmt.Printf("指令上传示例：curl -F \"file=@/root/xxx.log\" -F \"token=%s\" http://localhost:%s/upload\n", token, port)
}
func initsh() {
	//sh := "#!/bin/bash\ncmd=\"curl \"\nfor arg in \"$@\"; do\n  if [[ $arg == /* ]]; then\n      cmd+=\"-F \\\"file=@$arg\\\" \"\n  else\n      absolute_path=$(realpath \"$arg\")\n      cmd+=\"-F \\\"file=@$absolute_path\\\" \"\n  fi\ndone\ncmd+=\"-F \\\"token=het002402\\\" http://uuxia.cn:8087/upload\"\necho \"运行命令：$cmd\"\neval $cmd\n\n"
}
