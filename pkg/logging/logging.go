package logging

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

var L = NewGosibleDefaultLogger()

type GosibleLogger struct {
	// WireLogger is a handle for most verbose logging
	WireLogger *log.Logger
	// TraceLogger is a handle for trace logging
	TraceLogger *log.Logger
	// InfoLogger is a handle for info level logging
	InfoLogger *log.Logger
	// VerboseLogger is a handle for info level logging
	VerboseLogger *log.Logger
	// WarningLogger is a handle for warning level logging
	WarningLogger *log.Logger
	// ErrorLogger is a handle for error logging
	ErrorLogger *log.Logger

	Level int
}

func Info(format string, v ...interface{}) {
	L.InfoLogger.Printf(format, v...)
}

func Debug(format string, v ...interface{}) {
	L.VerboseLogger.Printf(format, v...)
}

func (g *GosibleLogger) SetLevel(level int) {
	g.Level = level
}

func (g *GosibleLogger) Info(msg interface{}) {
	g.InfoLogger.Print(msg)
}

func (g *GosibleLogger) Verbose(format string, v ...interface{}) {
	if g.Level > 0 {
		g.VerboseLogger.Printf(format, v...)
	}
}

func (g *GosibleLogger) Warn(msg interface{}) {
	g.WarningLogger.Print(msg)
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
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

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
