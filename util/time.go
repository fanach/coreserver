package util

import (
	"log"
	"time"
)

const (
	// TimeLayout1 ...
	TimeLayout1 = "2006-01-02 15:04:05"
	// TimeLayout2 ...
	TimeLayout2 = "20060102150405"
)

// FormatTime formats date, time to YYYY-MM-DD hh:mm:ss
func FormatTime(t time.Time, layout string) string {
	return t.Format(layout)
}

// ParseTime parse string "YYYY-MM-DD hh:mm:ss" to time
func ParseTime(date, layout string) (t time.Time, err error) {
	t, err = time.Parse(layout, date)
	if err != nil {
		log.Printf("parse date %s error: %v\n", date, err)
		return
	}
	return
}
