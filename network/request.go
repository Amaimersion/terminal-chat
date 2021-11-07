package network

import (
	"errors"

	"github.com/Amaimersion/terminal-chat/protocol"
)

var (
	// ErrMalformedRequest is returned when request
	// have either invalid value or missing properties.
	//
	// Note that this happens before sending.
	// So it indicates invalid logic at sender side,
	// not that receiver is not able to handle this request
	ErrMalformedRequest = errors.New("malformed request")
)

// Request is an incoming data from client
type Request struct {
	// UTF-8 text
	Text string

	// At which location handler is expected to exists
	// to handle request.
	//
	// For arrived requests it equal to the location of
	// potential handler who should handle this request
	// according to sender opinion.
	//
	// For outgoing requests it equal to the location of
	// handler who should handle potential response
	// (as separete request) according to sender opinion.
	HandlerLocation uint8

	// URL of remote peer.
	//
	// For arrived requests it equal to the sender URL.
	//
	// For outgoing requests it equal to the receiver URL.
	Remote protocol.URL
}
