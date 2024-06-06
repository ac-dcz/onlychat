package logger

import (
	"io"
	"log"
	"os"
)

const (
	NONE   = "\033[0m"
	GREEN  = "\033[1;33m"
	YELLOW = "\033[1;34m"
	RED    = "\033[1;31m"
)

type Level int

var (
	InfoLevel  Level = 0x1
	DebugLevel Level = 0x2
	WarnLevel  Level = 0x4
	ErrorLevel Level = 0x8
)

const LevelNum = 4

type Logger interface {
	Printf(format string, v ...any)
	Println(v ...any)
}

var (
	infoLog  = log.New(os.Stdin, GREEN+"[INFO]\t"+NONE, log.LstdFlags+log.Lshortfile)
	debugLog = log.New(os.Stdin, YELLOW+"[DEBUG]\t"+NONE, log.LstdFlags+log.Lshortfile)
	warnLog  = log.New(os.Stdin, RED+"[WARN]\t"+NONE, log.LstdFlags+log.Lshortfile)
	errorLog = log.New(os.Stdin, RED+"[ERROR]\t"+NONE, log.LstdFlags+log.Lshortfile)
	out      = []io.Writer{os.Stdin, os.Stdin, os.Stdin, os.Stdin}
	logs     = []*log.Logger{infoLog, debugLog, warnLog, errorLog}
)

func SetLevel(level Level) {
	for i := 0; i < LevelNum; i++ {
		if ((level >> i) & 1) == 1 {
			logs[i].SetOutput(out[i])
		} else {
			logs[i].SetOutput(io.Discard)
		}
	}
}

func SetOutput(level Level, w io.Writer) {
	for i := 0; i < LevelNum; i++ {
		if ((level >> i) & 1) == 1 {
			logs[i].SetOutput(w)
			out[i] = w
		}
	}
}

var (
	Info  Logger = infoLog
	Debug Logger = debugLog
	Warn  Logger = warnLog
	Error Logger = errorLog
)

func NewFileWriter(path string) io.Writer {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	return f
}
