package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

const dir = "./files"

func Start() {
	http.HandleFunc("/files", listFiles)
	http.HandleFunc("/files/", getFile)
	http.HandleFunc("/files/", saveFile)
	http.ListenAndServe(":8080", nil)
}

func listFiles(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(dir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	json.NewEncoder(w).Encode(fileNames)
}

func getFile(w http.ResponseWriter, r *http.Request) {
	fileName := filepath.Base(r.URL.Path)
	fileContent, err := os.ReadFile(filepath.Join(dir, fileName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Write(fileContent)
}

func saveFile(w http.ResponseWriter, r *http.Request) {
	fileName := filepath.Base(r.URL.Path)
	fileContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ioutil.WriteFile(filepath.Join(dir, fileName), fileContent, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
