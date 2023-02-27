package log

import (
	"fmt"
	"log"
	"os"
)

type logger struct {
	ClassicLogger *log.Logger
	VerboseLogger *log.Logger
	level         Level
}

var l = &logger{
	ClassicLogger: log.New(os.Stderr, "", 0),
	VerboseLogger: log.New(os.Stderr, "[verbose] ", 0),
	level:         LevelClassic,
}

func SetLogLevel(level Level) {
	l.level = level
	if l.level >= LevelVerbose {
		l.ClassicLogger.SetFlags(log.Lshortfile | log.Ltime)
		l.VerboseLogger.SetFlags(log.Lshortfile | log.Ltime)
	}
}

func Print(a ...any) {
	if l.shouldPrint(LevelClassic) {
		_ = l.ClassicLogger.Output(2, fmt.Sprintln(a...))
	}
}

func Printf(format string, a ...any) {
	if l.shouldPrint(LevelClassic) {
		_ = l.ClassicLogger.Output(2, fmt.Sprintf(format, a...))
	}
}

func Verbose(a ...any) {
	if l.shouldPrint(LevelVerbose) {
		_ = l.VerboseLogger.Output(2, fmt.Sprintln(a...))
	}
}

func Verbosef(format string, a ...any) {
	if l.shouldPrint(LevelVerbose) {
		_ = l.VerboseLogger.Output(2, fmt.Sprintf(format, a...))
	}
}

func (l *logger) shouldPrint(level Level) bool {
	return l.level >= level
}
