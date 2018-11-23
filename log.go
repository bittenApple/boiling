package boiling

import (
	"fmt"
	"log"
	"os"
)

var logger *log.Logger

func init() {
	setLogger(log.New(os.Stderr, "boiling: ", log.LstdFlags|log.Lshortfile))
}

func setLogger(l *log.Logger) {
	logger = l
}

func logErrf(s string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Output(2, fmt.Sprintf(s, args...))
}
