package chat

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Amaimersion/terminal-chat/network"
	"github.com/Amaimersion/terminal-chat/protocol"
)

// handleError handles non critical error that don't affects
// at program work.
//
// It returns an error message that can be logged to user.
func handleError(err error) string {
	m := fmt.Sprintf(
		"%v",
		err.Error(),
	)

	return m
}

// handleWelcome returns welcome message.
//
// Welcome message is smaller than help message.
// Welcome message intended to be displayed every time when user
// executes the program, while help message can be used to provide
// full documentation about the program.
func handleWelcome() string {
	var m string

	m += "Welcome to multiroom chat."
	m += "\n"
	m += "Type \"" + commandHelp.text + "\" for more information."

	return m
}

// handleHelp returns help message about the program.
//
// See handleWelcome() documentation for difference with welcome message.
func handleHelp() string {
	var m string

	m += "It is a multiroom chat over STTP protocol."
	m += "\n\n"
	m += "You can have many chat rooms, every chat room can have many receivers."
	m += " Room organization is independent."
	m += " It is means that room organization (room name, receivers in that room, messages in that room that are visible to you) will not match on receiver side, which will have its own organization."
	m += "\n\n"
	m += "Sender will be not able to contact receiver until receiver allows accepting of messages from this sender."
	m += " Every chat side (sender or receiver) may stop receiving of messages from another side without any notifications."
	m += "\n\n"
	m += "All interaction with the chat is performed using specific commands with possible positional arguments (denoted with <> signs)."
	m += " Everything else will be treated as UTF-8 messages and will be sended to all receivers in current room."
	m += "\n\n"
	m += "The commands are:"
	m += "\n"
	m += commandExit.text + " - exit the program"
	m += "\n"
	m += commandHelp.text + " - print this help message"
	m += "\n"
	m += commandStartRoom.text + " <name> - start a new room with specific name or switch to existing one. Maximum number of started rooms is " + fmt.Sprint(maxRooms) + "."
	m += "\n"
	m += commandListRooms.text + " - print information about all started rooms"
	m += "\n"
	m += commandDeleteRoom.text + " <name> - delete room with specific name"
	m += "\n"
	m += commandAddUser.text + " <name> <URL> - add a user with specific name in current room. This user will be allowed to send messages to you. URL is a this user response room URL, ask for it from him."
	m += "\n"
	m += commandListUsers.text + " - print information about all users in current room"
	m += "\n"
	m += commandDeleteUser.text + " <name> - delete user with specific name"

	return m
}

// handlePrompt returns prompt symbols.
func handlePrompt(rooms roomsState) string {
	m := ":"

	if r, ok := rooms.started[rooms.active]; ok {
		m = fmt.Sprintf(
			"[%v]:",
			r.name,
		)
	}

	return m
}

type roomID uint

type roomInfo struct {
	name string

	// Allowed values are from 0 to 15 due to protocol limitations.
	// Should be unique at runtime.
	location uint8
}

type roomsState struct {
	active  roomID
	nextNew roomID
	started map[roomID]roomInfo
}

const (
	maxRooms = 16
)

var (
	errTooMuchRooms = errors.New("reached maximum number of started rooms")
)

// handleStartRoom changes active room to another room with that name.
//
// If such room doesn't exists, then it will be created.
// Otherwise existing one will be picked.
//
// If rooms limit is reached, then nothing will be created and
// errTooMuchRooms will be returned.
func handleStartRoom(rooms roomsState, name string) (roomsState, error) {
	// Busy status of all available locations.
	// Array index is a location value, note that 0 is a valid value.
	// Array value is a status, false means not busy, true means busy.
	busyLocations := make([]bool, maxRooms)

	for id, info := range rooms.started {
		busyLocations[info.location] = true

		if info.name == name {
			rooms.active = id
			return rooms, nil
		}
	}

	location := -1

	for i, busy := range busyLocations {
		if !busy {
			location = i
			break
		}
	}

	if location == -1 {
		return rooms, errTooMuchRooms
	}

	info := roomInfo{
		name:     name,
		location: uint8(location),
	}
	rooms.started[rooms.nextNew] = info
	rooms.active = rooms.nextNew
	rooms.nextNew++

	return rooms, nil
}

