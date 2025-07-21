package main

import (
	"testing"
	"time"

	"snippetbox.gregor-pifko/internal/assert"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			"UTC",
			time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC),
			"17 Mar 2022 at 10:15",
		},
		{
			"Empty",
			time.Time{},
			"",
		},
		{
			"CET",
			time.Date(2022, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			"17 Mar 2022 at 09:15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			assert.Equal(t, hd, tt.want)
		})
	}
}
