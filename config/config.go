package config

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

// Web contains options for configuring the web server.
type Web struct {
	// Host is the address to listen on.
	Host string `json:"host"`
	// Port is the address port to listen on.
	Port int `json:"port"`
}

// ChannelType specifies how to parse channel information.
type ChannelType string

const (
	// ChannelTypeNFTV is for parsing NFTV channel information.
	ChannelTypeNFTV ChannelType = "nftv"
	// ChannelTypeBuzzr is for parsing Buzzr channel information.
	ChannelTypeBuzzr ChannelType = "buzzr"
)

// Channel specifies which channel(s) to include in xmltv output.
type Channel struct {
	Type ChannelType `json:"type"`
	ID   string      `json:"id"`
}

// Config represents the configuration for the application.
type Config struct {
	Web      Web       `json:"web"`
	Channels []Channel `json:"channels"`
	Dir      string    `json:"-"`
}

// NewConfig returns a new (default) config instance.
func NewConfig() Config {
	ep, err := os.Executable()
	if err != nil {
		log.Fatalf("new config error: %s", err.Error())
	}
	return Config{
		Web: Web{
			Host: "",
			Port: 7590,
		},
		Dir: filepath.Join(filepath.Dir(ep), "config"),
	}
}

func (c Config) Write() error {
	buf, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(c.path(), buf, fs.FileMode(0644)); err != nil {
		return err
	}
	return nil
}

func (c *Config) Read() error {
	buf, err := os.ReadFile(c.path())
	if err != nil {
		return err
	}
	if err := json.Unmarshal(buf, c); err != nil {
		return err
	}
	return nil
}

func (c Config) path() string {
	return filepath.Join(c.Dir, "config.json")
}
