package wire

import (
	"encoding/json"
	"time"

	"opphav.io/sdk-go/internal/wirehook"
	"opphav.io/sdk-go/opphav"
)

type jsonevent struct {
	Source    opphav.Subject `json:"source"`
	Subject   opphav.Subject `json:"subject"`
	Sign      string         `json:"sign"`
	Dimension string         `json:"dimension,omitempty"`
	TraceID   opphav.TraceID `json:"trace_id"`
	WallTime  time.Time      `json:"wall_time"`
	EmitTime  time.Time      `json:"emit_time"`
}

var _ Codec = (*JSONCodec)(nil)

type newEventFunc func(source, subject opphav.Subject, obs opphav.Observation, traceID opphav.TraceID, wallTime, emitTime int64) opphav.Event

type JSONCodec struct {
	newEvent newEventFunc
}

func NewJSONCodec() *JSONCodec {
	var newEvent newEventFunc
	if cast, ok := wirehook.NewEvent.(func(opphav.Subject, opphav.Subject, opphav.Observation, opphav.TraceID, int64, int64) opphav.Event); ok {
		newEvent = cast
	}

	return &JSONCodec{
		newEvent: newEvent,
	}
}

func (c *JSONCodec) Serialize(e opphav.Event) ([]byte, error) {
	if !e.IsValid() {
		return nil, ErrInvalidEvent
	}

	var wallTime, emitTime time.Time

	if e.WallTime() != 0 {
		wallTime = time.Unix(0, e.WallTime()).UTC()
	}

	if e.EmitTime() != 0 {
		emitTime = time.Unix(0, e.EmitTime()).UTC()
	}

	return json.Marshal(jsonevent{
		Source:    e.Source(),
		Subject:   e.Subject(),
		Sign:      e.Sign(),
		Dimension: e.Dimension(),
		TraceID:   e.TraceID(),
		WallTime:  wallTime,
		EmitTime:  emitTime,
	})
}

func (c *JSONCodec) Deserialize(data []byte) (opphav.Event, error) {
	var v jsonevent
	if err := json.Unmarshal(data, &v); err != nil {
		return opphav.Event{}, err
	}

	obs := opphav.NewObservation(v.Sign, v.Dimension)
	if obs.IsZero() {
		return opphav.Event{}, ErrInvalidEvent
	}

	var wallTime, emitTime int64

	if !v.WallTime.IsZero() {
		wallTime = v.WallTime.UnixNano()
	}

	if !v.EmitTime.IsZero() {
		emitTime = v.EmitTime.UnixNano()
	}

	return c.newEvent(
		v.Source,
		v.Subject,
		obs,
		v.TraceID,
		wallTime,
		emitTime,
	), nil
}
