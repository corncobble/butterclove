package artiflix

import (
	"context"
	"time"

	"github.com/sherif-fanous/xmltv"
)

const (
	channelID  = "Artiflix.us"
	channelURL = "https://artiflix.com/live"
)

var lang = new("en")

func ParseChannel(ctx context.Context, tv *xmltv.TV) error {
	apiChannel, err := getAPIChannel()
	if err != nil {
		return err
	}

	// add channel
	tv.Channels = append(tv.Channels, xmltv.Channel{
		ID: channelID,
		DisplayNames: []xmltv.DisplayName{
			{Text: apiChannel.ChannelName, Lang: lang},
		},
		Icons: []xmltv.Icon{
			{Source: apiChannel.Logo},
		},
		URLs: []xmltv.URL{
			{Text: channelURL},
		},
	})

	for _, p := range []APIProgram{apiChannel.NowPlaying, apiChannel.UpNext} {
		start, err := time.Parse(timeLayout, p.StartTime)
		if err != nil {
			return err
		}
		end, err := time.Parse(timeLayout, p.EndTime)
		if err != nil {
			return err
		}

		programme := xmltv.Programme{
			Channel:      channelID,
			Titles:       []xmltv.Title{{Text: p.Title, Lang: lang}},
			Descriptions: []xmltv.Description{{Text: p.Description, Lang: lang}},
			Start:        xmltv.Time{Time: start},
			Stop:         &xmltv.Time{Time: end},
			Icons:        []xmltv.Icon{{Source: p.Thumbnail}},
		}
		tv.Programmes = append(tv.Programmes, programme)
	}

	return nil
}
