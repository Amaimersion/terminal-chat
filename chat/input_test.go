package chat

import (
	"errors"
	"strings"
	"testing"
	"testing/iotest"
	"time"
)

func stringsAreEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestInputIsEqual(t *testing.T) {
	i1 := input{
		command: command{
			text:      "command1",
			argsCount: 1,
		},
		args: []string{"arg1"},
	}
	i2 := input{
		command: command{
			text:      "command2",
			argsCount: 1,
		},
		args: []string{"arg1"},
	}
	i3 := i1

	if i1.isEqual(i2) {
		t.Errorf("i1 is equal to i2, but non-equality was expected")
	}

	if !i1.isEqual(i3) {
		t.Errorf("i1 is not equal to i3, but equality was expected")
	}
}

func testReadInput(t *testing.T, in string, want input) {
	inputs := make(chan input)
	errs := make(chan error)
	timeout := time.After(time.Second)

	go func() {
		r := strings.NewReader(in)

		if err := readInput(r, inputs); err != nil {
			errs <- err
		}
	}()

	select {
	case in := <-inputs:
		if !in.isEqual(want) {
			t.Errorf("result input = %v, want = %v", in, want)
		}

		if !stringsAreEqual(in.args, want.args) {
			t.Errorf("result args = %v, want = %v", in.args, want.args)
		}
	case err := <-errs:
		t.Errorf("err = %v, want = nil", err)
	case <-timeout:
		t.Error("timeout")
	}
}

func TestReadInputPlainText(t *testing.T) {
	in := "plain text"
	want := input{
		command: commandSendText,
		args:    []string{in},
	}

	testReadInput(t, in, want)
}

func TestReadInputExitCommand(t *testing.T) {
	in := commandExit.text
	want := input{
		command: commandExit,
		args:    make([]string, 0),
	}

	testReadInput(t, in, want)
}

func TestReadInputError(t *testing.T) {
	errs := make(chan error)
	wantErr := errors.New("expected error")

	go func() {
		r := iotest.ErrReader(wantErr)

		if err := readInput(r, nil); err != nil {
			errs <- err
		}
	}()

	timeout := time.After(time.Second)

	select {
	case err := <-errs:
		if err != wantErr {
			t.Errorf("err = %v, want = %v", err, wantErr)
		}
	case <-timeout:
		t.Error("timeout")
	}
}

func testDeconstructInput(t *testing.T, in string, wantCommand command, wantArgs []string) {
	command, args, err := deconstructInput(in)

	if err != nil {
		t.Fatalf("err = %v, want = nil", err)
	}

	if command != wantCommand {
		t.Errorf("result command = %v, want = %v", command, wantCommand)
	}

	if !stringsAreEqual(args, wantArgs) {
		t.Errorf("result args = %v, want = %v", args, wantArgs)
	}
}

func TestDeconstructInputExitCommand(t *testing.T) {
	in := commandExit.text
	wantCommand := commandExit
	wantArgs := make([]string, 0)

	testDeconstructInput(t, in, wantCommand, wantArgs)
}

func TestDeconstructInputAddUserCommand(t *testing.T) {
	in := commandAddUser.text + " name url"
	wantCommand := commandAddUser
	wantArgs := []string{"name", "url"}

	testDeconstructInput(t, in, wantCommand, wantArgs)
}

func TestDeconstructInputUnnecessaryArguments(t *testing.T) {
	in := commandAddUser.text + " name url unnecessary arguments"
	wantCommand := commandAddUser
	wantArgs := []string{"name", "url"}

	testDeconstructInput(t, in, wantCommand, wantArgs)
}

func TestDeconstructInputNotEnoughArguments(t *testing.T) {
	in := commandAddUser.text
	_, _, err := deconstructInput(in)

	if err != errInvalidInput {
		t.Fatalf("err = %v, want = %v", err, errInvalidInput)
	}
}

func TestDeconstructInputSendTextCommand(t *testing.T) {
	in := "just a plain text to send"
	wantCommand := commandSendText
	wantArgs := []string{in}

	testDeconstructInput(t, in, wantCommand, wantArgs)
}

func TestDeconstructInputEmpty(t *testing.T) {
	_, _, err := deconstructInput("")

	if err != errInvalidInput {
		t.Fatalf("err = %v, want = %v", err, errInvalidInput)
	}
}

func testExtractInputArgs(t *testing.T, in, trimPrefix string, argsCount int, wantArgs []string, wantNil bool) {
	result := extractInputArgs(in, trimPrefix, argsCount)

	if wantNil {
		if result != nil {
			t.Fatalf("result = %v, want = nil", result)
		}
	} else {
		if result == nil {
			t.Fatalf("result = nil, want = %v", wantArgs)
		}

		if !stringsAreEqual(result, wantArgs) {
			t.Errorf("result = %v, want = %v", result, wantArgs)
		}
	}
}

func TestExtractInputArgs(t *testing.T) {
	trimPrefix := "/command"
	in := trimPrefix + " 11 2 3"
	wantArgs := []string{"11", "2"}
	argsCount := len(wantArgs)

	testExtractInputArgs(t, in, trimPrefix, argsCount, wantArgs, false)
}

func TestExtractInputArgsSmallLength(t *testing.T) {
	trimPrefix := "/command"
	in := trimPrefix + " 11"
	wantArgs := []string{"11", "2"}
	argsCount := len(wantArgs)

	testExtractInputArgs(t, in, trimPrefix, argsCount, wantArgs, true)
}

func TestExtractInputArgsBigLength(t *testing.T) {
	trimPrefix := "/command"
	in := trimPrefix + " 11"
	wantArgs := []string{"11", "2"}
	argsCount := 30

	testExtractInputArgs(t, in, trimPrefix, argsCount, wantArgs, true)
}

func TestExtractInputArgsZeroLength(t *testing.T) {
	trimPrefix := "/command"
	in := trimPrefix
	wantArgs := []string{}
	argsCount := len(wantArgs)

	testExtractInputArgs(t, in, trimPrefix, argsCount, wantArgs, false)
}

func TestExtractInputArgsMissingPrefix(t *testing.T) {
	trimPrefix := "/command"
	arg := "/another"
	in := arg + " " + trimPrefix
	wantArgs := []string{arg}
	argsCount := len(wantArgs)

	testExtractInputArgs(t, in, trimPrefix, argsCount, wantArgs, false)
}