// handleListRooms returns information about started rooms.
func handleListRooms(rooms roomsState, users usersState) string {
	m := ""

	for id, r := range rooms.started {
		if id == rooms.active {
			m += "> "
		}

		m += r.name
		m += " (users - " + fmt.Sprint(len(users.added[id])) + ")"
		m += "\n"
	}

	if len(m) == 0 {
		m = "No started rooms"
	} else {
		m = m[:len(m)-1] // remove last \n
	}

	return m
}

var (
	errNoSuchRoom = errors.New("no such room")
)

// handleDeleteRoom deletes specific room and it related data.
//
// errNoSuchRoom will be returned in case if room doesn't exists.
//
// Note that although it returns updated roomsState and usersState,
// some fields of original objects still will be modified.
// For example, roomsState.started will be changed in both original
// object and passed object, because it is how map type works.
// So, you shouldn't compare result object with original object.
// For better readability, this function explicitly returns data that
// was modified in some way (including reference fields).
func handleDeleteRoom(rooms roomsState, users usersState, name string) (roomsState, usersState, error) {
	var id roomID
	ok := false

	for rID, r := range rooms.started {
		if r.name == name {
			id = rID
			ok = true
			break
		}
	}

	if !ok {
		return rooms, users, errNoSuchRoom
	}

	delete(users.added, id)
	delete(rooms.started, id)

	// if possible, we will switch to arbitrary room
	for rID := range rooms.started {
		rooms.active = rID
		break
	}

	return rooms, users, nil
}

type userInfo struct {
	name string
	url  protocol.URL
}

type usersState struct {
	added map[roomID][]userInfo
}

type handleAddUserInput struct {
	rooms     roomsState
	users     usersState
	name, url string
	_         struct{} // to prevent unkeyed literals
}

var (
	errUserExists     = errors.New("such user already exsits")
	errRoomNotStarted = errors.New("room not started")
)

// handleAddUser adds new user in active room.
// It returns updated state.
//
// If user with such URL already exists, then errUserExists
// will be returned. If room not started, then errRoomNotStarted
// will be returned.
func handleAddUser(in handleAddUserInput) (usersState, error) {
	url := protocol.URL{}

	if err := url.FromString(in.url); err != nil {
		return in.users, err
	}

	roomID := in.rooms.active

	if _, ok := in.rooms.started[roomID]; !ok {
		return in.users, errRoomNotStarted
	}

	for _, u := range in.users.added[roomID] {
		if u.url.IsEqual(url) {
			return in.users, errUserExists
		}
	}

	info := userInfo{
		name: in.name,
		url:  url,
	}
	in.users.added[roomID] = append(in.users.added[roomID], info)

	return in.users, nil
}

// handleListUsers returns information about all users in active room.
func handleListUsers(rooms roomsState, users usersState) string {
	m := ""

	if usrs, ok := users.added[rooms.active]; ok {
		for _, u := range usrs {
			m += u.name
			m += " (URL - " + u.url.String() + ")"
			m += "\n"
		}
	}

	if len(m) == 0 {
		m = "No users in current room"
	} else {
		m = m[:len(m)-1] // remove last \n
	}

	return m
}

var (
	errNoSuchUser = errors.New("no such user")
)

// handleDeleteUser deletes all users with specific name in active room.
//
// errNoSuchUser will be returned in case if such user doesn't exists.
func handleDeleteUser(rooms roomsState, users usersState, name string) (usersState, error) {
	if _, ok := users.added[rooms.active]; !ok {
		return users, errNoSuchUser
	}

	oldAdded := users.added[rooms.active]
	newAdded := make([]userInfo, 0, cap(oldAdded)-1)

	for _, u := range oldAdded {
		if u.name != name {
			newAdded = append(newAdded, u)
		}
	}

	if len(oldAdded) == len(newAdded) {
		return users, errNoSuchUser
	}

	users.added[rooms.active] = newAdded

	return users, nil
}

// handleSendText handles sending of text to all users in active room.
//
// Channel which returns all errors that occurred during requests
// will be returned. It will be closed when all requests will be done
// (either with success or fail). Message composed on behalf of user
// will be returned.
func handleSendText(rooms roomsState, users usersState, text string) (<-chan error, message) {
	var wg sync.WaitGroup
	errs := make(chan error)
	responseRoom := rooms.started[rooms.active]

	// we will not allow zero-length text from user
	// because it is reserved for internal purposes.
	if len(text) != 0 {
		receivers := users.added[rooms.active]

		for _, user := range receivers {
			req := network.Request{
				Text:            text,
				Remote:          user.url,
				HandlerLocation: responseRoom.location,
			}

			wg.Add(1)
			go func() {
				defer wg.Done()

				if err := network.Send(req); err != nil {
					errs <- err
				}
			}()
		}
	}

	go func() {
		wg.Wait()
		close(errs)
	}()

	m := message{
		outgoing: true,
		text:     text,
		room:     responseRoom.name,
		at:       time.Now(),
	}

	return errs, m
}

