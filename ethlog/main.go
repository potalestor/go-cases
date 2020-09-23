package main

import (
	"bytes"
	"fmt"
	"log/syslog"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

var endpoints = []string{
	"/nsi/organizations/",
	"/nsi/currencies/",
	"/nsi/currencies/:code",
	"/nsi/okopf/",
	"/nsi/oktmo/",
	"/nsi/kladr/",
}

// NewWithSpan - создать новый экземпляр логгера со Span
// (аналогично tlog)
func NewWithSpan() log.Logger {
	span := rand.Int63()
	return log.New("span", span)
}

// CbgFormat - форматированный вывод в консоль и файл
func CbgFormat(r *log.Record) []byte {
	buf := bytes.NewBuffer(nil)
	location := fmt.Sprintf("%+v", r.Call)
	fmt.Fprintf(buf, "[%s] %v %v %v %s ", r.Lvl, r.Time.Format("02.01.2006 15:04:05.000"), location, r.Ctx[1], r.Msg)
	if len(r.Ctx) > 2 {
		buf.WriteByte('[')
		for i := 2; i < len(r.Ctx); i += 2 {
			fmt.Fprintf(buf, "%v=%v ", r.Ctx[i], r.Ctx[i+1])
		}
		buf.WriteByte(']')
	}
	buf.WriteByte('\n')
	return buf.Bytes()
}

// CbgColorFormat - форматированный вывод в консоль в цвете
func CbgColorFormat(r *log.Record) []byte {
	var color = 0
	switch r.Lvl {
	case log.LvlCrit:
		color = 35
	case log.LvlError:
		color = 31
	case log.LvlWarn:
		color = 33
	case log.LvlInfo:
		color = 32
	case log.LvlDebug:
		color = 36
	case log.LvlTrace:
		color = 34
	}
	buf := bytes.NewBuffer(nil)
	location := fmt.Sprintf("%+v", r.Call)
	fmt.Fprintf(buf, "\x1b[%dm%s\x1b[0m %v %v %v %s ", color, r.Lvl, r.Time.Format("02.01.2006 15:04:05.000"), location, r.Ctx[1], r.Msg)
	if len(r.Ctx) > 2 {
		buf.WriteByte('[')
		for i := 2; i < len(r.Ctx); i += 2 {
			fmt.Fprintf(buf, "%v=%v ", r.Ctx[i], r.Ctx[i+1])
		}
		buf.WriteByte(']')
	}
	buf.WriteByte('\n')
	return buf.Bytes()
}

// Rfc5424Format - форматированный вывод в syslog
func Rfc5424Format(r *log.Record) []byte {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte('[')
	for i := 0; i < len(r.Ctx); i += 2 {
		fmt.Fprintf(buf, "%v=%v ", r.Ctx[i], r.Ctx[i+1])
	}
	buf.WriteString("] ")
	fmt.Fprintf(buf, "%s\n", r.Msg)
	return buf.Bytes()
}

func main() {

	// вывод в syslog в формате Rfc5424
	syslogHandler, _ := log.SyslogHandler(syslog.LOG_USER, "cbg", log.FormatFunc(Rfc5424Format))
	// вывод в консоль в цвете
	consoleHandler := log.StreamHandler(os.Stdout, log.FormatFunc(CbgColorFormat))

	if err := os.MkdirAll("/var/log/tf", os.ModePerm); err != nil {
		panic(err)
	}

	// ротация файла
	fileWithRotateHandler := log.StreamHandler(&lumberjack.Logger{
		Filename:   "/var/log/tf/cbg_rotate.log",
		MaxSize:    1, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	},
		log.FormatFunc(CbgFormat),
	)

	// вывод в файл в формате Cbg
	fileHandler, err := log.FileHandler("/var/log/tf/cbg.log", log.FormatFunc(CbgFormat))
	if err != nil {
		panic(err)
	}

	log.Root().SetHandler(
		log.MultiHandler(
			consoleHandler,
			syslogHandler,
			fileHandler,
			fileWithRotateHandler,
		))

	SimulateCallingEnpoints(endpoints)
}

// SimulateCallingEnpoints - конкурентный пример работы логгера
func SimulateCallingEnpoints(e []string) {
	var wg sync.WaitGroup
	wg.Add(len(e))
	defer wg.Wait()
	for i := 0; i < len(e); i++ {
		go func(i int) {
			defer wg.Done()
			elog := NewWithSpan()
			elog.Info("this is information message", "endpoints", e[i])
			time.Sleep(time.Microsecond)
			elog.Error("and error example", "endpoints", e[i])
			time.Sleep(time.Microsecond)
			elog.Debug("it's debug step", "endpoints", e[i], "iter", i)
			time.Sleep(time.Microsecond)
		}(i)
	}
}
