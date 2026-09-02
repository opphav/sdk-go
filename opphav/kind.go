package opphav

import (
	"errors"
	"fmt"
)

type Kind int

const (
	KindUnspecified Kind = iota
	KindOperation
	KindOutcome
)

func (k Kind) String() string {
	switch k {
	case KindOperation:
		return "OPERATION"
	case KindOutcome:
		return "OUTCOME"
	default:
		return "UNSPECIFIED"
	}
}

func (k Kind) MarshalText() ([]byte, error) {
	return []byte(k.String()), nil
}

func (k *Kind) UnmarshalText(text []byte) error {
	if k == nil {
		return errors.New("opphav: nil Kind pointer")
	}
	switch string(text) {
	case "OPERATION":
		*k = KindOperation
	case "OUTCOME":
		*k = KindOutcome
	case "UNSPECIFIED", "UNKNOWN", "":
		*k = KindUnspecified
	default:
		return fmt.Errorf("opphav: unknown kind %q", text)
	}
	return nil
}
