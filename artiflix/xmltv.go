package artiflix

import (
	"context"
	"time"

	"github.com/sherif-fanous/xmltv"
)

const (
	channelID  = "Artiflix.us"
	channelURL = "https://artiflix.com/live"
	timeLayout = "2006-01-02T15:04:05.000Z"
)

var (
	client = newAPIClient()
	lang   = new("en")
)

func ParseChannel(ctx context.Context, tv *xmltv.TV) error {
	channel, err := client.getChannel(ctx)
	if err != nil {
		return err
	}

	// add channel
	tv.Channels = append(tv.Channels, xmltv.Channel{
		ID: channelID,
		DisplayNames: []xmltv.DisplayName{
			{Text: channel.ChannelName, Lang: lang},
		},
		Icons: []xmltv.Icon{
			{Source: channel.Logo},
		},
		URLs: []xmltv.URL{
			{Text: channelURL},
		},
	})

	for _, p := range []Program{channel.NowPlaying, channel.UpNext} {
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
