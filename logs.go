package main

import (
	rlog "log"
)

type logger struct{}

func init() {
	rlog.SetFlags(rlog.Lshortfile)
}

func (*logger) Printf(format string, a ...interface{}) {
	if *verbose {
		rlog.Printf(format, a...)
	}
}

func (*logger) Fatal(a ...interface{}) {
	rlog.Fatal(a...)
}

var (
	log logger
)
