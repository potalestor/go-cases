package activity

import (
	"log"
	"log/syslog"

	"github.com/nikandfor/tlog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type handler interface {
	log(activity Activity, data ...interface{})
}

func handlerFactory(cfg *Config) handler {
	tlog.Printf("initializing user activity logger")
	if cfg != nil {
		switch cfg.Type {
		case syslogType:
			h, err := newSyslogHandler(cfg, &syslogFormat{})
			if err != nil {
				tlog.Printf("user activity syslog is not initialized: %v", err)
			}
			return h
		case fileType:
			return newFilelogHandler(cfg, &filelogFormat{})
		}
	}
	return nil
}

type filelogHandler struct {
	cfg *Config
	l   *log.Logger
	f   formatter
}

func newFilelogHandler(cfg *Config, f formatter) *filelogHandler {

	l := log.New(&lumberjack.Logger{
		Filename:   cfg.Name,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	},
		"",
		log.Lmsgprefix,
	)
	tlog.Printf("user activity filelog initialized")
	return &filelogHandler{cfg: cfg, l: l, f: f}
}

func (l *filelogHandler) log(activity Activity, data ...interface{}) {
	// <PRI>VERSION TIMESTAMP HOSTNAME APP-NAME PROCID MSGID [STRUCTURED_DATA] MSG
	m := newMessageBuilder().
		buildPriority(syslog.LOG_USER | activity.severity()).
		buildTimestamp().
		buildHostname().
		buildAppName(l.cfg.Tag).
		buildProcessID().
		buildStructuredData(data).
		buildMsg(activity).
		build()

	l.l.Println(l.f.Format(m))
}

type syslogHandler struct {
	l *syslog.Writer
	f formatter
}

func newSyslogHandler(cfg *Config, f formatter) (*syslogHandler, error) {
	l, err := syslog.New(syslog.LOG_USER, cfg.Tag)
	if err != nil {
		return nil, err
	}
	tlog.Printf("user activity syslog initialized")
	return &syslogHandler{l: l, f: f}, nil
}

func (l *syslogHandler) log(activity Activity, data ...interface{}) {
	m := newMessageBuilder().
		buildStructuredData(data).
		buildMsg(activity).
		build()

	if activity.severity() == syslog.LOG_INFO {
		if err := l.l.Info(l.f.Format(m)); err != nil {

			tlog.Printf("%v", err)
		}
	} else {
		if err := l.l.Err(l.f.Format(m)); err != nil {
			tlog.Printf("%v", err)
		}
	}
}
