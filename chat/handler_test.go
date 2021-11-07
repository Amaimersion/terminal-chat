package chat

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/Amaimersion/terminal-chat/protocol"
)

func TestHandleError(t *testing.T) {
	e := errors.New("err")
	m := handleError(e)

	if len(m) == 0 {
		t.Errorf("empty message was returned")
	}
}

func TestHandleWelcome(t *testing.T) {
	m := handleWelcome()

	if len(m) == 0 {
		t.Errorf("empty message was returned")
	}
}

func TestHandleHelp(t *testing.T) {
	m := handleHelp()

	if len(m) == 0 {
		t.Errorf("empty message was returned")
	}
}

func TestHandlePrompt(t *testing.T) {
	var rooms roomsState
	handlePrompt(rooms)

	// We will not check length here here,
	// because empty prompt is allowed.
}

func TestHandleStartRoom(t *testing.T) {
	state := roomsState{
		active:  0,
		nextNew: 1,
		started: make(map[roomID]roomInfo),
	}
	result, err := handleStartRoom(state, "test")

	if err != nil {
		t.Fatalf("err = %v, want = nil", err)
	}

	if result.active != 1 {
		t.Errorf("active = %v, want = 1", result.active)
	}

	if l := len(result.started); l != 1 {
		t.Errorf("len(started) = %v, want = 1", l)
	}
}

func TestHandleStartRoomDuplicates(t *testing.T) {
	state := roomsState{
		active:  0,
		nextNew: 0,
		started: make(map[roomID]roomInfo),
	}
	var err error

	state, err = handleStartRoom(state, "test")

	if err != nil {
		t.Fatalf("err = %v, want = nil", err)
	}

	state, err = handleStartRoom(state, "test")

	if err != nil {
		t.Fatalf("err = %v, want = nil", err)
	}

	if l := len(state.started); l != 1 {
		t.Errorf("len(started) = %v, want = 1", l)
	}
}

func TestHandleStartRoomMaxLength(t *testing.T) {
	state := roomsState{
		active:  0,
		nextNew: 1,
		started: make(map[roomID]roomInfo),
	}

	for i := 0; i != maxRooms; i++ {
		var err error
		name := fmt.Sprintf("test%v", i)
		state, err = handleStartRoom(state, name)

		if err != nil {
			t.Fatalf("err = %v, want = nil", err)
			return
		}
	}

	if l := len(state.started); l != maxRooms {
		t.Errorf("len(started) = %v, want = %v", l, maxRooms)
	}
}

func TestHandleStartRoomTooBigLength(t *testing.T) {
	state := roomsState{
		active:  0,
		nextNew: 1,
		started: make(map[roomID]roomInfo),
	}

	for i := 0; i != maxRooms; i++ {
		var err error
		name := fmt.Sprintf("test%v", i)
		state, err = handleStartRoom(state, name)

		if err != nil {
			t.Fatalf("err = %v, want = nil", err)
			return
		}
	}

	_, err := handleStartRoom(state, "overflow")

	if err != errTooMuchRooms {
		t.Errorf("err = %v, want = %v", err, errTooMuchRooms)
	}
}

func TestHandleStartRoomDifferentLocations(t *testing.T) {
	state := roomsState{
		active:  0,
		nextNew: 1,
		started: make(map[roomID]roomInfo),
	}

	for i := 0; i != maxRooms; i++ {
		var err error
		name := fmt.Sprintf("test%v", i)
		state, err = handleStartRoom(state, name)

		if err != nil {
			t.Fatalf("err = %v, want = nil", err)
			return
		}
	}

	var locations []uint8

	for _, r := range state.started {
		for _, l := range locations {
			if r.location == l {
				t.Error("locations are not unique")
				return
			}
		}

		locations = append(locations, r.location)
	}
}

func TestHandleListRooms(t *testing.T) {
	rooms := roomsState{
		active:  0,
		nextNew: 1,
		started: map[roomID]roomInfo{
			0: {
				name:     "room1",
				location: 0,
			},
		},
	}
	users := usersState{
		added: map[roomID][]userInfo{},
	}
	m := handleListRooms(rooms, users)

	if len(m) == 0 {
		t.Errorf("result string is empty, want some content")
	}
}

func TestHandleDeleteRoom(t *testing.T) {
	rooms := roomsState{
		active:  2,
		nextNew: 3,
		started: map[roomID]roomInfo{
			0: {
				name:     "room1",
				location: 0,
			},
			1: {
				name:     "room2",
				location: 1,
			},
			2: {
				name:     "room3",
				location: 2,
			},
		},
	}
	users := usersState{
		added: map[roomID][]userInfo{
			1: {
				{
					name: "user1",
					url:  protocol.URL{},
				},
			},
		},
	}
	rooms, users, err := handleDeleteRoom(rooms, users, "room2")

	if err != nil {
		t.Fatalf("err = %v, want = nil", err)
	}

	if _, ok := rooms.started[1]; ok {
		t.Error("room have not been deleted")
	}

	if _, ok := users.added[1]; ok {
		t.Error("users have not been deleted")
	}
}

