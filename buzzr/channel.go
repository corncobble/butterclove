package buzzr

import (
	"fmt"
	"time"
)

const (
	channelID         = "BUZZR.us"
	channelName       = "BUZZR"
	channelURL        = "https://buzzrtv.com/schedule"
	channelIcon       = "https://upload.wikimedia.org/wikipedia/commons/d/d6/Buzzr_logo.svg"
	channelDataURLfmt = "https://buzzrtv.com/schedule/schedule_by_date/%s/Etc/UTC"
	channelDaysOfData = 7
)

type channelData struct {
	date time.Time
	url  string
}

func channelDataAll() []channelData {
	var all []channelData
	for i := range channelDaysOfData {
		t := time.Now().UTC().AddDate(0, 0, i)
		all = append(all, channelData{
			date: t,
			url:  fmt.Sprintf(channelDataURLfmt, t.Format(time.DateOnly)),
		})
	}
	return all
}
