package session

import (
	"testing"
	"time"
	"twtsh/config"
)

var (
	user       = "user"
	pw         = "abc123"
	timeout, _ = time.ParseDuration("5m")
)

func init() {
	config.TwtSh.Password = &pw
	config.TwtSh.SessionTimeout = timeout
}

func TestAuthorizeWithNonExistentUser(t *testing.T) {
	if msg, ok := Authorize(user); ok || msg != ErrNotAuthorized {
		t.Fatal()
	}
}

func TestAuthorizeWithExistentUser(t *testing.T) {
	sessions[user] = time.Now()

	if msg, ok := Authorize(user); !ok {
		t.Fatal(msg)
	}
}

func TestAuthorizeWithExpiredSession(t *testing.T) {
	sessions[user] = time.Now().Add(-timeout * 2)

	if msg, ok := Authorize(user); ok || msg != ErrSessionExpired {
		t.Fatal()
	}
}
