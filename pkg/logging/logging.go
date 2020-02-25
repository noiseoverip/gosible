package logging

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var Global = NewGosibleDefaultLogger()

type GosibleLogger struct {
	WireLogger    *log.Logger
	TraceLogger   *log.Logger
	InfoLogger    *log.Logger
	VerboseLogger *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
	Level         int
}

func Info(format string, v ...interface{}) {
	Global.InfoLogger.Output(2, fmt.Sprintf(format, v...))
}

func Error(format string, v ...interface{}) {
	Global.InfoLogger.Output(2, "ERROR: "+fmt.Sprintf(format, v...))
}

func Display(format string, v ...interface{}) {
	width, _, _ := terminal.GetSize(int(os.Stdout.Fd()))
	if width < 1 {
		width = 130
	}
	msg := fmt.Sprintf(format, v...)
	Info("\n%s %s", msg, strings.Repeat("*", width-len(msg)-1))
}

func Debug(format string, v ...interface{}) {
	Global.VerboseLogger.Output(2, fmt.Sprintf(format, v...))
}

func (g *GosibleLogger) SetLevel(level int) {
	g.Level = level
}

func (g *GosibleLogger) Info(msg interface{}) {
	g.InfoLogger.Print(msg)
}

func (g *GosibleLogger) Verbose(format string, v ...interface{}) {
	if g.Level > 0 {
		g.VerboseLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

func (g *GosibleLogger) Warn(format string, v ...interface{}) {
	g.WarningLogger.Output(2, fmt.Sprintf(format, v...))
}

func NewGosibleLogger(
	wireHandle io.Writer,
	traceHandle io.Writer,
	verboseHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) *GosibleLogger {

	l := GosibleLogger{}
	l.TraceLogger = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	l.InfoLogger = log.New(infoHandle,
		"", 0)

	l.SetVerbose(verboseHandle)

	l.WarningLogger = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	l.ErrorLogger = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	l.WireLogger = log.New(wireHandle,
		"WIRE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	return &l
}

func (g *GosibleLogger) SetVerbose(writer io.Writer) {
	g.VerboseLogger = log.New(writer,
		"DEBUG: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func NewGosibleDefaultLogger() *GosibleLogger {
	return NewGosibleLogger(ioutil.Discard, ioutil.Discard, ioutil.Discard, os.Stdout, os.Stdout, os.Stdout)
}

func NewGosibleSilentLogger() *GosibleLogger {
	return NewGosibleLogger(ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard, os.Stdout)
}

func NewGosibleVerboseLogger(verbosity int) *GosibleLogger {
	logger := NewGosibleDefaultLogger()

	switch {
	case verbosity > 1:
		logger.VerboseLogger.SetOutput(os.Stdout)
	case verbosity > 2:
		logger.TraceLogger.SetOutput(os.Stdout)
	}
	return logger
}
