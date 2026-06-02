package nftv

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/corncobble/butterclove/config"
	"github.com/sherif-fanous/xmltv"
	"golang.org/x/net/html"
)

var (
	lang                   = new("en")
	attrProgram            = html.Attribute{Key: "class", Val: "program"}
	attrProgramTitle       = html.Attribute{Key: "class", Val: "program-title"}
	attrProgramDescription = html.Attribute{Key: "class", Val: "program-description"}
)

func ParseChannel(ctx context.Context, tv *xmltv.TV, channel config.Channel) error {
	c := getChannelByID(channel.ID)

	resp, err := http.Get(c.url())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// add channel
	tv.Channels = append(tv.Channels, xmltv.Channel{
		ID: c.id(),
		DisplayNames: []xmltv.DisplayName{
			{Text: c.name(), Lang: lang},
		},
		Icons: []xmltv.Icon{
			{Source: c.icon()},
		},
		URLs: []xmltv.URL{
			{Text: c.url()},
		},
	})

	z := html.NewTokenizer(resp.Body)
	var (
		tt        html.TokenType
		t         html.Token
		programme xmltv.Programme
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
			case attrProgram:

				for _, attr := range t.Attr {
					switch attr.Key {

					// start of a new program
					case "class":
						programme = xmltv.Programme{Channel: c.id()}

					// Val is in unix time (msec)
					case "data-start-time":
						unix, err := strconv.ParseInt(attr.Val, 10, 0)
						if err != nil {
							slog.ErrorContext(ctx, "Unable to parse int", "val", attr.Val, "err", err)
							break
						}
						programme.Start.Time = time.UnixMilli(unix).UTC()

					// Val is in unix time (msec)
					case "data-end-time":
						unix, err := strconv.ParseInt(attr.Val, 10, 0)
						if err != nil {
							slog.ErrorContext(ctx, "Unable to parse int", "val", attr.Val, "err", err)
							break
						}
						programme.Stop = new(xmltv.Time)
						programme.Stop.Time = time.UnixMilli(unix).UTC()

					// Val is a url string
					case "data-thumbnail":
						programme.Icons = []xmltv.Icon{
							{Source: attr.Val},
						}
					}
				}

			// program title is contained in subsequent token
			case attrProgramTitle:
				if programme.Channel == "" {
					slog.InfoContext(ctx, "Program title attribute found with no associated program, ignoring.")
					break
				}

				next()
				if tt != html.TextToken {
					slog.WarnContext(ctx, "No program title found", "expected", html.TextToken.String(), "got", tt.String())
					break
				}
				programme.Titles = []xmltv.Title{
					{Text: string(z.Text()), Lang: lang},
				}

			// program description is contained in subsequent token
			case attrProgramDescription:
				if programme.Channel == "" {
					slog.InfoContext(ctx, "Program title attribute found with no associated program, ignoring.")
					break
				}

				next()
				if tt != html.TextToken {
					slog.WarnContext(ctx, "No program description found", "expected", html.TextToken.String(), "got", tt.String())
					break
				}
				programme.Descriptions = []xmltv.Description{
					{Text: string(z.Text()), Lang: lang},
				}

				// append complete program entry
				tv.Programmes = append(tv.Programmes, programme)
				// set current program to zero value
				programme = xmltv.Programme{}
			}
		}
	}
}
