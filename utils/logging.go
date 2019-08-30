package utils

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

type GosibleLogger struct {
	// Wire is a handle for most verbose logging
	Wire *log.Logger
	// Trace is a handle for trace logging
	Trace *log.Logger
	// Info is a handle for info level logging
	Info *log.Logger
	// Warning is a handle for warning level logging
	Warning *log.Logger
	// Error is a handle for error logging
	Error *log.Logger
}

func NewGosibleLogger(wireHandle io.Writer,
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) *GosibleLogger {

	l := GosibleLogger{}
	l.Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	l.Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	l.Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	l.Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	l.Wire = log.New(wireHandle,
		"WIRE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	return &l
}

func NewGosibleDefaultLogger() *GosibleLogger {
	return NewGosibleLogger(ioutil.Discard, ioutil.Discard, os.Stdout, os.Stdout, os.Stdout)
}