func TestHandleDeleteRoomNoSuchRoom(t *testing.T) {
	rooms := roomsState{}
	users := usersState{}
	_, _, err := handleDeleteRoom(rooms, users, "room2")

	if err != errNoSuchRoom {
		t.Fatalf("err = %v, want = %v", err, errNoSuchRoom)
	}
}

var handleAddUserInpt = handleAddUserInput{
	rooms: roomsState{
		active:  0,
		nextNew: 1,
		started: map[roomID]roomInfo{
			0: {
				name:     "room1",
				location: 0,
			},
		},
	},
	users: usersState{
		added: make(map[roomID][]userInfo),
	},
	name: "user1",
	url: protocol.URL{
		Address:  []byte{127, 0, 0, 1},
		Port:     4444,
		Location: 1,
	}.String(),
}

func TestHandleAddUser(t *testing.T) {
	users, err := handleAddUser(handleAddUserInpt)

	if err != nil {
		t.Fatalf("err = %v, want = nil", err)
	}

	if l := len(users.added); l != 1 {
		t.Errorf("len(added) = %v, want = 1", l)
	}
}

func TestHandleAddUserDuplicates(t *testing.T) {
	inpt1 := handleAddUserInpt
	inpt2 := inpt1
	inpt2.name = "user2"

	users, _ := handleAddUser(inpt1)
	users, err := handleAddUser(inpt2)

	if err != errUserExists {
		t.Fatalf("err = %v, want = %v", err, errUserExists)
	}

	if l := len(users.added); l != 1 {
		t.Errorf("len(added) = %v, want = 1", l)
	}
}

func TestHandleAddUserNoSuchRoom(t *testing.T) {
	inpt := handleAddUserInpt
	inpt.rooms.active = 1
	inpt.rooms.nextNew = 2

	_, err := handleAddUser(inpt)

	if err != errRoomNotStarted {
		t.Errorf("err = %v, want = %v", err, errRoomNotStarted)
	}
}

func TestHandleListUsers(t *testing.T) {
	rooms := roomsState{
		active:  0,
		nextNew: 1,
		started: map[roomID]roomInfo{
			0: {
				name:     "room1",
				location: 0,
			},
		},
	}
	users := usersState{
		added: map[roomID][]userInfo{
			0: {
				{
					name: "user1",
					url:  protocol.URL{},
				},
			},
		},
	}
	m := handleListUsers(rooms, users)

	if len(m) == 0 {
		t.Errorf("result string is empty, want some content")
	}
}

func TestHandleDeleteUser(t *testing.T) {
	rooms := roomsState{
		active:  0,
		nextNew: 1,
		started: map[roomID]roomInfo{
			0: {
				name:     "room1",
				location: 0,
			},
		},
	}
	users := usersState{
		added: map[roomID][]userInfo{
			0: {
				{
					name: "user1",
					url:  protocol.URL{},
				},
				{
					name: "user2",
					url:  protocol.URL{},
				},
			},
		},
	}
	_, err := handleDeleteUser(rooms, users, "user2")

	if err != nil {
		t.Fatalf("err = %v, want = nil", err)
	}

	if l := len(users.added[0]); l != 1 {
		t.Errorf("invalid count of existing elements = %v, want = 1", l)
	}

	if n := users.added[0][0].name; n != "user1" {
		t.Errorf("invalud name of remaining element = %v, want = user1", n)
	}
}

func TestHandleDeleteUserNoSuchUser(t *testing.T) {
	rooms := roomsState{
		active:  0,
		nextNew: 1,
		started: map[roomID]roomInfo{
			0: {
				name:     "room1",
				location: 0,
			},
		},
	}
	users := usersState{
		added: map[roomID][]userInfo{
			0: {
				{
					name: "user1",
					url:  protocol.URL{},
				},
			},
		},
	}
	_, err := handleDeleteUser(rooms, users, "user2")

	if err != errNoSuchUser {
		t.Fatalf("err = %v, want = %v", err, errNoSuchUser)
	}
}

