package web

import (
	"encoding/xml"
	"log/slog"
	"net/http"

	"github.com/corncobble/butterclove/config"
	"github.com/corncobble/butterclove/nftv"
	"github.com/sherif-fanous/xmltv"
)

var handler http.Handler

func SetupHandler(channels []config.Channel) {
	mux := http.NewServeMux()

	xh := epgHandler(channels)
	mux.Handle("/output/epg", xh)

	handler = mux
}

func epgHandler(channels []config.Channel) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Parse channels into xmltv.
		var tv xmltv.TV
		for _, c := range channels {
			if c.Type != config.ChannelTypeNFTV {
				continue
			}
			if err := nftv.ParseChannel(r.Context(), &tv, c); err != nil {
				slog.ErrorContext(r.Context(), "Cannot parse channel", "channel", c, "err", err)
			}
		}

		// Set headers.
		w.Header().Set("Content-Type", "application/xml")
		w.Header().Set("Content-Disposition", "attachment; filename=\"butterclove.xml\"")

		// Write generic xml header.
		w.Write([]byte(xml.Header))

		// Encode xmltv object into xml.
		e := xml.NewEncoder(w)
		defer e.Close()
		e.Indent("", "  ")
		if err := e.Encode(tv); err != nil {
			slog.ErrorContext(r.Context(), "Cannot encode xml", "err", err)
			return
		}
	}
	return http.HandlerFunc(fn)
}

// Content-Type: application/xml
// Content-Length: 2438834
// Connection: keep-alive
// Content-Disposition: attachment; filename="Dispatcharr.xml"
// Cache-Control: no-cache
// Vary: origin
// X-Frame-Options: DENY
// X-Content-Type-Options: nosniff
// Referrer-Policy: same-origin
// Cross-Origin-Opener-Policy: same-origin
