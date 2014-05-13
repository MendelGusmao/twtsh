package session

import (
	"fmt"
	"strings"
	"time"
	"twtsh/config"
)

var (
	sessions = make(map[string]time.Time)

	AuthModeWhitelist = "whitelist"
	AuthModeFriends   = "friends"
	AuthModeTwoFactor = "two-factor"
	AuthModePassword  = "password"
	AuthModeNone      = "none"

	ErrNotAuthorized       = "Not authorized"
	ErrSessionExpired      = "Session expired"
	ErrPasswordNotSupplied = "Password not supplied"
	ErrWrongPassword       = "Wrong password"
	MsgByebye              = "Bye, bye!"
	MsgFailedAuth          = "User %s failed %s authentication"
	MsgAuthOK              = "User %s successfully authenticated using %s"
)

func Authorize(sender string) (string, bool) {
	createdAt, ok := sessions[sender]

	if !ok {
		return ErrNotAuthorized, false
	}

	if time.Now().Unix() > createdAt.Add(config.TwtSh.SessionTimeout).Unix() {
		return ErrSessionExpired, false
	}

	sessions[sender] = time.Now()

	return "", true
}

func Authenticate(senderId int64, sender, message string) (string, bool) {
	ok := true
	failed := make([]string, 0)

	for _, authMode := range config.TwtSh.AuthModes {
		switch authMode {
		case AuthModeWhitelist:
			if ok = ok && whitelist(sender); !ok {
				failed = append(failed, authMode)
			}
		case AuthModeFriends:
			if ok = ok && friends(senderId); !ok {
				failed = append(failed, authMode)
			}
		case AuthModePassword:
			fmt.Println(password(sender, message))
			if ok = ok && password(sender, message); !ok {
				failed = append(failed, authMode)
			}
		case AuthModeTwoFactor:
			//
		}
	}

	authModes := strings.Join(config.TwtSh.AuthModes, "+")

	if len(config.TwtSh.AuthModes) == 0 {
		ok = true
		authModes = AuthModeNone
	}

	if !ok {
		return fmt.Sprintf(MsgFailedAuth, sender, authModes), false
	}

	sessions[sender] = time.Now()

	return fmt.Sprintf(MsgAuthOK, sender, authModes), true
}

func Exit(sender string) (string, bool) {
	delete(sessions, sender)

	return MsgByebye, true
}
