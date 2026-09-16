package main

import (
	"log"
	"net/http"

	"opphav.io/sdk-go/opphav"
	"opphav.io/sdk-go/symbols/operations"
	"opphav.io/sdk-go/symbols/services"
	"opphav.io/sdk-go/transport"
	"opphav.io/sdk-go/wire"
)

const (
	url     = "https://httpbin.org/get"
	traceID = opphav.TraceID("867d85aa-63b3-4edf-9bc2-d1f384ae1a1e")
)

func main() {
	dispatcher := opphav.NewDispatcher(
		transport.NewStdoutTransport(wire.NewJSONCodec()),
	)
	defer dispatcher.Close()

	var (
		source  = opphav.NewName("example.client").Subject()
		subject = opphav.NewName("httpbin.org").Subject()
	)

	percept := opphav.NewPercept(dispatcher, source, traceID)

	percept.Emit(subject, operations.BEGIN)
	defer percept.Emit(subject, operations.END)

	percept.Emit(subject, services.Call(services.CALLER))

	resp, err := http.Get(url)
	if err != nil {
		percept.Emit(subject, services.Fail(services.CALLER))
		log.Printf("HTTP request failed: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		percept.Emit(subject, services.Success(services.CALLER))
	} else {
		percept.Emit(subject, services.Fail(services.CALLER))
	}
}
