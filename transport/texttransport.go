package transport

import (
	"fmt"
	"io"
	"os"

	"opphav.io/sdk-go/opphav"
	"opphav.io/sdk-go/wire"
)

var _ Transport = (*TextTransport)(nil)

type TextTransport struct {
	writer     io.Writer
	serializer wire.Serializer
}

func NewTextTransport(writer io.Writer, serializer wire.Serializer) *TextTransport {
	return &TextTransport{
		writer:     writer,
		serializer: serializer,
	}
}

func (t *TextTransport) Send(e opphav.Event) error {
	data, err := t.serializer.Serialize(e)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(t.writer, "%s\n", data)

	return err
}

func (t *TextTransport) Close() error {
	return nil
}

func NewStdoutTransport(serializer wire.Serializer) *TextTransport {
	return NewTextTransport(os.Stdout, serializer)
}

func NewStderrTransport(serializer wire.Serializer) *TextTransport {
	return NewTextTransport(os.Stderr, serializer)
}
