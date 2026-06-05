package artiflix

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type APIChannel struct {
	ChannelName string     `json:"channel_name"`
	LiveURL     string     `json:"live_url"`
	Logo        string     `json:"logo"`
	NowPlaying  APIProgram `json:"now_playing"`
	UpNext      APIProgram `json:"up_next"`
}

type APIProgram struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Thumbnail   string `json:"thumbnail"`
}

const (
	timeLayout = "2006-01-02T15:04:05.000Z"
)

var header = http.Header{
	"uid":   []string{"6790380"},
	"pubid": []string{"50117"},
}

func getAPIChannel() (APIChannel, error) {
	var data APIChannel

	token, err := getAuthToken()
	if err != nil {
		return data, err
	}

	r, err := http.NewRequest(http.MethodGet, "https://api.gizmott.com/api/v1/fastchannel/1159", nil)
	if err != nil {
		return data, err
	}
	r.Header = header
	r.Header.Add("access-token", token)

	w, err := http.DefaultClient.Do(r)
	if err != nil {
		return data, err
	}
	defer w.Body.Close()

	body := struct {
		Data []APIChannel `json:"data"`
	}{}
	dec := json.NewDecoder(w.Body)
	if err = dec.Decode(&body); err != nil {
		return data, err
	}

	if len(body.Data) < 1 {
		return data, fmt.Errorf("expected data length > 0, got 0")
	}
	return body.Data[0], nil
}

func getAuthToken() (string, error) {
	r, err := http.NewRequest(http.MethodGet, "https://api.gizmott.com/api/v1/account/authenticate", nil)
	if err != nil {
		return "", err
	}
	r.Header = header

	w, err := http.DefaultClient.Do(r)
	if err != nil {
		return "", err
	}
	defer w.Body.Close()

	data := struct {
		Token string `json:"token"`
	}{}
	dec := json.NewDecoder(w.Body)
	if err = dec.Decode(&data); err != nil {
		return "", err
	}

	return data.Token, nil
}
