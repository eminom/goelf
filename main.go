package main

import (
	"goelf/cipher"
	"goelf/compress"

	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"time"
)

var (
	verbose = flag.Bool("v", false, "verbose output")
	indir   = flag.String("i", "", "input directory")
	outfile = flag.String("o", "", "output file")

	intar  = flag.String("d", "", "tar to decode")
	outdir = flag.String("c", "", "folder to unfold")

	pkey = flag.String("k", "koko", "key for me")
)

func init() {
	flag.Parse()
}

func main() {
	//files := collectfiles(`C:\Users\baihai`)

	if len(*indir) != 0 {
		indirStat, err := os.Stat(*indir)
		if err != nil {
			log.Fatal(err)
		}
		if indirStat.IsDir() {
			goEncode()
		} else {
			log.Fatal("folder to compress only")
		}
		return
	}

	if len(*intar) != 0 {
		intarStat, err := os.Stat(*intar)
		if err != nil {
			log.Fatal(err)
		}
		if !intarStat.IsDir() {
			goDecode()
		} else {
			log.Fatal("error decode a folder")
		}
		return
	}

	flag.PrintDefaults()
}

func goEncode() {
	files, tr := collectfiles(*indir)
	startTime := time.Now()
	log.Printf("%v file(s) found.", len(files))
	durTime := time.Now().Sub(startTime)
	log.Printf("%v elapsed.", durTime)

	recvBuff := bytes.NewBuffer(nil)
	compress.CreateTar(files, tr, recvBuff)

	fout, err := os.Create(*outfile)
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

func goDecode() {
	chunk, err := ioutil.ReadFile(*intar)
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
	}

	inputBuff := bytes.NewBuffer(decoded)
	compress.Detar(inputBuff, *outdir)
}
