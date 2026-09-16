package opphav

import (
	"time"

	"opphav.io/sdk-go/internal/wirehook"
)

type Event struct {
	source      Subject
	subject     Subject
	observation Observation
	traceID     TraceID
	wallTime    int64
	emitTime    int64
}

func init() {
	wirehook.NewEvent = newEvent
}

// NewEvent constructs a validated Event with current timestamps.
// Like NewEventAt, it panics on absence or empty violations to fail fast.
func NewEvent(source, subject Subject, observable Observable, traceID TraceID) Event {
	now := time.Now().UnixNano()

	return NewEventAt(source, subject, observable, traceID, now, now)
}

// NewEventAt constructs a validated Event with explicit timestamps.
// Used for explicit event creation (e.g. sidecars or gateways), it panics
// on absence or empty violations to fail fast on misconfiguration.
func NewEventAt(source, subject Subject, observable Observable, traceID TraceID, wallTime, emitTime int64) Event {
	if isNil(observable) {
		panic("opphav: absence violation: observation is nil")
	}

	obs := observable.Observation()
	if obs.IsZero() {
		panic("opphav: absence violation: observation is empty")
	}

	return newEvent(source, subject, obs, traceID, wallTime, emitTime)
}

func (e Event) Source() Subject {
	return e.source
}

func (e Event) Subject() Subject {
	return e.subject
}

func (e Event) Sign() string {
	return e.observation.Sign()
}

func (e Event) Dimension() string {
	return e.observation.Dimension()
}

func (e Event) Observation() Observation {
	return e.observation
}

func (e Event) TraceID() TraceID {
	return e.traceID
}

func (e Event) WallTime() int64 {
	return e.wallTime
}

func (e Event) EmitTime() int64 {
	return e.emitTime
}

func (e Event) IsZero() bool {
	return e == Event{}
}

func (e Event) IsValid() bool {
	return !e.observation.IsZero()
}

func newEvent(source, subject Subject, obs Observation, traceID TraceID, wallTime, emitTime int64) Event {
	return Event{
		source:      source,
		subject:     subject,
		observation: obs,
		traceID:     traceID,
		wallTime:    wallTime,
		emitTime:    emitTime,
	}
}
