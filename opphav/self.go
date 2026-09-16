package opphav

import "time"

const (
	selfDomain          = "self"
	selfTraceID TraceID = "9bf4bdb3-4863-5f21-beb6-f7b6867052af"
)

type selfSign struct {
	SignEntry
}

var selfSigns = NewSignSet(
	selfDomain,
	func(e SignEntry) selfSign {
		return selfSign{e}
	},
)

var (
	selfSubject   = NewName(selfDomain).SubjectWithInstance("opphav")
	heartbeatSign = selfSigns.Define("HEARTBEAT")
)

func heartbeat() Event {
	return Event{
		source:      selfSubject,
		subject:     selfSubject,
		observation: heartbeatSign.Observation(),
		traceID:     selfTraceID,
		wallTime:    time.Now().UnixNano(),
	}
}

func IsHeartbeat(e Event) bool {
	return e.observation.sign == heartbeatSign.name && e.traceID == selfTraceID
}
