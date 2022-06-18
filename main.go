package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
)

const BufferSize = 5

type aContainsFile struct {
	http.File
}

type aList struct {
	index int
	fname string
	file  fs.FileInfo
}

func (f aContainsFile) Readdir(n int) (fis []fs.FileInfo, err error) {
	var wg sync.WaitGroup
	var fileList []aList
	// min := make(chan bool)
	var min int
	min = math.MaxInt32

	files, err := f.File.Readdir(n)
	if err != nil {
		log.Fatal(err)
	}
	wg.Add(len(files))
	for _, file := range files {
		go readByChunkAndFindA(&fis, file, file.Name(), &wg, &fileList, &min)
	}
	wg.Wait()

	fmt.Println("fileList: ", fileList)
	if len(fileList) == 0 {
		return nil, errors.New("no files found")
	}

	fmt.Println("min :", min)

	for _, file := range fileList {
		if file.index == min {
			fis = append(fis, fs.FileInfo(file.file))
		}
	}

	return fis, nil
}

func readByChunkAndFindA(fis *[]fs.FileInfo, file fs.FileInfo, filename string, wg *sync.WaitGroup, fileList *[]aList, min *int) {
	defer wg.Done()

	f, err := os.OpenFile(fmt.Sprintf("temp/%s", filename), os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	buffer := make([]byte, BufferSize)

	counter := 0
	for {
		bytesread, err := f.Read(buffer)

		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}

		if *min < BufferSize*counter {
			return
		}

		if strings.Contains(string(buffer[:bytesread]), "a") {
			idx := strings.Index(string(buffer[:bytesread]), "a") + counter*BufferSize
			if idx < *min {
				*min = idx
				fmt.Println("idx:", idx, " min:", *min)
			}
			*fileList = append(*fileList, aList{idx, filename, file})
			return
		}

		counter += 1
	}
}

type aContainsFileSystem struct {
	http.FileSystem
}

func (fsys aContainsFileSystem) Open(name string) (http.File, error) {

	file, err := fsys.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	return aContainsFile{file}, err
}

func main() {
	fsys := aContainsFileSystem{http.Dir("./temp")}
	http.Handle("/", http.FileServer(fsys))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