func TestHandleDeleteUserMultiple(t *testing.T) {
	rooms := roomsState{
		active:  0,
		nextNew: 1,
		started: map[roomID]roomInfo{
			0: {
				name:     "room1",
				location: 0,
			},
		},
	}
	users := usersState{
		added: map[roomID][]userInfo{
			0: {
				{
					name: "user1",
					url: protocol.URL{
						Location: 0,
					},
				},
				{
					name: "user1",
					url: protocol.URL{
						Location: 1,
					},
				},
				{
					name: "user2",
					url:  protocol.URL{},
				},
			},
		},
	}
	_, err := handleDeleteUser(rooms, users, "user1")

	if err != nil {
		t.Fatalf("err = %v, want = nil", err)
	}

	if l := len(users.added[0]); l != 1 {
		t.Errorf("invalid count of existing elements = %v, want = 1", l)
	}

	if n := users.added[0][0].name; n != "user2" {
		t.Errorf("invalud name of remaining element = %v, want = user2", n)
	}
}

var (
	handleSendTextInputRooms = roomsState{
		active:  1,
		nextNew: 2,
		started: map[roomID]roomInfo{
			1: {
				name:     "room",
				location: 1,
			},
		},
	}
	handleSendTextInputUsers = usersState{
		added: map[roomID][]userInfo{
			1: {
				{
					name: "user1",
					url: protocol.URL{
						Address:  []byte{127, 0, 0, 1},
						Port:     1234,
						Location: 5,
					},
				},
			},
		},
	}
	handleSendTextInputText = "text to send"
)

func TestHandleSendTextNoReceivers(t *testing.T) {
	r := handleSendTextInputRooms
	r.active = r.nextNew
	u := handleSendTextInputUsers
	u.added[r.active] = make([]userInfo, 0)
	tx := handleSendTextInputText

	errs, _ := handleSendText(r, u, tx)

	for err := range errs {
		t.Errorf("error = %v, want = no errors at all", err)
	}
}

func TestHandleSendTextZeroLength(t *testing.T) {
	r := handleSendTextInputRooms
	u := handleSendTextInputUsers
	tx := ""

	errs, _ := handleSendText(r, u, tx)

	for err := range errs {
		t.Errorf("error = %v, want = no errors at all", err)
	}
}

var handleReceiveTextInpt = handleReceiveTextInput{
	rooms: roomsState{
		active:  1,
		nextNew: 2,
		started: map[roomID]roomInfo{
			1: {
				name:     "room1",
				location: 1,
			},
		},
	},
	users: usersState{
		added: map[roomID][]userInfo{
			1: {
				{
					name: "user1",
					url: protocol.URL{
						Address:  []byte{127, 0, 0, 1},
						Port:     3333,
						Location: 5,
					},
				},
			},
		},
	},
	from: protocol.URL{
		Address:  []byte{127, 0, 0, 1},
		Port:     55687,
		Location: 5,
	},
	location: 1,
	text:     "text to receive",
}

func TestHandleReceiveText(t *testing.T) {
	inpt := handleReceiveTextInpt
	_, msg, err := handleReceiveText(inpt)

	if err != nil {
		t.Fatalf("err = %v, want = nil", err)
	}

	subText := handleReceiveTextInpt.text

	if !strings.Contains(msg.text, subText) {
		t.Errorf("message text = %v, don't contains sub text = %v", msg.text, subText)
	}
}

func TestHandleReceiveTextNoSuchLocation(t *testing.T) {
	inpt := handleReceiveTextInpt
	inpt.location = 2

	_, _, err := handleReceiveText(inpt)

	if err != errNoDestinationRoom {
		t.Fatalf("err = %v, want = %v", err, errNoDestinationRoom)
	}
}

func TestHandleReceiveTextNoSuchUser(t *testing.T) {
	inpt := handleReceiveTextInpt
	inpt.from.Address = []byte{192, 168, 1, 235}

	_, _, err := handleReceiveText(inpt)

	if err != errNoUserInDestinationRoom {
		t.Fatalf("err = %v, want = %v", err, errNoUserInDestinationRoom)
	}
}

func TestHandleReceiveTextInternalText(t *testing.T) {
	inpt := handleReceiveTextInpt
	inpt.text = ""

	_, _, err := handleReceiveText(inpt)

	if err != errReceivedTextIsInternal {
		t.Fatalf("err = %v, want = %v", err, errReceivedTextIsInternal)
	}
}

// func TestHandleReceiveTextActualizeUser(t *testing.T) {
// 	inpt := handleReceiveTextInpt
// 	inpt.from.Location = 6
//
// 	users, _, err := handleReceiveText(inpt)
//
// 	if err != nil {
// 		t.Fatalf("err = %v, want = nil", err)
// 	}
//
// 	user := users.added[1][0]
//
// 	if user.url.Location != inpt.from.Location {
// 		t.Errorf("user data wasn't actualized")
// 	}
// }
