package ps

import (
	"encoding/base64"
	"fmt"
	"os/exec"
)

// Run ...
func Run(script []byte) ([]byte, error) {
	cmd := exec.Command(
		"powershell.exe",
		"-EncodedCommand",
		base64.StdEncoding.EncodeToString(
			utf16LeEncode(script),
		),
	)

	outCapturer := newCapturer(cmd)
	errCapturer := newErrCapturer(cmd)

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("cmd.Run() failed with %v\n%v", err, errCapturer.Error())
	}
	return outCapturer.Bytes(), errCapturer.Error()
}
