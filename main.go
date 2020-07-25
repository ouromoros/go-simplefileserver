package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"

	// "html/template"
	"flag"
	"net/http"
	"os"
)

func main() {
	fdir := flag.String("dir", ".", "The directory the file server serves.")
	fport := flag.Int("port", 8888, "The port number the file server uses.")
	flag.Parse()

	dir, _ := filepath.Abs(*fdir)
	port := *fport

	log.Printf("Serving directory %v at port %v...", dir, port)

	serve(dir, port)
}

type fileServerHandler struct {
	basedir string
}

func (s fileServerHandler) ServeHTTP(respw http.ResponseWriter, req *http.Request) {
	path := filepath.Join(s.basedir, req.URL.Path)

	// info, err := os.Stat(path)
	_, err := os.Stat(path)
	respw.Header().Add("Cache-Control", "no-cache")
	log.Println("Serving " + path)
	if os.IsNotExist(err) {
		fmt.Fprintf(respw, "<h1>404</h1>")
		// } else if info.IsDir() {
		// 	fmt.Fprintf(respw, pagedir(path))
	} else {
		http.ServeFile(respw, req, path)
	}
}

func serve(dir string, port int) {
	handler := fileServerHandler{dir}
	s := http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: handler,
	}
	log.Fatal(s.ListenAndServe())
}
