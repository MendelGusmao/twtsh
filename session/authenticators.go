package session

import (
	"strings"
	"twtsh/config"
)

func whitelist(sender string) bool {
	ok := false

	for _, whitelisted := range config.TwtSh.Whitelist {
		if whitelisted == sender {
			ok = true
		}
	}

	return ok
}

func friends(senderId int64) bool {
	ok := false

	for _, friend := range config.TwtSh.Friends {
		if friend == senderId {
			ok = true
		}
	}

	return ok
}

func password(sender, message string) bool {
	parts := strings.Split(message, " ")

	if len(parts) < 2 {
		return false
	}

	if config.TwtSh.Password != nil && parts[1] != *config.TwtSh.Password {
		return false
	}

	return true
}
