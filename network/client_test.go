package network_test

import (
	"strings"
	"testing"

	"github.com/Amaimersion/terminal-chat/network"
	"github.com/Amaimersion/terminal-chat/protocol"
)

func TestSendWithEmptyRemote(t *testing.T) {
	req := network.Request{
		Text:            "test",
		HandlerLocation: 1,
	}
	err := network.Send(req)

	if err != network.ErrMalformedRequest {
		t.Errorf("err = %v, want = %v", err, network.ErrMalformedRequest)
	}
}

func TestSendBigPayload(t *testing.T) {
	url := protocol.URL{
		Address: []byte{127, 0, 0, 1},
	}
	req := network.Request{
		Text:   strings.Repeat("a", protocol.MaxPayloadLength+1),
		Remote: url,
	}
	err := network.Send(req)

	if err != network.ErrMalformedRequest {
		t.Errorf("err = %v, want = %v", err, network.ErrMalformedRequest)
	}
}
