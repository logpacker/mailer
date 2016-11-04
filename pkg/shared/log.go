package shared

import (
	"log"
)

// LogErr func
func LogErr(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

// Logf func
func Logf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