type handleReceiveTextInput struct {
	rooms    roomsState
	users    usersState
	from     protocol.URL
	location uint8
	text     string
}

var (
	errNoDestinationRoom       = errors.New("destination room doesn't exists")
	errNoUserInDestinationRoom = errors.New("no such user in destination room")
	errReceivedTextIsInternal  = errors.New("received text is for internal purposes only")
)

// handleReceiveText handles receiving of text from remote user.
//
// If destination room doesn't exists, then errNoDestinationRoom
// will be returned. If sender doesn't exists in destination room,
// then messages from him are not allowed due to security reasons
// and errNoUserInDestinationRoom will be returned. If received text
// is reserved to be used only for internal purposes, then
// errReceivedTextIsInternal will be returned, but all internal actions
// will be maded.
//
// If actual sender information doesn't match to local information
// about sender, then local info will be updated with actual data,
// that's why usersState will be returned.
//
// Composed message from sender will be returned.
func handleReceiveText(in handleReceiveTextInput) (usersState, message, error) {
	var destRoom roomInfo
	var destRoomID roomID
	ok := false

	for id, r := range in.rooms.started {
		if r.location == in.location {
			destRoom = r
			destRoomID = id
			ok = true
			break
		}
	}

	if !ok {
		return in.users, message{}, errNoDestinationRoom
	}

	var sender userInfo
	var senderIndx int
	ok = false

	for i, u := range in.users.added[destRoomID] {
		if u.url.IsEqualIP(in.from) {
			// TODO:
			// upcoming behavior (actualization of location),
			// doesn't works correctly with multiple clients
			// (different TCP ports, locations, etc.).
			// So, at the moment we will not allow different location.
			if in.from.Location != u.url.Location {
				continue
			}

			sender = u
			senderIndx = i
			ok = true

			break
		}
	}

	if !ok {
		return in.users, message{}, errNoUserInDestinationRoom
	}

	// sender changed its location, old one not valid anymore.
	// TODO:
	// doesn't works correctly.
	if sender.url.Location != in.from.Location {
		sender.url.Location = in.from.Location
		in.users.added[destRoomID][senderIndx] = sender
	}

	if len(in.text) == 0 {
		return in.users, message{}, errReceivedTextIsInternal
	}

	m := message{
		text:     in.text,
		from:     sender.name,
		room:     destRoom.name,
		at:       time.Now(),
		outgoing: false,
	}

	return in.users, m, nil
}

var (
	errRoomUnavailable = errors.New("unable to retrieve room URL")
)

// handleGetRoomURL returns message that contains all URLs by which
// active room can be reached by other peers.
//
// port is a TCP port of program.
//
// In case if unable to get at least one URL, errRoomUnavailable will
// be returned.
func handleGetRoomURL(rooms roomsState, port uint16) (string, error) {
	info, ok := rooms.started[rooms.active]

	if !ok {
		return "", errRoomUnavailable
	}

	// do not care about errors here, they will be handled later
	outbound, _ := network.LookupOutbound()
	systemNet, _ := network.LookupSystemNetwork()

	result := ""
	outboundWritten := false

	for _, ip := range systemNet {
		ipType := "local"

		if outbound != nil && outbound.Equal(ip) {
			ipType = "outbound"
			outboundWritten = true
		}

		url := protocol.URL{
			Address:  ip,
			Port:     port,
			Location: info.location,
		}
		s := fmt.Sprintf("%v (%v)", url.String(), ipType)
		result += s + "\n"
	}

	if outbound != nil && !outboundWritten {
		ipType := "outbound"
		url := protocol.URL{
			Address:  outbound,
			Port:     port,
			Location: info.location,
		}
		s := fmt.Sprintf("%v (%v)", url.String(), ipType)
		result += s + "\n"
	}

	if len(result) == 0 {
		return "", errRoomUnavailable
	}

	result = result[:len(result)-1] // trim last \n
	result =
		"Users can reach this room using following URLs.\n" +
			"Pick \"outbound\" one if available.\n" +
			"URLs:\n" +
			result

	return result, nil
}
