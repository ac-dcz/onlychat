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

var (
	Info  = log.New(os.Stdin, GREEN+"[INFO]\t"+NONE, log.LstdFlags+log.Lshortfile)
	Debug = log.New(os.Stdin, YELLOW+"[DEBUG]\t"+NONE, log.LstdFlags+log.Lshortfile)
	Warn  = log.New(os.Stdin, RED+"[WARN]\t"+NONE, log.LstdFlags+log.Lshortfile)
	Error = log.New(os.Stdin, RED+"[ERROR]\t"+NONE, log.LstdFlags+log.Lshortfile)
	out   = []io.Writer{os.Stdin, os.Stdin, os.Stdin, os.Stdin}
	logs  = []*log.Logger{Info, Debug, Warn, Error}
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
