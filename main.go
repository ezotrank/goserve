package main

import (
	"flag"
	"net/http"
	"os"
	"log"
	"io"
	"path"
)

var (
	bind = flag.String("bind", "0.0.0.0:8080", "bind address 0.0.0.0:8080")
	folder = flag.String("folder", "files", "root folder")
	token = flag.String("token", "test-token", "token for auth")
)

func handler(w http.ResponseWriter, r *http.Request) {
	currentDir, _ := os.Getwd()
	switch r.Method {
	case "GET":
		fp := path.Join(currentDir, *folder, r.URL.Path[1:])
		http.ServeFile(w, r, fp)
	case "POST":
		if r.URL.Query().Get("token") != *token  {
			w.WriteHeader(401)
			return
		}
		fp := path.Join(currentDir, *folder, r.URL.Path[1:])
		file, err := os.OpenFile(fp, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			w.WriteHeader(500)
		}
		defer file.Close()
		_, err = io.Copy(file, r.Body)
		if err != nil {
			w.WriteHeader(500)
			log.Println(err)
		}
	}
}

func main() {
	flag.Parse()

	http.HandleFunc("/", handler)
	http.ListenAndServe(*bind, nil)
}