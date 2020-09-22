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

// import (
// 	"bytes"
// 	"fmt"
// 	"log/syslog"
// 	"math"
// 	"math/rand"
// 	"os"
// 	"runtime"
// 	"sync"
// 	"time"

// 	"github.com/ethereum/go-ethereum/log"
// )

// var l log.Logger

// func parallel(n int) {
// 	var wg sync.WaitGroup
// 	wg.Add(n)
// 	defer wg.Wait()

// 	for i := 0; i < n; i++ {
// 		go func() {
// 			defer wg.Done()
// 			span := rand.Intn(10)
// 			for {
// 				b := math.MaxUint64 / (uint64(span) + 2)
// 				if b <= 100 {
// 					break
// 				}
// 				// l.Info(strconv.Itoa(10), "span", span)
// 				// time.Sleep(time.Microsecond * time.Duration(rand.Intn(100)))
// 			}
// 		}()
// 	}
// }

// // FormatSpan ...
// func FormatSpan(r *log.Record) []byte {
// 	buf := bytes.NewBuffer(nil)
// 	fmt.Fprintf(buf, "[%v] %v %v: %s\n", r.Lvl, r.Ctx[3], r.Time.Format("15:04:05"), r.Msg)
// 	return buf.Bytes()
// }

// // FormatSyslog ...
// func FormatSyslog(r *log.Record) []byte {
// 	buf := bytes.NewBuffer(nil)
// 	fmt.Fprintf(buf, "[%v=%v] %s\n", r.Ctx[2], r.Ctx[3], r.Msg)
// 	return buf.Bytes()
// }

// func Metrics() {
// 	span := rand.Intn(10)

// 	for {
// 		l.Info(fmt.Sprintf("CPU: %v GOROUT: %v", runtime.NumCPU(), runtime.NumGoroutine()), "span", span)
// 		time.Sleep(time.Second)
// 	}
// }

// func main() {
// 	runtime.GOMAXPROCS(1)

// 	l = log.New("app", "cbg")
// 	h, _ := log.SyslogHandler(syslog.LOG_USER, "cbg", log.FormatFunc(FormatSyslog))
// 	l.SetHandler(log.MultiHandler(
// 		log.StreamHandler(os.Stdout, log.FormatFunc(FormatSpan)),
// 		h,
// 	),
// 	)
// 	go Metrics()

// 	parallel(30)
// }
