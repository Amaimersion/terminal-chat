package chat

import (
	"errors"
	"io"
	"math"
	"strconv"

	"github.com/Amaimersion/terminal-chat/network"
)

type Flags struct {
	// Where to read input from.
	In io.Reader

	// Where to write output.
	Out io.Writer

	// IP address to use for server.
	Address string

	// TCP port to for server.
	Port string
}

type chatState struct {
	rooms roomsState
	users usersState
	port  uint16
}

// Run starts an interactive chat in terminal.
//
// It is blocking function. nil will be returned
// in case of normal exit, error will be returned
// in case of unexpected critical error.
func Run(flags Flags) error {
	var err error

	state := chatState{
		rooms: roomsState{
			active:  0,
			nextNew: 0,
			started: make(map[roomID]roomInfo),
		},
		users: usersState{
			added: make(map[roomID][]userInfo),
		},
		port: 0,
	}
	state.port, err = parsePort(flags.Port)

	if err != nil {
		return err
	}

	errs := make(chan error)
	inputs := make(chan input)

	// If chAreOpen is true, don't send anything to channels,
	// because they are already closed at this moment.
	// We use such logic because we can't gracefully end
	// goroutines from outside.
	chAreOpen := true

	defer func() {
		chAreOpen = false
	}()
	defer close(errs)
	defer close(inputs)

	go func() {
		var h network.Handler = func(req network.Request) {
			state, err = handleRequest(flags.Out, state, req)

			if !chAreOpen {
				return
			}

			if err != nil {
				errs <- err
			}
		}

		network.HandleAll(h)

		addr := flags.Address + ":" + flags.Port
		err = network.ListenAndServe(addr)

		if !chAreOpen {
			return
		}

		if err != nil {
			errs <- err
		}
	}()

	go func() {
		err = readInput(flags.In, inputs)

		if !chAreOpen {
			return
		}

		if err != nil {
			errs <- err
		}
	}()

	state, err = initChat(flags.Out, state)

	if err != nil {
		return err
	}

	for {
		select {
		case err = <-errs:
			return err
		case in := <-inputs:
			if in.command == commandExit {
				return nil
			}

			state, err = handleInput(flags.Out, state, in)

			if err != nil {
				return err
			}
		}
	}
}

func parsePort(s string) (uint16, error) {
	i, err := strconv.Atoi(s)

	if err != nil {
		return 0, errors.New("unable to parse port: " + err.Error())
	}

	if i > math.MaxUint16 {
		return 0, errors.New("invalid port value")
	}

	return uint16(i), nil
}

func initChat(w io.Writer, state chatState) (chatState, error) {
	var err error

	s := handleWelcome()
	err = writeWithFormat(
		w,
		s,
		wEndParagraph,
	)

	if err != nil {
		return state, err
	}

	// at start we will create default room for fast usage
	state, err = handleInput(
		w,
		state,
		input{commandStartRoom, []string{"main"}},
	)

	return state, err
}

func handleInput(w io.Writer, st chatState, in input) (chatState, error) {
	// If you want to end the program, then return error.
	// If you want to just notify user about occured error,
	// then set values to these variables.
	var err error = nil
	var errs <-chan error = nil

	// If you want to log with default formatting, then use this variable.
	// Otherwise log by yourself with your formatting.
	var str string = ""

	switch in.command {
	case commandHelp:
		str = handleHelp()
	case commandStartRoom:
		st.rooms, err = handleStartRoom(st.rooms, in.args[0])

		if err == nil {
			str, err = handleGetRoomURL(st.rooms, st.port)
		}
	case commandListRooms:
		str = handleListRooms(st.rooms, st.users)
	case commandDeleteRoom:
		st.rooms, st.users, err = handleDeleteRoom(st.rooms, st.users, in.args[0])
	case commandAddUser:
		i := handleAddUserInput{
			rooms: st.rooms,
			users: st.users,
			name:  in.args[0],
			url:   in.args[1],
		}
		st.users, err = handleAddUser(i)
	case commandListUsers:
		str = handleListUsers(st.rooms, st.users)
	case commandDeleteUser:
		st.users, err = handleDeleteUser(st.rooms, st.users, in.args[0])
	case commandSendText:
		var m message
		errs, m = handleSendText(st.rooms, st.users, in.args[0])
		s := m.string()
		err = writeWithFormat(
			w,
			s,
			wEndNewline|wDeleteCurrentLine|wAboveCurrentLine,
		)
	}

	if err != nil {
		s := handleError(err)
		err = writeWithFormat(
			w,
			s,
			wEndParagraph|wRedColor,
		)

		if err != nil {
			return st, err
		}
	}

	if len(str) != 0 {
		err = writeWithFormat(
			w,
			str,
			wEndParagraph,
		)

		if err != nil {
			return st, err
		}
	}

	s := handlePrompt(st.rooms)
	err = writeWithFormat(
		w,
		s,
		wEndSpace,
	)

	if err != nil {
		return st, err
	}

	// We will not wait for async errors in order to not block thread.
	// They will be printed under prompt. When they are done,
	// we will print prompt once again.
	go func() {
		if errs == nil {
			return
		}

		oneWritten := false

		for err := range errs {
			if !oneWritten {
				writeWithFormat(
					w,
					"",
					wStartNewline,
				)
			}

			s := handleError(err)
			writeWithFormat(
				w,
				s,
				wEndNewline|wRedColor,
			)
			oneWritten = true
		}

		if oneWritten {
			m := handlePrompt(st.rooms)
			writeWithFormat(
				w,
				m,
				wEndSpace|wStartNewline,
			)
		}
	}()

	return st, nil
}

func handleRequest(w io.Writer, st chatState, req network.Request) (chatState, error) {
	inpt := handleReceiveTextInput{
		rooms:    st.rooms,
		users:    st.users,
		from:     req.Remote,
		location: req.HandlerLocation,
		text:     req.Text,
	}
	users, message, err := handleReceiveText(inpt)
	st.users = users

	if err == nil {
		s := message.string()
		writeWithFormat(
			w,
			s,
			wEndNewline|wDeleteCurrentLine,
		)

		s = handlePrompt(st.rooms)
		writeWithFormat(
			w,
			s,
			wEndSpace,
		)
	} else {
		// These errors can occur because of spam.
		// We will ignore them due to security reasons.
		errShouldBeIgnored :=
			err == errNoDestinationRoom ||
				err == errNoUserInDestinationRoom ||
				err == errReceivedTextIsInternal

		if errShouldBeIgnored {
			err = nil
		}
	}

	return st, err
}
