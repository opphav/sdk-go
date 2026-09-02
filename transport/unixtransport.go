//go:build unix

package transport

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"opphav.io/sdk-go/opphav"
	"opphav.io/sdk-go/wire"
)

var (
	ErrTransportClosed = errors.New("opphav/transport: transport closed")
	ErrNotConnected    = errors.New("opphav/transport: not connected")
	ErrWriteFailed     = errors.New("opphav/transport: write failed")
)

type unixTransportOptions struct {
	dialTimeout    time.Duration
	writeTimeout   time.Duration
	dialRetryDelay time.Duration
}

var defaultUnixTransportOptions = unixTransportOptions{
	dialTimeout:    250 * time.Millisecond,
	writeTimeout:   100 * time.Millisecond,
	dialRetryDelay: 250 * time.Millisecond,
}

type UnixTransportOption func(*unixTransportOptions)

func WithDialTimeout(duration time.Duration) UnixTransportOption {
	if duration <= 0 {
		panic("opphav/transport: dial timeout must be greater than zero")
	}

	return func(opts *unixTransportOptions) {
		opts.dialTimeout = duration
	}
}

func WithWriteTimeout(duration time.Duration) UnixTransportOption {
	if duration <= 0 {
		panic("opphav/transport: write timeout must be greater than zero")
	}

	return func(opts *unixTransportOptions) {
		opts.writeTimeout = duration
	}
}

func WithDialRetryDelay(duration time.Duration) UnixTransportOption {
	if duration <= 0 {
		panic("opphav/transport: dial retry delay must be greater than zero")
	}

	return func(opts *unixTransportOptions) {
		opts.dialRetryDelay = duration
	}
}

// UnixTransport is not safe for concurrent use.
type UnixTransport struct {
	address        string
	dialTimeout    time.Duration
	writeTimeout   time.Duration
	dialRetryDelay time.Duration
	conn           net.Conn

	serializer wire.Serializer

	closed      bool
	nextAttempt time.Time
}

var _ Transport = (*UnixTransport)(nil)

func NewUnixTransport(
	address string,
	serializer wire.Serializer,
	options ...UnixTransportOption,
) *UnixTransport {
	cfg := defaultUnixTransportOptions
	for _, option := range options {
		option(&cfg)
	}

	return &UnixTransport{
		address:        address,
		dialTimeout:    cfg.dialTimeout,
		writeTimeout:   cfg.writeTimeout,
		dialRetryDelay: cfg.dialRetryDelay,
		serializer:     serializer,
	}
}

func (t *UnixTransport) Send(e opphav.Event) error {
	// ensureConnection rejects a closed transport, so resolving the
	// conn first avoids serializing events that cannot be sent.
	conn, err := t.ensureConnection()
	if err != nil {
		return err
	}

	data, err := t.serializer.Serialize(e)
	if err != nil {
		return err
	}

	if err := conn.SetWriteDeadline(time.Now().Add(t.writeTimeout)); err != nil {
		t.dropConnection()
		return fmt.Errorf("%w: %w", ErrWriteFailed, err)
	}

	if n, err := wire.WriteFrame(conn, data); err != nil {
		if errors.Is(err, wire.ErrPayloadTooLarge) {
			return err
		}

		if n > 0 || !errors.Is(err, os.ErrDeadlineExceeded) {
			t.dropConnection()
			return fmt.Errorf("%w: %w", ErrWriteFailed, err)
		}

		return fmt.Errorf("%w: %w: %w", ErrRetryable, ErrWriteFailed, err)
	}

	return nil
}

func (t *UnixTransport) ensureConnection() (net.Conn, error) {
	if t.closed {
		return nil, ErrTransportClosed
	}

	if t.conn != nil {
		return t.conn, nil
	}

	if time.Now().Before(t.nextAttempt) {
		return nil, fmt.Errorf("%w: %w", ErrRetryable, ErrNotConnected)
	}

	conn, err := net.DialTimeout("unix", t.address, t.dialTimeout)
	if err != nil {
		t.nextAttempt = time.Now().Add(t.dialRetryDelay)
		return nil, fmt.Errorf("%w: %w: %w", ErrRetryable, ErrNotConnected, err)
	}

	t.conn = conn
	t.nextAttempt = time.Time{}

	return conn, nil
}

func (t *UnixTransport) dropConnection() {
	if t.conn != nil {
		_ = t.conn.Close()
		t.conn = nil
	}

	t.nextAttempt = time.Now().Add(t.dialRetryDelay)
}

func (t *UnixTransport) Close() error {
	if t.closed {
		return nil
	}

	t.closed = true

	conn := t.conn
	t.conn = nil

	if conn == nil {
		return nil
	}

	return conn.Close()
}
