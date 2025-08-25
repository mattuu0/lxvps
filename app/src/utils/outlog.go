package utils

import (
	"log"
	// "path/filepath"
	"runtime"
)

func Println(data interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		// fname := filepath.Base(file)
		log.Printf("%s:%d   %s\n", file, line,data)
	}
}