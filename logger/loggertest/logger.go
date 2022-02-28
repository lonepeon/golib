package loggertest

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/lonepeon/golib/logger"
)

func NewFake(t *testing.T) (*logger.Logger, *Writer, func()) {
	var buf bytes.Buffer
	writer := Writer{
		t:       t,
		input:   &buf,
		decoder: json.NewDecoder(&buf),
	}

	log, logCloser := logger.NewLogger(&writer)

	closer := func() {
		if err := logCloser(); err != nil {
			t.Errorf("can't flush logger: %v", err)
		}
	}

	return log, &writer, closer
}

type Line map[string]interface{}

type Writer struct {
	t       *testing.T
	lines   []Line
	input   io.ReadWriter
	decoder *json.Decoder
}

func (w *Writer) Write(p []byte) (int, error) {
	n, err := w.input.Write(p)
	if err != nil {
		return n, err
	}

	var line Line
	if err := w.decoder.Decode(&line); err != nil {
		return n, err
	}

	w.lines = append(w.lines, line)

	return n, nil
}
