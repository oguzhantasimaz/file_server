package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	port := flag.String("p", "8100", "port to serve on")
	flag.Parse()

	// log.Println(fileNames)
	http.HandleFunc("/", home)
	http.HandleFunc("/download", download)
	// http.HandleFunc("/test", tst)

	// log.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
	s := &http.Server{
		Addr:           ":" + *port,
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}

func home(w http.ResponseWriter, r *http.Request) {
	directory := flag.String("d", "../temp", "the directory of static file to host")

	http.FileServer(http.Dir(*directory)).ServeHTTP(w, r)
}

// func tst(w http.ResponseWriter, r *http.Request) {
// 	directory := fmt.Sprintf("../temp")

// 	// http.StripPrefix(directory, http.FileServer(http.Dir(directory)))
// 	//download files in directory with FileServer
// 	http.FileServer()
// }

func download(w http.ResponseWriter, r *http.Request) {
	directory := fmt.Sprintf("../temp")

	fileNames := getFileNames(directory)
	if len(fileNames) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for _, name := range fileNames {
		w.Header().Set("Content-Disposition", "attachment; filename="+name)
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

		f, err := os.OpenFile(fmt.Sprintf("%s/%s", directory, name), os.O_RDONLY, 0666)
		fmt.Println(fmt.Sprintf("%s/%s", directory, name))
		if err != nil {
			log.Fatal(err)
		}
		http.ServeFile(w, r, fmt.Sprintf("%s/%s", directory, name))
		defer f.Close()
	}
}

func getFileNames(dir string) []string {
	fileNames := make([]string, 0)

	f, err := os.Open(dir)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	files, err := f.Readdir(0)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	for _, v := range files {
		fileNames = append(fileNames, v.Name())
	}
	return fileNames
}
