package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/araddon/httpstream"
	"github.com/mrjones/oauth"
	"log"
	"os"
	"strings"
	"twtsh/config"
	"twtsh/shell"
	"twtsh/types"
)

var (
	configFilename *string = flag.String("config", "/etc/twtsh", "Configuration file")
	logLevel       *string = flag.String("log-level", "info", "Which log level: [debug,info,warn,error,fatal]")

	directMessage = []byte(`{"direct_message"`)
	friends       = []byte(`{"friends"`)
)

func main() {
	flag.Parse()

	if err := config.Load(*configFilename); err != nil {
		fmt.Printf("Error loading configuration file %s: %s\n", *configFilename, err)
		os.Exit(1)
	}

	if logLevel != nil {
		config.TwtSh.LogLevel = logLevel
	}

	httpstream.SetLogger(log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile), *logLevel)

	httpstream.OauthCon = oauth.NewConsumer(
		config.TwtSh.OAuth.ConsumerKey,
		config.TwtSh.OAuth.ConsumerSecret,
		oauth.ServiceProvider{},
	)

	anaconda.SetConsumerKey(config.TwtSh.OAuth.ConsumerKey)
	anaconda.SetConsumerSecret(config.TwtSh.OAuth.ConsumerSecret)

	httpstream.Log(httpstream.INFO, "Starting TwtSh")
	httpstream.Log(httpstream.INFO, "Authentication Modes:", strings.Join(config.TwtSh.AuthModes, ", "))
	httpstream.Log(httpstream.INFO, "Whitelist:", strings.Join(config.TwtSh.Whitelist, ", "))

	watchStream()
}

func watchStream() {
	stream := make(chan []byte, 1000)
	done := make(chan bool)

	token := oauth.AccessToken{
		Token:  config.TwtSh.OAuth.Token,
		Secret: config.TwtSh.OAuth.TokenSecret,
	}

	api := anaconda.NewTwitterApi(
		token.Token,
		token.Secret,
	)

	client := httpstream.NewOAuthClient(&token, func(line []byte) {
		stream <- line
	})

	err := client.User(done)

	if err != nil {
		httpstream.Log(httpstream.ERROR, err.Error())
		os.Exit(1)
	}

	go func() {
		for body := range stream {
			httpstream.Log(httpstream.DEBUG, "json:", string(body))

			switch {
			case bytes.HasPrefix(body, directMessage):
				object := make(map[string]types.DirectMessage)

				if err := json.Unmarshal(body, &object); err != nil {
					httpstream.Log(httpstream.ERROR, err.Error())
					continue
				}

				directMessage := object["direct_message"]
				httpstream.Log(httpstream.INFO, "Received", directMessage.Text, "from", directMessage.SenderScreenName)

				msg, _ := shell.Handle(directMessage)
				tweet := fmt.Sprintf("d %s %s", directMessage.SenderScreenName, msg)
				if _, err := api.PostTweet(tweet, nil); err != nil {
					httpstream.Log(httpstream.ERROR, err.Error())
				}

			case bytes.HasPrefix(body, friends):
				object := make(map[string][]int64)

				if err := json.Unmarshal(body, &object); err != nil {
					httpstream.Log(httpstream.ERROR, err.Error())
					continue
				}

				config.TwtSh.Friends = object["friends"]
			default:
				httpstream.Log(httpstream.DEBUG, "Unrecognized JSON", string(body))
			}

		}
	}()

	<-done
}
