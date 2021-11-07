package protocol_test

import (
	"testing"

	"github.com/Amaimersion/terminal-chat/protocol"
)

func TestURLIsEmpty(t *testing.T) {
	u1 := protocol.URL{
		Address:  nil,
		Port:     44,
		Location: 1,
	}
	u2 := protocol.URL{
		Address:  []byte{0, 0, 0, 0},
		Port:     0,
		Location: 0,
	}

	if !u1.IsEmpty() {
		t.Errorf("u1 is not empty, but should be empty")
	}

	if u2.IsEmpty() {
		t.Errorf("u2 is empty, but should be not empty")
	}
}

func TestURLIsEqual(t *testing.T) {
	u1 := protocol.URL{
		Address:  []byte{127, 0, 0, 1},
		Port:     44,
		Location: 0,
	}
	u2 := protocol.URL{
		Address:  []byte{127, 0, 0, 1},
		Port:     44,
		Location: 1,
	}
	u3 := u1

	if u1.IsEqual(u2) {
		t.Errorf("u1 is equal to u2, but should be not equal")
	}

	if !u1.IsEqual(u3) {
		t.Errorf("u1 is not equal to u3, but should be equal")
	}
}

func TestURLIsEqualIP(t *testing.T) {
	u1 := protocol.URL{
		Address:  []byte{127, 0, 0, 1},
		Port:     33,
		Location: 0,
	}
	u2 := protocol.URL{
		Address:  []byte{192, 168, 1, 235},
		Port:     44,
		Location: 1,
	}
	u3 := u1

	if u1.IsEqualIP(u2) {
		t.Errorf("u1 is equal to u2, but should be not equal")
	}

	if !u1.IsEqualIP(u3) {
		t.Errorf("u1 is not equal to u3, but should be equal")
	}
}

func TestUrlString(t *testing.T) {
	url := protocol.URL{
		Address:  []byte{192, 168, 1, 235},
		Port:     1234,
		Location: 21,
	}
	s := url.String()
	want := "sttp://192.168.1.235:1234/21"

	if s != want {
		t.Errorf("result = %v, want = %v", s, url)
	}
}

func TestUrlStringTCPIP(t *testing.T) {
	url := protocol.URL{
		Address:  []byte{192, 168, 1, 235},
		Port:     1234,
		Location: 21,
	}
	s := url.StringTCPIP()
	want := "192.168.1.235:1234"

	if s != want {
		t.Errorf("result = %v, want = %v", s, url)
	}
}

func TestUrlFromString(t *testing.T) {
	url := protocol.URL{}
	err := url.FromString("sttp://0.0.0.0:3333/12")

	if err != nil {
		t.Fatalf("err = %v, want = %v", err, nil)
	}

	wantAddress := []byte{0, 0, 0, 0}
	var wantPort uint16 = 3333
	var wantLocation uint8 = 12

	if !url.Address.Equal(wantAddress) {
		t.Errorf("address = %v, want = %v", url.Address, wantAddress)
	}

	if url.Port != wantPort {
		t.Errorf("port = %v, want = %v", url.Port, wantPort)
	}

	if url.Location != wantLocation {
		t.Errorf("location = %v, want = %v", url.Location, wantLocation)
	}
}

func TestUrlFromStringWithoutScheme(t *testing.T) {
	url := protocol.URL{}
	err := url.FromString("0.0.0.0:3333/12")

	if err != nil {
		t.Fatalf("err = %v, want = %v", err, nil)
	}

	wantAddress := []byte{0, 0, 0, 0}
	var wantPort uint16 = 3333
	var wantLocation uint8 = 12

	if !url.Address.Equal(wantAddress) {
		t.Errorf("address = %v, want = %v", url.Address, wantAddress)
	}

	if url.Port != wantPort {
		t.Errorf("port = %v, want = %v", url.Port, wantPort)
	}

	if url.Location != wantLocation {
		t.Errorf("location = %v, want = %v", url.Location, wantLocation)
	}
}

func TestUrlFromStringWithDefaults(t *testing.T) {
	url := protocol.URL{}
	err := url.FromString("sttp://1.2.3.4")

	if err != nil {
		t.Fatalf("err = %v, want = %v", err, nil)
	}

	wantAddress := []byte{1, 2, 3, 4}
	var wantPort uint16 = 4444
	var wantLocation uint8 = 0

	if !url.Address.Equal(wantAddress) {
		t.Errorf("address = %v, want = %v", url.Address, wantAddress)
	}

	if url.Port != wantPort {
		t.Errorf("port = %v, want = %v", url.Port, wantPort)
	}

	if url.Location != wantLocation {
		t.Errorf("location = %v, want = %v", url.Location, wantLocation)
	}
}

func TestIsEmptyWithoutRequiredFields(t *testing.T) {
	url := protocol.URL{}
	empty := url.IsEmpty()

	if empty != true {
		t.Error("expected true, got false")
	}
}

func TestIsEmptyWithRequiredFields(t *testing.T) {
	url := protocol.URL{
		Address: []byte{0, 0, 0, 0},
	}
	empty := url.IsEmpty()

	if empty != false {
		t.Error("expected false, got true")
	}
}
