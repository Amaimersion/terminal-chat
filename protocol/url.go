package protocol

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// URL implements protocol URL scheme.
//
// Format: sttp://<IP address>:<TCP port>/<STTP location>
//
// It is mostly based on RFC 1738, not RFC 3986
type URL struct {
	// IP address
	Address net.IP

	// TCP port
	Port uint16

	// STTP location
	Location uint8
}

// IsEmpty indicates if struct is empty due to
// missing initialization of required fields.
func (u URL) IsEmpty() bool {
	isEmpty :=
		u.Address == nil

	return isEmpty
}

// IsEqual reports whether u and x are the same URL.
func (u URL) IsEqual(x URL) bool {
	eq :=
		u.Address.Equal(x.Address) &&
			u.Port == x.Port &&
			u.Location == x.Location

	return eq
}

// IsEqual reports whether u and x are the same URL according to their IP.
func (u URL) IsEqualIP(x URL) bool {
	eq :=
		u.Address.Equal(x.Address)

	return eq
}

// String returns full URL as string
func (u URL) String() string {
	result := fmt.Sprintf(
		"sttp://%v:%v/%v",
		u.Address.String(),
		u.Port,
		u.Location,
	)

	return result
}

// StringTCPIP returns TCP/IP URL as string
func (u URL) StringTCPIP() string {
	result := fmt.Sprintf(
		"%v:%v",
		u.Address.String(),
		u.Port,
	)

	return result
}

// ErrInvalidURL is returned in case if
// URL argument have invalid format
var ErrInvalidURL = errors.New("url have invalid format")

const (
	defaultPort     uint16 = 4444
	defaultLocation uint8  = 0
)

// FromString initializes fields from string URL
func (u *URL) FromString(s string) error {
	s = strings.ToLower(s)
	s = strings.TrimPrefix(s, "sttp://")
	parts := strings.Split(s, "/")
	location := defaultLocation

	if l := len(parts); l > 2 {
		return ErrInvalidURL
	} else if l == 2 {
		parsedLocation, err := strconv.Atoi(parts[1])

		if err != nil {
			return ErrInvalidURL
		}

		location = uint8(parsedLocation)
	}

	tcpIP := strings.Split(parts[0], ":")
	port := defaultPort

	if l := len(tcpIP); l > 2 {
		return ErrInvalidURL
	} else if l == 2 {
		parsedPort, err := strconv.Atoi(tcpIP[1])

		if err != nil {
			return ErrInvalidURL
		}

		port = uint16(parsedPort)
	}

	address := net.ParseIP(tcpIP[0])

	if address == nil {
		return ErrInvalidURL
	}

	u.Address = address
	u.Port = port
	u.Location = location

	return nil
}
