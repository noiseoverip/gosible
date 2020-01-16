package logging

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

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
}

func (g *GosibleLogger) Info(msg interface{}) {
	g.InfoLogger.Print(msg)
}

func (g *GosibleLogger) Verbose(msg interface{}) {
	g.VerboseLogger.Print(msg)
}

func (g *GosibleLogger) Verbosef(format string, v ...interface{}) {
	g.VerboseLogger.Printf(format, v...)
}

func (g *GosibleLogger) Warn(msg interface{}) {
	g.WarningLogger.Print(msg)
}

func NewGosibleLogger(
	wireHandle io.Writer,
	traceHandle io.Writer,
	infoHandle io.Writer,
	verboseHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) *GosibleLogger {

	l := GosibleLogger{}
	l.TraceLogger = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	l.InfoLogger = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	l.VerboseLogger = log.New(verboseHandle,
		"DEBUG: ",
		log.Ldate|log.Ltime|log.Lshortfile)

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

func NewGosibleDefaultLogger() *GosibleLogger {
	return NewGosibleLogger(ioutil.Discard, ioutil.Discard, ioutil.Discard, os.Stdout, os.Stdout, os.Stdout)
}
