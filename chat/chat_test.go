package chat

import (
	"errors"
	"io"
	"math"
	"strconv"
	"strings"
	"testing"
	"testing/iotest"
	"time"
)

func TestRunBadAddress(t *testing.T) {
	done := make(chan bool)
	timer := time.After(time.Second * 3)

	go func() {
		defer func() {
			done <- true
		}()

		f := Flags{
			In:   strings.NewReader(""),
			Out:  io.Discard,
			Port: "1234",

			// For bad address we can use either non-resolvable hostname or
			// IP address with TCP port. IP along with TCP port is not allowed
			// because we have separate Port flag that intended for TCP port.
			Address: "127.0.0.1:4321",
		}
		err := Run(f)

		if err == nil {
			t.Errorf("err = nil, want net error")
		}
	}()

	select {
	case <-done:
		return
	case <-timer:
		t.Errorf("timeout")
	}
}

func TestRunBadIn(t *testing.T) {
	done := make(chan bool)
	timer := time.After(time.Second * 3)

	go func() {
		defer func() {
			done <- true
		}()

		e := errors.New("simulated error")
		r := iotest.ErrReader(e)
		f := Flags{
			In:      r,
			Out:     io.Discard,
			Address: "127.0.0.1",
			Port:    "1234",
		}
		err := Run(f)

		if err != e {
			t.Errorf("err = %v, want = %v", err, e)
		}
	}()

	select {
	case <-done:
		return
	case <-timer:
		t.Errorf("timeout")
	}
}

func TestParsePort(t *testing.T) {
	p, err := parsePort("1234")

	if err != nil {
		t.Fatalf("err = %v, want = nil", err)
	}

	var want uint16 = 1234

	if p != want {
		t.Errorf("result = %v, want = %v", p, want)
	}
}

func TestParsePortTooBig(t *testing.T) {
	_, err := parsePort(strconv.Itoa(math.MaxUint16 + 1))

	if err == nil {
		t.Errorf("err = nil, want some error")
	}
}

func TestInitChat(t *testing.T) {
	state := chatState{
		rooms: roomsState{
			active:  0,
			nextNew: 0,
			started: make(map[roomID]roomInfo),
		},
		users: usersState{
			added: make(map[roomID][]userInfo),
		},
		port: 4444,
	}
	_, err := initChat(io.Discard, state)

	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}
}
