package shared

import (
	"log"
	"time"
)

// LogErr func
func LogErr(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

// Logf func
func Logf(format string, v ...interface{}) {
	log.Printf("["+time.Now().Format(time.RFC822)+"]: "+format, v...)
}
