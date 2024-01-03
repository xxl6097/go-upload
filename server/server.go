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
	"strings"
)

var (
	fsys http.FileSystem
)

func init() {
	fsys = assets.Load()
}

func Ok(data interface{}) map[string]interface{} {
	return map[string]interface{}{"code": 0, "msg": "sucess", "data": data}
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
		filearr := utils.VisitDir("./files")
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
		var filearr []utils.FileStruct
		for i, _ := range files {
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			dir := "./files/" + utils.GetTimeDir()
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
	router.PathPrefix("/files/").Handler(http.StripPrefix("/files/", http.FileServer(http.Dir("./files"))))

	router.Use(mux.CORSMethodMiddleware(router))
	router.HandleFunc("/upload", upload).Methods(http.MethodPost, http.MethodOptions) // view
	router.HandleFunc("/upload", upload).Methods(http.MethodGet, http.MethodOptions)  // view
	router.Handle("/favicon.ico", http.FileServer(fsys)).Methods("GET")
	router.PathPrefix("/").Handler(utils.MakeHTTPGzipHandler(http.StripPrefix("/", http.FileServer(fsys)))).Methods("GET")
}

func FileUploadWebServer(port int, token string) {
	router := mux.NewRouter()
	initRouter(router)
	address := fmt.Sprintf(":%d", port)
	server := &http.Server{
		Addr:    address,
		Handler: router,
	}
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return
	}
	fmt.Printf("please open http://localhost%s\n", server.Addr)
	fmt.Printf("please open http://localhost%s/files/\n", server.Addr)
	_ = server.Serve(ln)
}
