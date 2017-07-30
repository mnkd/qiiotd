package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

var (
	ErrInvalidJson = errors.New("ErrInvalidJson")
)

type Config struct {
	Qiita struct {
		Domain      string `json:"domain"`
		AccessToken string `json:"access_token"`
		PerPage     int    `json:"per_page"`
	} `json:"qiita"`
	SlackWebhooks []struct {
		Channel    string `json:"channel"`
		IconEmoji  string `json:"icon_emoji"`
		Username   string `json:"username"`
		WebhookUrl string `json:"webhook_url"`
	} `json:"slack_webhooks"`
}

func (c *Config) validate() error {
	if len(c.Qiita.Domain) == 0 {
		fmt.Fprintln(os.Stderr, "Invalid config.json. You should set domain.")
		return ErrInvalidJson
	}
	if len(c.Qiita.AccessToken) == 0 {
		fmt.Fprintln(os.Stderr, "Invalid config.json. You should set access_token.")
		return ErrInvalidJson
	}
	return nil
}

func NewConfig(path string) (Config, error) {
	var config Config

	usr, err := user.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not get current user.", err)
		return config, err
	}

	if len(path) == 0 {
		path = filepath.Join(usr.HomeDir, "/.config/qiotd/config.json")
	} else {
		p, err := filepath.Abs(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "[Error] Could not return absolute representation of path:", err, path)
			return config, err
		}
		path = p
	}

	str, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not read config.json. ", err)
		return config, err
	}

	if err := json.Unmarshal(str, &config); err != nil {
		fmt.Fprintln(os.Stderr, "JSON Unmarshal Error:", err)
		return config, err
	}

	if err = config.validate(); err != nil {
		return config, err
	}

	return config, nil
}
