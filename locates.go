package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

func walkdir(dirnow string, outname chan<- string, wg *sync.WaitGroup, resGroup chan struct{}) {
	resGroup <- struct{}{} // If full, must be blocked

	go func() {
		defer func() {
			<-resGroup // release the resource
			wg.Done()
		}()
		din, err := os.Open(dirnow)
		if err != nil {
			log.Printf("error open: %v", err)
			return
		}
		defer din.Close()
		files, err := din.Readdir(-1)
		if err != nil {
			log.Printf("error readdir: %v", err)
			return
		}
		for _, file := range files {
			if file.IsDir() {
				wg.Add(1)
				go walkdir(path.Join(dirnow, file.Name()), outname, wg, resGroup)
			} else {
				outname <- path.Join(dirnow, file.Name())
			}
		}
	}()
}

func collectfiles(startd string) ([]string, func(string) string) {
	nd, err := filepath.Abs(startd)
	if err != nil {
		panic(fmt.Errorf("error filepath.Abs: %v", err))
	}
	nd = filepath.ToSlash(nd)
	startd = nd

	fi, err := os.Stat(startd)
	if err != nil {
		panic(fmt.Errorf("error: %v", err))
	}
	if !fi.IsDir() {
		panic(fmt.Errorf("no a directory"))
	}

	var fs []string

	nombresDeFiles := make(chan string, 1024)
	bufferedchan := make(chan struct{}, runtime.GOMAXPROCS(runtime.NumCPU()))

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go walkdir(startd, nombresDeFiles, &waitGroup, bufferedchan)

	terminarSig := make(chan bool)
	var wg1 sync.WaitGroup
	wg1.Add(1)
	go func() {
		defer wg1.Done()
		for {
			select {
			case <-terminarSig:
				return
			case name := <-nombresDeFiles:
				fs = append(fs, name)
			}
		}
	}()
	waitGroup.Wait()

	close(terminarSig)
	wg1.Wait()

	//Drain
A100:
	for {
		select {
		case name := <-nombresDeFiles:
			fs = append(fs, name)
		default:
			break A100
		}
	}

	lead := fi.Name()
	//log.Printf("origin lead: <%v> now is: <%v>", startd, lead)
	return fs, func(name string) string {
		return filepath.ToSlash(path.Join(lead, strings.TrimPrefix(name, startd)))
	}
}
