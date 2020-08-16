package ps

import (
	"bytes"
	"os/exec"
)

type capturer struct {
	bytes.Buffer
}

func newCapturer(cmd *exec.Cmd) *capturer {
	c := &capturer{}
	cmd.Stdout = c
	return c
}

func (c *capturer) Bytes() []byte {
	return cp866Decode(c.Buffer.Bytes())
}

func (c *capturer) String() string {
	return string(c.Bytes())
}
