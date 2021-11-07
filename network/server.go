package network

import (
	"io"
	"net"

	"github.com/Amaimersion/terminal-chat/protocol"
)

// Handler handles incoming request
type Handler = func(req Request)

var (
	handlers           = make(map[uint8]Handler)
	allHandler Handler = nil
)

// Handle binds specific handler to specific location.
// Think of "location" like "port".
// There can be only one handler on single location.
// If you will try to bind more than one handler to
// single location, then only last handler will be binded,
// and previous handler will be deleted silently
func Handle(location uint8, handler Handler) {
	handlers[location] = handler
}

// HandleAll is similar to Handle, but with difference that handler
// will handle any request,	it will not depend on specific location.
// There can be only one any handler.
// If for request location exists handler that was binded using Handle,
// then that handler will be called first, and only after its end
// will be called any handler. Calls are synchronous.
// If you will try to bind more than one any handler,
// then only last handler will be binded, and previous handler
// will be deleted silently.
func HandleAll(handler Handler) {
	allHandler = handler
}

// ListenAndServe listens for TCP connections on
// provided TCP network address and serves each
// connection using registered handlers.
//
// If unable to start listen or serve, appropriate
// error will be returned
func ListenAndServe(address string) error {
	listener, err := net.Listen("tcp", address)

	if err != nil {
		return err
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			return err
		}

		go serve(conn)
	}
}

const (
	bufferSize = 1024 * 1
)

func serve(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, bufferSize)
	data := make([]byte, 0)

	for {
		n, err := conn.Read(buffer)

		if err != nil && err != io.EOF {
			return
		}

		if n == 0 {
			break
		}

		data = append(data, buffer[:n]...)
	}

	packet, err := protocol.Unmarshal(data)

	if err != nil {
		return
	}

	remoteTCPIP := conn.RemoteAddr().String()
	remoteURL := protocol.URL{}
	remoteURL.FromString(remoteTCPIP)
	remoteURL.Location = packet.SourcePort

	request := Request{
		Text:            packet.Payload,
		HandlerLocation: packet.DestinationPort,
		Remote:          remoteURL,
	}
	handler, callHandler := handlers[packet.DestinationPort]
	callAllHandler := allHandler != nil

	if callHandler {
		handler(request)
	}

	if callAllHandler {
		allHandler(request)
	}
}
