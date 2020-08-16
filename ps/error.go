package ps

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os/exec"
	"strings"
)

// Error ...
type Error struct {
	Error []S `xml:"S"`
}

// S ...
type S struct {
	S     string `xml:"S,attr"`
	Value string `xml:",chardata"`
}

// Bytes ...
func (e *Error) Bytes() []byte {
	var builder bytes.Buffer
	for _, s := range e.Error {
		if strings.EqualFold(s.S, "error") {
			builder.WriteString(strings.TrimRight(s.Value, "_x000D__x000A_"))
			builder.WriteString("\r\n")
		}
	}
	return builder.Bytes()
}

type errCapturer struct {
	capturer
}

func newErrCapturer(cmd *exec.Cmd) *errCapturer {
	c := &errCapturer{}
	cmd.Stderr = c
	return c
}

func (c *errCapturer) Bytes() []byte {
	data := cp866Decode(c.capturer.Buffer.Bytes())
	var e Error
	idx := bytes.Index(data, []byte(`<Objs`))
	if idx > -1 {
		if err := xml.Unmarshal(data[idx:], &e); err != nil {
			panic(err)
		}
	}
	return e.Bytes()
}

func (c *errCapturer) String() string {
	return string(c.Bytes())
}

func (c *errCapturer) Error() error {
	errString := c.String()
	if len(errString) > 0 {
		return fmt.Errorf(errString)
	}
	return nil
}
