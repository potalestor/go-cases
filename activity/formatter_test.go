package activity

import (
	"log/syslog"
	"testing"
	"time"
)

func Test_filelogFormat_String(t *testing.T) {
	tests := []struct {
		name string
		m    *message
		want string
	}{
		{
			"with STRUCTURED-DATA",
			newMessageBuilder().
				buildPriority(syslog.LOG_USER|syslog.LOG_INFO).
				buildTimestamp(time.Date(2020, time.September, 21, 10, 50, 35, 0, time.Local)).
				buildHostname().
				buildAppName("cbg").
				buildProcessID(10).
				buildStructuredData(NewParam("key", "value"), "string", 10).
				buildMsg(Auth).
				build(),
			"<14>1 2020-09-21T07:50:35Z SF314-57 cbg 10 - [key=value string 10] Authentication OK",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &filelogFormat{}
			if got := f.Format(tt.m); got != tt.want {
				t.Errorf("filelogFormat.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_syslogFormat_String(t *testing.T) {
	tests := []struct {
		name string
		m    *message
		want string
	}{
		{
			"with STRUCTURED-DATA",
			newMessageBuilder().
				buildStructuredData(NewParam("key", "value"), "string", 10).
				buildMsg(Auth | AuthUserNotRegistered).
				build(),
			"[key=value string 10] Authentication failed: user is not registered",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &syslogFormat{}
			if got := f.Format(tt.m); got != tt.want {
				t.Errorf("syslogFormat.String() = %v, want %v", got, tt.want)
			}
		})
	}

}
