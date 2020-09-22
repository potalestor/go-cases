package activity

import (
	"fmt"
	"log/syslog"
	"os"
	"strings"
	"time"
)

// message corresponds to the RFC 5424 message format
type message struct {
	Priority       syslog.Priority
	Timestamp      time.Time
	Hostname       string
	AppName        string
	ProcessID      int
	MessageID      int
	StructuredData string
	Msg            string
}

// messageBuilder  is used to create message
type messageBuilder struct {
	m *message
}

func newMessageBuilder() *messageBuilder {
	return &messageBuilder{m: &message{}}
}

func (b *messageBuilder) buildPriority(p syslog.Priority) *messageBuilder {
	b.m.Priority = p
	return b
}

func (b *messageBuilder) buildTimestamp(t ...time.Time) *messageBuilder {
	if len(t) == 0 {
		b.m.Timestamp = time.Now()
	} else {
		b.m.Timestamp = t[0]
	}
	return b
}

func (b *messageBuilder) buildHostname() *messageBuilder {
	name, err := os.Hostname()
	if err != nil {
		b.m.Hostname = "-"
	}
	b.m.Hostname = name
	return b
}

func (b *messageBuilder) buildAppName(s string) *messageBuilder {
	b.m.AppName = s
	return b
}

func (b *messageBuilder) buildProcessID(i ...int) *messageBuilder {
	if len(i) == 0 {
		b.m.ProcessID = os.Getpid()
	} else {
		b.m.ProcessID = i[0]
	}
	return b
}

func (b *messageBuilder) buildMessageID(i int) *messageBuilder {
	b.m.MessageID = i
	return b
}

func (b *messageBuilder) buildStructuredData(data ...interface{}) *messageBuilder {
	l := len(data)
	if l > 0 {
		format := "%v "
		var sbuilder strings.Builder
		for i := range data {
			sbuilder.WriteString(fmt.Sprintf(format, data[i]))
		}
		b.m.StructuredData = strings.TrimSpace(sbuilder.String())
	}
	return b
}

func (b *messageBuilder) buildMsg(activity Activity) *messageBuilder {
	b.m.Msg = activity.String()
	return b
}

func (b *messageBuilder) build() *message {
	return b.m
}
