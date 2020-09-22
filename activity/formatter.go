package activity

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

type formatter interface {
	Format(m *message) string
}

type filelogFormat struct{}

func (f *filelogFormat) Format(m *message) string {
	b := bytes.NewBuffer(nil)
	fmt.Fprintf(b, "<%d>1 %s %s %s %s %s [%s] %s",
		m.Priority,
		m.Timestamp.UTC().Format(time.RFC3339),
		m.Hostname,
		m.AppName,
		nilify(m.ProcessID),
		nilify(m.MessageID),
		m.StructuredData,
		m.Msg,
	)
	return b.String()
}

func nilify(i int) string {
	if i == 0 {
		return "-"
	}
	return strconv.Itoa(i)
}

type syslogFormat struct{}

func (f *syslogFormat) Format(m *message) string {
	b := bytes.NewBuffer(nil)
	fmt.Fprintf(b, "[%s] %s",
		m.StructuredData,
		m.Msg,
	)
	return b.String()
}
