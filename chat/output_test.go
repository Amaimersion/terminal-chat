package chat

import (
	"bytes"
	"testing"
	"time"
)

func TestWrite(t *testing.T) {
	var w bytes.Buffer
	m := "message"
	err := write(&w, m)

	if err != nil {
		t.Fatalf("err = %v, want = nil", err)
	}

	if r := w.String(); r != m {
		t.Errorf("result = %v, want = %v", r, m)
	}
}

func TestMessageString(t *testing.T) {
	m := message{
		text:     "text",
		room:     "room",
		at:       time.Now(),
		from:     "from",
		outgoing: false,
	}
	s := m.string()

	if len(s) == 0 {
		t.Errorf("result string is empty")
	}
}
