package util

import (
	"testing"
	"time"

	"github.com/bmizerany/assert"
)

func TestFormatTime(t *testing.T) {
	var cases = []struct {
		t         time.Time
		layout    string
		expectLen int
	}{
		{time.Now(), TimeLayout1, len(TimeLayout1)},
		{time.Now(), TimeLayout2, len(TimeLayout2)},
	}

	for _, c := range cases {
		got := FormatTime(c.t, c.layout)
		assert.Equal(t, c.expectLen, len(got))
	}
}

func TestParseTime(t *testing.T) {
	var cases = []struct {
		timeStr    string
		layout     string
		expectTime time.Time
		expectErr  error
	}{
		{"2017-03-01 16:40:28", TimeLayout1, time.Date(2017, 3, 1, 16, 40, 28, 0, time.Local), nil},
		{"20170301164028", TimeLayout2, time.Date(2017, 3, 1, 16, 40, 28, 0, time.Local), nil},
	}

	for _, c := range cases {
		_, gotErr := ParseTime(c.timeStr, c.layout)
		assert.Equal(t, c.expectErr, gotErr)
		// fmt.Printf("%+v\n", c.expectTime)
		// fmt.Printf("%+v\n", gotTime)
		// eq := reflect.DeepEqual(c.expectTime, gotTime)
		// if !eq {
		// 	t.Fail()
		// }
		// assert.Equal(t, c.expectTime, gotTime)
	}
}
