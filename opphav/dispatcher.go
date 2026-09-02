package opphav

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type dispatcherOptions struct {
	queueSize         int
	heartbeatInterval time.Duration
	sendRetryInterval time.Duration
	drainTimeout      time.Duration
}

var defaultDispatcherOptions = dispatcherOptions{
	queueSize:         8192,
	heartbeatInterval: time.Second,
	sendRetryInterval: 100 * time.Millisecond,
	drainTimeout:      time.Second,
}

type DispatcherOption func(*dispatcherOptions)

func WithQueueSize(value int) DispatcherOption {
	if value <= 0 {
		panic("opphav: queue size must be greater than zero")
	}

	return func(opts *dispatcherOptions) {
		opts.queueSize = value
	}
}

func WithHeartbeatInterval(duration time.Duration) DispatcherOption {
	if duration <= 0 {
		panic("opphav: heartbeat interval must be greater than zero")
	}

	return func(opts *dispatcherOptions) {
		opts.heartbeatInterval = duration
	}
}

func WithSendRetryInterval(duration time.Duration) DispatcherOption {
	if duration <= 0 {
		panic("opphav: send retry interval must be greater than zero")
	}

	return func(opts *dispatcherOptions) {
		opts.sendRetryInterval = duration
	}
}

func WithDrainTimeout(duration time.Duration) DispatcherOption {
	if duration <= 0 {
		panic("opphav: drain timeout must be greater than zero")
	}

	return func(opts *dispatcherOptions) {
		opts.drainTimeout = duration
	}
}

type Dispatcher struct {
	transport         Transport
	queue             chan Event
	done              chan struct{}
	heartbeatInterval time.Duration
	sendRetryInterval time.Duration
	drainTimeout      time.Duration

	mu           sync.Mutex
	closed       bool
	lastEmitTime int64

	admitted      atomic.Uint64
	rejected      atomic.Uint64
	failed        atomic.Uint64
	dropped       atomic.Uint64
	discarded     atomic.Uint64
	retries       atomic.Uint64
	drainDeadline atomic.Int64

	wg   sync.WaitGroup
	once sync.Once
}

func NewDispatcher(transport Transport, options ...DispatcherOption) *Dispatcher {
	if transport == nil {
		panic("opphav: transport must not be nil")
	}

	cfg := defaultDispatcherOptions
	for _, option := range options {
		option(&cfg)
	}

	dispatcher := &Dispatcher{
		transport:         transport,
		queue:             make(chan Event, cfg.queueSize),
		done:              make(chan struct{}),
		heartbeatInterval: cfg.heartbeatInterval,
		sendRetryInterval: cfg.sendRetryInterval,
		drainTimeout:      cfg.drainTimeout,
	}

	dispatcher.wg.Go(func() {
		for e := range dispatcher.queue {
			dispatcher.deliver(e)
		}
	})

	dispatcher.submit(heartbeat())

	dispatcher.wg.Go(func() {
		ticker := time.NewTicker(dispatcher.heartbeatInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				dispatcher.submit(heartbeat())
			case <-dispatcher.done:
				return
			}
		}
	})

	return dispatcher
}

func (d *Dispatcher) deliver(e Event) {
	for {
		err := d.transport.Send(e)
		if err == nil {
			return
		}

		if !errors.Is(err, ErrRetryable) {
			d.failed.Add(1)
			return
		}

		var (
			deadline = d.drainDeadline.Load()
			sleep    = d.sendRetryInterval
		)

		if deadline > 0 {
			now := time.Now().UnixNano()

			if now > deadline {
				d.failed.Add(1)
				return
			}

			left := time.Duration(deadline - now)
			if left < sleep {
				sleep = left
			}
		}

		d.retries.Add(1)
		time.Sleep(sleep)
	}
}

func (d *Dispatcher) Queued() int {
	return len(d.queue)
}

func (d *Dispatcher) Admitted() uint64 {
	return d.admitted.Load()
}

func (d *Dispatcher) Dropped() uint64 {
	return d.dropped.Load()
}

func (d *Dispatcher) Rejected() uint64 {
	return d.rejected.Load()
}

func (d *Dispatcher) Discarded() uint64 {
	return d.discarded.Load()
}

func (d *Dispatcher) Failed() uint64 {
	return d.failed.Load()
}

func (d *Dispatcher) Retries() uint64 {
	return d.retries.Load()
}

func (d *Dispatcher) submit(e Event) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.closed {
		d.discarded.Add(1)
		return
	}

	if !e.IsValid() {
		d.rejected.Add(1)
		return
	}

	// A supplied stamp is kept. The wall clock can step backwards,
	// so the stamp never goes below the last one emitted.
	if e.emitTime == 0 {
		d.lastEmitTime = max(time.Now().UnixNano(), d.lastEmitTime)
		e.emitTime = d.lastEmitTime
	}

	select {
	case d.queue <- e:
		d.admitted.Add(1)
	default:
		d.dropped.Add(1)
	}
}

func (d *Dispatcher) Close() error {
	var err error

	d.once.Do(func() {
		d.shutdown()
		d.wg.Wait()
		err = d.transport.Close()
	})

	return err
}

func (d *Dispatcher) shutdown() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.closed = true
	d.drainDeadline.Store(time.Now().Add(d.drainTimeout).UnixNano())

	close(d.done)
	close(d.queue)
}
