package opphav

import "errors"

var ErrRetryable = errors.New("opphav: retryable")

type Transport interface {
	Send(Event) error
	Close() error
}
