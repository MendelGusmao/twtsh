package config

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

var (
	TwtSh TwtShConf
)

type TwtShConf struct {
	OAuth          OAuthConf
	AuthModes      []string      `json:"auth-modes"`
	Whitelist      []string      `json:"whitelist"`
	Friends        []int64       `json:"-"`
	LogLevel       *string       `json:"log-level"`
	Password       *string       `json:"password"`
	SessionTimeout time.Duration `json:"-"`
	sessionTimeout *string       `json:"session-timeout"`
}

type OAuthConf struct {
	ConsumerKey    string `json:"consumer-key"`
	ConsumerSecret string `json:"consumer-secret"`
	Token          string `json:"token"`
	TokenSecret    string `json:"token-secret"`
}

func init() {
	TwtSh.SessionTimeout, _ = time.ParseDuration("5m")
}

func Load(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &TwtSh)
	if err != nil {
		return err
	}

	if TwtSh.sessionTimeout != nil {
		timeout, err := time.ParseDuration(*TwtSh.sessionTimeout)
		if err == nil {
			return err
		}

		TwtSh.SessionTimeout = timeout
	}

	return nil
}
