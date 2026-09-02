package opphav

type TraceID string

func (traceID TraceID) String() string {
	return string(traceID)
}
