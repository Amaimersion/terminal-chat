package network

import (
	"testing"
)

var handler Handler = func(_ Request) {}

func TestHandlerRegistering(t *testing.T) {
	Handle(0, handler)

	_, ok := handlers[0]

	if !ok {
		t.Error("request handler was not registered")
	}
}

func TestHandlerAllRegistering(t *testing.T) {
	HandleAll(handler)

	exists := allHandler != nil

	if !exists {
		t.Error("all requests handler was not registered")
	}
}
