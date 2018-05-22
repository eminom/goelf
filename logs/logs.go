package logs

import (
	"fmt"
	"os"
	"sync"

	rlog "log"
)

type logger struct {
	mu sync.Mutex
}

var (
	fVerbose bool
)

func SetVerbose(enabled bool) {
	fVerbose = enabled
}

func init() {
	rlog.SetFlags(rlog.Lshortfile)
}

func (*logger) Printf(format string, a ...interface{}) {
	if fVerbose {
		rlog.Output(3, fmt.Sprintf(format, a...))
	}
}

func (*logger) Fatal(i interface{}) {
	rlog.Output(3, fmt.Sprintf("%v", i))
	os.Exit(1)
}

func (*logger) Fatalf(f string, a ...interface{}) {
	rlog.Output(3, fmt.Sprintf(f, a...))
	os.Exit(1)
}

var (
	esteLog logger
)

// wrappers to go
func Printf(a string, p ...interface{}) {
	esteLog.Printf(a, p...)
}

func Fatal(i interface{}) {
	esteLog.Fatal(i)
}

func Fatalf(f string, a ...interface{}) {
	esteLog.Fatalf(f, a...)
}
