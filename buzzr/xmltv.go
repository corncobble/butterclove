package buzzr

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/sherif-fanous/xmltv"
	"golang.org/x/net/html"
)

var (
	lang         = new("en")
	attrTime     = html.Attribute{Key: "class", Val: "time"}
	attrTitle    = html.Attribute{Key: "class", Val: "title"}
	attrSubTitle = html.Attribute{Key: "class", Val: "sub_title"}
)

func ParseChannel(ctx context.Context, tv *xmltv.TV) error {
	// add channel
	tv.Channels = append(tv.Channels, xmltv.Channel{
		ID: channelID,
		DisplayNames: []xmltv.DisplayName{
			{Text: channelName, Lang: lang},
		},
		Icons: []xmltv.Icon{
			{Source: channelIcon},
		},
		URLs: []xmltv.URL{
			{Text: channelURL},
		},
	})

	var programme xmltv.Programme

	parseChannelData := func(cd channelData) error {
		resp, err := http.Get(cd.url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		z := html.NewTokenizer(resp.Body)
		var (
			tt html.TokenType
			t  html.Token
		)

		next := func() {
			tt = z.Next()
		}

		for {
			next()
			switch tt {

			case html.ErrorToken:
				if !errors.Is(z.Err(), io.EOF) {
					return fmt.Errorf("error during tokenization: %s", z.Err())
				}
				return nil

			case html.StartTagToken:
				t = z.Token()
				if len(t.Attr) == 0 {
					continue
				}
				attr := t.Attr[0]
				switch attr {

				// start of a new program
				case attrTime:

					// time is contained in subsequent (text) token
					next()
					if tt != html.TextToken {
						slog.WarnContext(ctx, "No time found", "expected", html.TextToken.String(), "got", tt.String())
						break
					}

					// parse time token (3:04pm -> 3:04PM)
					text := strings.TrimSpace(string(z.Text()))
					ktime, err := time.Parse(time.Kitchen, strings.ToUpper(text))
					if err != nil {
						slog.WarnContext(ctx, "Unable to parse time token", "text", text, "err", err)
						break
					}
					// create programme date using channel data date and parsed time
					date := time.Date(cd.date.Year(), cd.date.Month(), cd.date.Day(), ktime.Hour(), ktime.Minute(), ktime.Second(), 0, time.UTC)

					// programme contains start time, treat current time token as stop time and start time for next programme
					if !programme.Start.IsZero() {
						programme.Stop = &xmltv.Time{Time: date}
						tv.Programmes = append(tv.Programmes, programme)
					}
					programme = xmltv.Programme{
						Channel: channelID,
						Start:   xmltv.Time{Time: date},
					}

				// program title is in next attribute (e.g. data="Card Sharks 86_0258")
				case attrTitle:
					if len(t.Attr) < 2 {
						slog.WarnContext(ctx, "No data attribute found for program title", "expected length", 2, "got", len(t.Attr))
						break
					}
					// format data value into program title
					title, episodeNum, found := strings.Cut(t.Attr[1].Val, "_")
					if found {
						programme.EpisodeNumbers = []xmltv.EpisodeNumber{{Text: episodeNum}}
					}
					programme.Titles = []xmltv.Title{
						{Text: title, Lang: lang},
					}

				// program description is contained in subsequent (text) token
				case attrSubTitle:
					next()
					if tt != html.TextToken {
						slog.WarnContext(ctx, "No program description found", "expected", html.TextToken.String(), "got", tt.String())
						break
					}
					programme.Descriptions = []xmltv.Description{
						{Text: string(z.Text()), Lang: lang},
					}
				}
			}
		}
	}

	all := channelDataAll()
	for _, cd := range all {
		if err := parseChannelData(cd); err != nil {
			return err
		}
	}
	return nil
}
