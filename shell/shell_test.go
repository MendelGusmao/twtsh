package shell

import (
	"encoding/json"
	"fmt"
	"testing"
	"twtsh/config"
	"twtsh/session"
	"twtsh/types"
)

const (
	messageTemplate = `
	{
		"id": 0,
		"text": "%s",
		"sender_id": 0,
		"sender_screen_name": "test",
		"created_at": ""
	}`
)

var (
	loginMessage           = createMessage("!login abc123")
	loginMessageNoPassword = createMessage("!login")
)

func init() {
	password := "abc123"

	config.TwtSh.Password = &password
	config.TwtSh.AuthModes = []string{"password"}
}

func TestLogin(t *testing.T) {
	message := types.DirectMessage{}
	json.Unmarshal(loginMessage, &message)

	msg, ok := Handle(message)

	if !ok {
		t.Fatal(msg)
	}
}

func TestLoginNoPassword(t *testing.T) {
	message := types.DirectMessage{}
	json.Unmarshal(loginMessageNoPassword, &message)

	msg, ok := Handle(message)

	if ok && msg != session.ErrPasswordNotSupplied {
		t.Fatal(msg)
	}
}

func TestExecute(t *testing.T) {
	rsp, _, err := execute("echo hello, world!")

	if err != nil {
		t.Fatal(err)
	}

	if rsp != "hello, world!\n" {
		t.Fatal()
	}
}

func TestExecuteNonExistentCommand(t *testing.T) {
	_, _, err := execute("e_c_h_o")

	if err == nil {
		t.Fatal(err)
	}
}

func TestExecuteExitStatus(t *testing.T) {
	_, state, _ := execute("awk '}'")

	if state != 1 {
		t.Fatal()
	}
}

func createMessage(message string) []byte {
	return []byte(fmt.Sprintf(messageTemplate, message))
}
