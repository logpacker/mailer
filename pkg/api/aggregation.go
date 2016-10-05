package api

import (
	"bytes"
	"github.com/logpacker/mailer/pkg/shared"
)

// Insert \n after every 500 characters
func prepareEmail(email *shared.Email) error {
	email.Body = insertNth(email.Body, 500, "\r\n")

	return nil
}

// insertNth inserts rune after each N characters
func insertNth(s string, n int, insert string) string {
	var buffer bytes.Buffer
	var prevN = n - 1
	var prevL = len(s) - 1

	for i := range s {
		buffer.WriteString(string(s[i]))
		if i%n == prevN && i != prevL {
			buffer.WriteString(insert)
		}
	}

	return buffer.String()
}
