package chat

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

// command is a template of user input that intended only to the program.
type command struct {
	text      string
	argsCount int
}

var (
	commandSendText = command{
		text:      "", // everything else except specific commands can be treated as plain text to send
		argsCount: 1,
	}
	commandExit = command{
		text:      "/exit",
		argsCount: 0,
	}
	commandHelp = command{
		text:      "/help",
		argsCount: 0,
	}
	commandStartRoom = command{
		text:      "/room",
		argsCount: 1,
	}
	commandListRooms = command{
		text:      "/rooms",
		argsCount: 0,
	}
	commandDeleteRoom = command{
		text:      "/del_room",
		argsCount: 1,
	}
	commandAddUser = command{
		text:      "/user",
		argsCount: 2,
	}
	commandListUsers = command{
		text:      "/users",
		argsCount: 0,
	}
	commandDeleteUser = command{
		text:      "/del_user",
		argsCount: 1,
	}
)

// input is a parsed and structured user input.
type input struct {
	command command
	args    []string
}

func (a input) isEqual(b input) bool {
	equal :=
		a.command == b.command &&
			len(a.args) == len(b.args)

	return equal
}

var (
	errInvalidInput = errors.New("invalid input")
)

// readInput reads input from r until EOF or unexpected non-EOF error.
// Next it parses input data to make it structured.
// If input is valid, then it will be sended to ch.
func readInput(r io.Reader, ch chan<- input) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		t := scanner.Text()
		command, args, err := deconstructInput(t)

		if err != nil {
			// let's ignore all input errors and
			// continue waiting for next input
			continue
		}

		in := input{
			command: command,
			args:    args,
		}

		ch <- in
	}

	err := scanner.Err()

	return err
}

// deconstructInput deconstructs input into separate pieces.
// It returns command that should handle that input,
// appropriate number of unhandled input arguments.
// errInvalidInput will be returned in case if unable
// to deconstruct input properly.
func deconstructInput(i string) (c command, args []string, err error) {
	if strings.HasPrefix(i, commandHelp.text) {
		if args := extractInputArgs(
			i,
			commandHelp.text,
			commandHelp.argsCount,
		); args != nil {
			return commandHelp, args, nil
		}
	} else if strings.HasPrefix(i, commandExit.text) {
		if args := extractInputArgs(
			i,
			commandExit.text,
			commandExit.argsCount,
		); args != nil {
			return commandExit, args, nil
		}
	} else if strings.HasPrefix(i, commandListRooms.text) {
		if args := extractInputArgs(
			i,
			commandListRooms.text,
			commandListRooms.argsCount,
		); args != nil {
			return commandListRooms, args, nil
		}
	} else if strings.HasPrefix(i, commandStartRoom.text) {
		if args := extractInputArgs(
			i,
			commandStartRoom.text,
			commandStartRoom.argsCount,
		); args != nil {
			return commandStartRoom, args, nil
		}
	} else if strings.HasPrefix(i, commandDeleteRoom.text) {
		if args := extractInputArgs(
			i,
			commandDeleteRoom.text,
			commandDeleteRoom.argsCount,
		); args != nil {
			return commandDeleteRoom, args, nil
		}
	} else if strings.HasPrefix(i, commandListUsers.text) {
		if args := extractInputArgs(
			i,
			commandListUsers.text,
			commandListUsers.argsCount,
		); args != nil {
			return commandListUsers, args, nil
		}
	} else if strings.HasPrefix(i, commandAddUser.text) {
		if args := extractInputArgs(
			i,
			commandAddUser.text,
			commandAddUser.argsCount,
		); args != nil {
			return commandAddUser, args, nil
		}
	} else if strings.HasPrefix(i, commandDeleteUser.text) {
		if args := extractInputArgs(
			i,
			commandDeleteUser.text,
			commandDeleteUser.argsCount,
		); args != nil {
			return commandDeleteUser, args, nil
		}
	} else if len(i) > 0 {
		return commandSendText, []string{i}, nil
	}

	return command{}, nil, errInvalidInput
}

// extractInputArgs extracts arguments from input.
// trimPrefix is a string prefix that should be trimmed.
// argsCount is a number of expected arguments.
// If number of extracted arguments doesn't match to argsCount,
// then nil will be returned.
func extractInputArgs(i, trimPrefix string, argsCount int) []string {
	i = strings.TrimPrefix(i, trimPrefix)
	i = strings.TrimSpace(i)

	if len(i) == 0 && argsCount != 0 {
		return nil
	}

	result := strings.SplitN(i, " ", argsCount+1)

	if len(result) < argsCount {
		return nil
	}

	args := result[:argsCount]

	return args
}
