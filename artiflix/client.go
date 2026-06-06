package artiflix

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type Channel struct {
	ChannelName string  `json:"channel_name"`
	LiveURL     string  `json:"live_url"`
	Logo        string  `json:"logo"`
	NowPlaying  Program `json:"now_playing"`
	UpNext      Program `json:"up_next"`
}

type Program struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Thumbnail   string `json:"thumbnail"`
}

type authResponse struct {
	Token string `json:"token"`
}

type channelResponse struct {
	Data    []Channel `json:"data"`
	Success bool      `json:"success"`
	Message string    `json:"message"`
}

type apiClient struct {
	client *http.Client
	header http.Header

	mu        sync.RWMutex
	token     string
	expiresAt time.Time
}

func newAPIClient() *apiClient {
	header := http.Header{
		"uid":          []string{"6790380"},
		"pubid":        []string{"50117"},
		"country_code": []string{"US"},
		"device_type":  []string{"web"},
		"dev_id":       []string{"756e9bd72db9a7f3feeb3de4a8e5f291"},
		"channelid":    []string{"451"},
		"ip":           []string{"0.0.0.0"},
	}
	return &apiClient{
		client: &http.Client{Timeout: 10 * time.Second},
		header: header,
	}
}

func (c *apiClient) getChannel(ctx context.Context) (Channel, error) {
	var channel Channel

	t, err := c.getToken(ctx)
	if err != nil {
		return channel, err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.gizmott.com/api/v1/fastchannel/1159", nil)
	if err != nil {
		return channel, err
	}
	r.Header = c.header
	r.Header.Set("access-token", t)

	w, err := c.client.Do(r)
	if err != nil {
		return channel, err
	}
	defer w.Body.Close()

	var data channelResponse
	dec := json.NewDecoder(w.Body)
	if err = dec.Decode(&data); err != nil {
		return channel, err
	}

	if !data.Success {
		slog.WarnContext(ctx, "Unable to get channel data, clearing token (try again)", "uri", r.RequestURI, "message", data.Message)
		c.clearToken()
		return channel, errors.New(data.Message)
	}
	if len(data.Data) < 1 {
		return channel, fmt.Errorf("expected data length > 0, got 0")
	}
	return data.Data[0], nil
}

func (c *apiClient) authenticate(ctx context.Context) error {
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.gizmott.com/api/v1/account/authenticate", nil)
	if err != nil {
		return err
	}
	r.Header = c.header

	w, err := c.client.Do(r)
	if err != nil {
		return err
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusOK {
		return fmt.Errorf("authentication failed, status %d", w.StatusCode)
	}

	var data authResponse
	dec := json.NewDecoder(w.Body)
	if err = dec.Decode(&data); err != nil {
		return err
	}

	c.mu.Lock()
	c.token = data.Token
	// TODO: Is 80 days adequate here?
	c.expiresAt = time.Now().Add(time.Hour * 24 * 80)
	c.mu.Unlock()
	return nil
}

func (c *apiClient) getToken(ctx context.Context) (string, error) {
	c.mu.RLock()
	t := c.token
	exp := c.expiresAt
	c.mu.RUnlock()

	if t != "" && time.Now().Before(exp) {
		return t, nil
	}
	if err := c.authenticate(ctx); err != nil {
		return "", err
	}

	c.mu.RLock()
	t = c.token
	c.mu.RUnlock()
	return t, nil
}

func (c *apiClient) clearToken() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.token = ""
}
