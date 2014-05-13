package shell

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"twtsh/session"
	"twtsh/types"
)

func Handle(dm types.DirectMessage) (string, bool) {
	switch strings.Split(dm.Text, " ")[0] {
	case "!login":
		return session.Authenticate(dm.SenderId, dm.SenderScreenName, dm.Text)
	case "!exit":
		return session.Exit(dm.SenderScreenName)
	default:
		msg, ok := session.Authorize(dm.SenderScreenName)

		if !ok {
			return msg, ok
		}

		result, state, err := execute(dm.Text)

		if err != nil {
			return fmt.Sprintf("<error> %s status %s", err, state), false
		}

		if state != 0 {
			return fmt.Sprintf("%d %s", result, state), true
		}

		return string(result), true
	}

	return "", true
}

func execute(message string) (string, int, error) {
	commandLine := strings.Split(message, " ")
	command := exec.Command(commandLine[0], commandLine[1:]...)
	output, err := command.CombinedOutput()
	state := 0

	if command.ProcessState != nil {
		processState := command.ProcessState.String()
		parts := strings.Split(processState, " ")
		state, _ = strconv.Atoi(parts[2])
	}

	return string(output), state, err
}
