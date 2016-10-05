package shared

import "log"

// LogErr func
func LogErr(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}
