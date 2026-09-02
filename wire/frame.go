package wire

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

const (
	HeaderSize     = 4
	MaxPayloadSize = 64 << 10 // 64 KiB
)

var ErrPayloadTooLarge = errors.New("opphav/wire: payload exceeds limit")

func ReadFrame(r io.Reader) ([]byte, error) {
	var header [HeaderSize]byte

	if _, err := io.ReadFull(r, header[:]); err != nil {
		return nil, err
	}

	length := binary.LittleEndian.Uint32(header[:])
	if length > MaxPayloadSize {
		return nil, fmt.Errorf("%w: %d", ErrPayloadTooLarge, length)
	}

	payload := make([]byte, length)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, err
	}

	return payload, nil
}

func WriteFrame(w io.Writer, payload []byte) (int64, error) {
	if len(payload) > MaxPayloadSize {
		return 0, fmt.Errorf("%w: %d", ErrPayloadTooLarge, len(payload))
	}

	var header [HeaderSize]byte
	binary.LittleEndian.PutUint32(header[:], uint32(len(payload)))

	buffers := net.Buffers{header[:], payload}

	return buffers.WriteTo(w)
}
