package chat

import (
	"fmt"
	"io"
	"time"
)

// write writes string data into writer.
//
// It is equivalent to writeWithFormat() without any flags.
func write(w io.Writer, s string) error {
	_, err := fmt.Fprint(w, s)

	return err
}

const (
	wEndNewline        = 1 << 1
	wEndParagraph      = 1 << 2
	wEndSpace          = 1 << 3
	wAboveCurrentLine  = 1 << 4
	wDeleteCurrentLine = 1 << 5
	wStartNewline      = 1 << 6
	wRedColor          = 1 << 7
)

const (
	keyEsc = "\u001B["
)

const (
	keyCursorUp          = keyEsc + "1A"
	keyDeleteCurrentLine = keyEsc + "M"
	keyColorRed          = keyEsc + "31m"
	keyColorReset        = keyEsc + "0m"
)

// Same as write(), but implements some output formatting
// using combination of flags.
func writeWithFormat(w io.Writer, s string, flag int) error {
	if flag&wDeleteCurrentLine > 0 {
		s = keyDeleteCurrentLine + s
	}

	if flag&wAboveCurrentLine > 0 {
		s = keyCursorUp + s
	}

	if flag&wEndSpace > 0 {
		s += " "
	}

	if flag&wRedColor > 0 {
		s = keyColorRed + s + keyColorReset
	}

	if flag&wStartNewline > 0 {
		s = "\n" + s
	}

	if flag&wEndNewline > 0 {
		s += "\n"
	}

	if flag&wEndParagraph > 0 {
		s += "\n\n"
	}

	err := write(w, s)

	return err
}

type message struct {
	text string
	room string
	at   time.Time

	// If outgoing is true, then this value may be omitted.
	from string

	// If true, message is intended to be sent.
	// If false, message is considered as received.
	outgoing bool
}

func (m message) string() string {
	t := m.at.Format("15:04")
	s := ""

	if m.outgoing {
		s = fmt.Sprintf(
			"< %v %v: %v",
			t,
			m.room,
			m.text,
		)
	} else {
		s = fmt.Sprintf(
			"> %v %v %v: %v",
			t,
			m.room,
			m.from,
			m.text,
		)
	}

	return s
}
