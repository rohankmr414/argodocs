package logger

import (
	"log"
	"os"
)

func GetLogger(prefix string) *log.Logger {
	return log.New(os.Stdout, prefix, log.LstdFlags)
}
