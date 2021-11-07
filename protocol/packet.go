package protocol

import (
	"encoding/binary"
	"errors"
	"math"
)

// Packet is a structured representation
// of both protocol header and payload
type Packet struct {
	// UTF-8 text
	Payload string

	// Destination of packet at application level
	DestinationPort uint8

	// Source of packet at application level.
	// Can be used for response by receiver
	SourcePort uint8
}

var (
	// ErrTooBigPacket is returned by Marshal
	// if passed packet is too big to transmit
	ErrTooBigPacket = errors.New("packet data exceeds its maximum size or value")

	// ErrCorruptedPacket is returned by Unmarshal
	// if arrived packet data have invalid structure
	ErrCorruptedPacket = errors.New("packet data is corrupted")
)

const (
	// Maximum length of Payload
	MaxPayloadLength = math.MaxUint16

	// Maximum value for DestinationPort or SourcePort
	MaxPortValue = 15 // max of 4 bits

	fixedHeaderLength = 4
	maxPacketLength   = MaxPayloadLength + fixedHeaderLength
)

// Marshal converts packet to byte stream
func Marshal(data Packet) ([]byte, error) {
	payload := []byte(data.Payload)
	payloadLength := len(payload)
	packetLength := fixedHeaderLength + payloadLength

	if packetLength > maxPacketLength {
		return nil, ErrTooBigPacket
	}

	packet := make([]byte, packetLength)

	binary.BigEndian.PutUint16(packet[0:2], uint16(payloadLength))

	packet[2] = fixedHeaderLength
	packet[3] = 0
	packet[3] |= data.SourcePort
	packet[3] <<= 4
	packet[3] |= data.DestinationPort

	for i, value := range payload {
		packet[fixedHeaderLength+i] = value
	}

	return packet, nil
}

const (
	destinationPortMask = 0b00001111
	sourcePortMask      = 0b11110000
)

// Unmarshal converts byte stream to packet
func Unmarshal(data []byte) (Packet, error) {
	packet := Packet{}

	if len(data) < fixedHeaderLength {
		return packet, ErrCorruptedPacket
	}

	headerLength := int(data[2])

	if len(data) < headerLength {
		return packet, ErrCorruptedPacket
	}

	var payload []byte

	if len(data) == headerLength {
		payload = make([]byte, 0)
	} else {
		payload = data[headerLength:]
	}

	payloadLength := int(binary.BigEndian.Uint16(data[:2]))

	if len(payload) != payloadLength {
		return packet, ErrCorruptedPacket
	}

	destinationPort := uint8(data[3] & destinationPortMask)
	sourcePort := uint8((data[3] & sourcePortMask) >> 4)

	packet.Payload = string(payload)
	packet.DestinationPort = destinationPort
	packet.SourcePort = sourcePort

	return packet, nil
}
