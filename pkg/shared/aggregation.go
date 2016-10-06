package shared

import "bytes"

// PrepareEmail inserts \n after every 500 characters
func PrepareEmail(email *Email) error {
	email.Body = InsertNth(email.Body, 500, "\r\n")

	return nil
}

// InsertNth inserts rune after each N characters
func InsertNth(s string, n int, insert string) string {
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
