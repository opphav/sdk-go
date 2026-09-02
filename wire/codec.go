package wire

import (
	"errors"

	"opphav.io/sdk-go/opphav"
)

var ErrInvalidEvent = errors.New("opphav/wire: invalid event")

type Serializer interface {
	Serialize(e opphav.Event) ([]byte, error)
}

type Deserializer interface {
	Deserialize(data []byte) (opphav.Event, error)
}

type Codec interface {
	Serializer
	Deserializer
}
