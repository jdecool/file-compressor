package logger

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	isVerbose bool
	stdLogger *log.Logger
	errLogger *log.Logger
}

func NewLogger(isVerbose bool) *Logger {
	return &Logger{
		isVerbose: isVerbose,
		stdLogger: log.New(os.Stdout, "INFO: ", log.LstdFlags),
		errLogger: log.New(os.Stderr, "ERROR: ", log.LstdFlags),
	}
}

func NewLoggerWithOutput(isVerbose bool, output io.Writer) *Logger {
	return &Logger{
		isVerbose: isVerbose,
		stdLogger: log.New(output, "INFO: ", log.LstdFlags),
		errLogger: log.New(output, "ERROR: ", log.LstdFlags),
	}
}

func (l *Logger) SetVerbose(verbose bool) {
	l.isVerbose = verbose
}

func (l *Logger) IsVerbose() bool {
	return l.isVerbose
}

func (l *Logger) Println(a ...interface{}) {
	l.stdLogger.Println(a...)
}

func (l *Logger) Printf(format string, a ...interface{}) {
	l.stdLogger.Printf(format, a...)
}

func (l *Logger) PrintlnVerbose(a ...interface{}) {
	if l.isVerbose {
		l.stdLogger.Println(a...)
	}
}

func (l *Logger) PrintfVerbose(format string, a ...interface{}) {
	if l.isVerbose {
		l.stdLogger.Printf(format, a...)
	}
}

func (l *Logger) PrintlnError(a ...interface{}) {
	l.errLogger.Println(a...)
}

func (l *Logger) PrintfError(format string, a ...interface{}) {
	l.errLogger.Printf(format, a...)
}
