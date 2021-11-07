package network

import (
	"net"

	"github.com/Amaimersion/terminal-chat/protocol"
)

// Send sends request to the specified Remote from req.
//
// ErrMalformedRequest will be returned before sending in case
// if request is malformed. Appropriate error will be returned
// in case of net error.
func Send(req Request) error {
	if req.Remote.IsEmpty() {
		return ErrMalformedRequest
	}

	packet := protocol.Packet{
		Payload:         req.Text,
		SourcePort:      req.HandlerLocation,
		DestinationPort: req.Remote.Location,
	}
	data, err := protocol.Marshal(packet)

	if err != nil {
		return ErrMalformedRequest
	}

	address := req.Remote.StringTCPIP()
	conn, err := net.Dial("tcp", address)

	if err != nil {
		return err
	}

	defer conn.Close()

	_, err = conn.Write(data)

	if err != nil {
		return err
	}

	return nil
}
