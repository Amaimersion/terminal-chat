package protocol_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Amaimersion/terminal-chat/protocol"
)

func TestMarshalPacket(t *testing.T) {
	payload := "test"
	packet := protocol.Packet{
		Payload:         payload,
		DestinationPort: 0b100,
		SourcePort:      0b10,
	}
	data, err := protocol.Marshal(packet)
	expectedData := []byte{
		0,
		byte(len(payload)),
		4,
		0b00100100,
	}
	expectedData = append(expectedData, []byte(payload)...)

	if err != nil {
		t.Fatalf("err = %v, want = %v", err, nil)
	}

	if !bytes.Equal(data, expectedData) {
		t.Errorf("result data = %v, want = %v", data, expectedData)
	}
}

func TestMarshalMaxSizePacket(t *testing.T) {
	payload := strings.Repeat("a", protocol.MaxPayloadLength)
	packet := protocol.Packet{
		Payload:         payload,
		DestinationPort: 0,
		SourcePort:      0,
	}
	data, err := protocol.Marshal(packet)
	expectedData := []byte{
		255,
		255,
		4,
		0,
	}
	expectedData = append(expectedData, []byte(payload)...)

	if err != nil {
		t.Fatalf("err = %v, want = %v", err, nil)
	}

	if !bytes.Equal(data, expectedData) {
		t.Logf("result data length = %v, want = %v", len(data), len(expectedData))
		t.Logf("result data payload length = %v, want = %v", data[:2], expectedData[:2])
		t.Fail()
	}
}

func TestMarshalTooBigPacket(t *testing.T) {
	payload := strings.Repeat("a", protocol.MaxPayloadLength+1)
	packet := protocol.Packet{
		Payload: payload,
	}
	_, err := protocol.Marshal(packet)

	if err != protocol.ErrTooBigPacket {
		t.Errorf("err = %v, want = %v", err, protocol.ErrTooBigPacket)
	}
}

func TestMarshalEmptyPayload(t *testing.T) {
	packet := protocol.Packet{
		Payload:         "",
		DestinationPort: 0,
		SourcePort:      0,
	}
	data, err := protocol.Marshal(packet)
	expectedData := []byte{0, 0, 4, 0}

	if err != nil {
		t.Fatalf("err = %v, want = %v", err, nil)
	}

	if !bytes.Equal(data, expectedData) {
		t.Errorf("result data = %v, want = %v", data, expectedData)
	}
}

func TestUnmarshalStream(t *testing.T) {
	expectedPacket := protocol.Packet{
		Payload:         "test",
		DestinationPort: 8,
		SourcePort:      2,
	}
	secondByte := byte(len(expectedPacket.Payload))

	var fourthByte byte = 0
	fourthByte |= expectedPacket.SourcePort
	fourthByte <<= 4
	fourthByte |= expectedPacket.DestinationPort

	payload := []byte(expectedPacket.Payload)
	data := []byte{0, secondByte, 4, fourthByte}
	data = append(data, payload...)
	resultPacket, err := protocol.Unmarshal(data)

	if err != nil {
		t.Fatalf("err = %v, want = %v", err, nil)
	}

	if resultPacket.Payload != expectedPacket.Payload {
		t.Errorf("result payload = %v, want = %v", resultPacket.Payload, expectedPacket.Payload)
	}

	if resultPacket.DestinationPort != expectedPacket.DestinationPort {
		t.Errorf("result destination port = %v, want = %v", resultPacket.DestinationPort, expectedPacket.DestinationPort)
	}

	if resultPacket.SourcePort != expectedPacket.SourcePort {
		t.Errorf("result source port = %v, want = %v", resultPacket.SourcePort, expectedPacket.SourcePort)
	}
}

func TestUnmarshalEmptyStream(t *testing.T) {
	data := make([]byte, 0)
	_, err := protocol.Unmarshal(data)

	if err != protocol.ErrCorruptedPacket {
		t.Errorf("err = %v, want = %v", err, protocol.ErrCorruptedPacket)
	}
}

func TestUnmarshalShortStream(t *testing.T) {
	data := []byte{0, 0, 2}
	_, err := protocol.Unmarshal(data)

	if err != protocol.ErrCorruptedPacket {
		t.Errorf("err = %v, want = %v", err, protocol.ErrCorruptedPacket)
	}
}

func TestUnmarshalInvalidHeaderLength(t *testing.T) {
	data := []byte{0, 0, 3, 0}
	_, err := protocol.Unmarshal(data)

	if err != protocol.ErrCorruptedPacket {
		t.Errorf("err = %v, want = %v", err, protocol.ErrCorruptedPacket)
	}
}

func TestUnmarshalInvalidPayloadLength(t *testing.T) {
	data := []byte{0, 5, 4, 0, 1, 2, 3}
	_, err := protocol.Unmarshal(data)

	if err != protocol.ErrCorruptedPacket {
		t.Errorf("err = %v, want = %v", err, protocol.ErrCorruptedPacket)
	}
}

func TestUnmarshalEmptyPayload(t *testing.T) {
	expectedPacket := protocol.Packet{
		Payload: "",
	}
	data := []byte{0, 0, 4, 0}
	resultPacket, err := protocol.Unmarshal(data)

	if err != nil {
		t.Fatalf("err = %v, want = %v", err, nil)
	}

	if resultPacket.Payload != expectedPacket.Payload {
		t.Errorf("result payload = %v, want = %v", resultPacket.Payload, expectedPacket.Payload)
	}
}

func TestUnmarshalBigPayload(t *testing.T) {
	expectedPacket := protocol.Packet{
		Payload: strings.Repeat("a", 256),
	}
	payload := []byte(expectedPacket.Payload)
	data := []byte{1, 0, 4, 0}
	data = append(data, payload...)
	resultPacket, err := protocol.Unmarshal(data)

	if err != nil {
		t.Fatalf("err = %v, want = %v", err, nil)
	}

	if resultPacket.Payload != expectedPacket.Payload {
		t.Errorf("result payload = %v, want = %v", resultPacket.Payload, expectedPacket.Payload)
	}
}
