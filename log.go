package boiling

import (
	"fmt"
	"log"
	"os"
)

var logger *log.Logger

func init() {
	SetLogger(log.New(os.Stderr, "boiling: ", log.LstdFlags|log.Lshortfile))
}

func SetLogger(l *log.Logger) {
	logger = l
}

func LogErrf(s string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Output(2, fmt.Sprintf(s, args...))
}
