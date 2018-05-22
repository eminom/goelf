package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"./cipher"
	"./compress"
	"./locates"
	log "./logs"
)

var (
	fVerbose = flag.Bool("v", false, "verbose output")

	// default key: no key. Just do simple tar.
	pkey = flag.String("k", "", "key for me")
)

func init() {
	flag.Parse()

	log.SetVerbose(*fVerbose)
}

func main() {
	//files := collectfiles(`C:\Users\baihai`)

	if len(flag.Args()) < 2 {
		flag.PrintDefaults()
		log.Fatal("not enough parameter")
	}

	cmd := flag.Args()[0]
	input := flag.Args()[1]
	switch cmd {
	case "c":
		goEncode(input)
	case "x":
		goDecode(input)
	default:
		log.Fatalf("unsupported cmd: %v", cmd)
	}
}

func goEncode(startupDir string) {

	stat, err := os.Stat(startupDir)
	if err != nil {
		log.Fatal(err)
	}

	if !stat.IsDir() {
		log.Fatalf("error: not a directory")
	}

	startTime := time.Now()
	files, tr := locates.Collectfiles(startupDir)
	log.Printf("%v file(s) found.", len(files))
	durTime := time.Now().Sub(startTime)
	log.Printf("%v elapsed.", durTime)

	recvBuff := bytes.NewBuffer(nil)
	compress.CreateTar(files, tr, recvBuff)

	outFile := fmt.Sprintf("%v.%v.tar", filepath.Base(startupDir), time.Now().Format("15-04-05.Jan-2-2006"))

	fout, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer fout.Close()
	encrypted := recvBuff.Bytes()

	if len(*pkey) > 0 {
		var ee error
		encrypted, ee = cipher.AesEncrypt(recvBuff.Bytes(), *pkey)
		if ee != nil {
			log.Fatal(ee)
		}
	}
	fout.Write(encrypted)
}

func goDecode(targetFile string) {
	stat, err := os.Stat(targetFile)
	if err != nil {
		log.Fatal(err)
	}
	if stat.IsDir() {
		log.Fatal("this is a folder not a file(to decode)")
	}

	var outFolder string
	pos := strings.Index(targetFile, ".")
	if pos >= 0 {
		outFolder = targetFile[:pos]
	} else {
		outFolder = targetFile
	}
	outFolder += time.Now().Format("15-04-05.Jan-2-2006") + ".d"

	stat, err = os.Stat(outFolder)
	if nil == err {
		log.Fatalf("error: %v exists", outFolder)
	}

	chunk, err := ioutil.ReadFile(targetFile)
	if err != nil {
		log.Fatal(err)
	}

	var decoded []byte
	if len(*pkey) > 0 {
		var err error
		decoded, err = cipher.AesDecrypt(chunk, *pkey)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		decoded = chunk
	}

	inputBuff := bytes.NewBuffer(decoded)
	err = compress.Detar(inputBuff, outFolder)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("decoded to <%v>", outFolder)
}
