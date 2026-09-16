package opphav

import (
	"time"
)

type Percept struct {
	dispatcher *Dispatcher
	source     Subject
	traceID    TraceID
}

func NewPercept(dispatcher *Dispatcher, source Subject, traceID TraceID) *Percept {
	if dispatcher == nil {
		panic("opphav: dispatcher must not be nil")
	}

	return &Percept{
		dispatcher: dispatcher,
		source:     source,
		traceID:    traceID,
	}
}

// Emit records an observation for in-process instrumentation.
// Like EmitAt, it never panics on invalid or zero observables to protect
// the host service.
func (p *Percept) Emit(subject Subject, observable Observable) {
	p.emit(subject, observable, time.Now().UnixNano(), 0)
}

// EmitAt records an observation with explicit timestamps.
// Unlike NewEvent, it never panics on invalid or zero observables to protect
// the host service; invalid submissions are rejected by the dispatcher.
func (p *Percept) EmitAt(subject Subject, observable Observable, wallTime, emitTime int64) {
	p.emit(subject, observable, wallTime, emitTime)
}

func (p *Percept) Source() Subject {
	return p.source
}

func (p *Percept) TraceID() TraceID {
	return p.traceID
}

func (p *Percept) emit(subject Subject, observable Observable, wallTime, emitTime int64) {
	var obs Observation
	if observable != nil && !isNil(observable) {
		obs = observable.Observation()
	}

	p.dispatcher.submit(newEvent(p.source, subject, obs, p.traceID, wallTime, emitTime))
}
