package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jhoonb/archivex"
)

func downloadFromUrl(url, dir string, throttle chan bool, wg *sync.WaitGroup) {
	defer func() {
		<-throttle
		wg.Done()
	}()
	tokens := strings.Split(url, "/")
	path := filepath.Join(dir, tokens[len(tokens)-1])
	fmt.Println("Downloading", url, "to", path)

	output, err := os.Create(path)
	if err != nil {
		fmt.Println("Error while creating", path, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}

	fmt.Println(n, "bytes downloaded.")
}

func zip(dst, src string) {
	zip := new(archivex.ZipFile)
	zip.Create(dst)
	zip.AddAll(src, false)
	zip.Close()
}

func Download(name string, urls []string) error {
	workDir, err := ioutil.TempDir("", "downloads")
	if err != nil {
		return err
	}
	throttle := make(chan bool, 5)
	var wg sync.WaitGroup
	for _, url := range urls {
		throttle <- true
		wg.Add(1)
		go downloadFromUrl(url, workDir, throttle, &wg)
	}
	wg.Wait()

	zip(name, workDir)
	return nil
}
