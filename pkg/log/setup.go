package log

import (
	"github.com/hashicorp/logutils"
	stdlog "log"
	"os"
)

func Setup(mode logutils.LogLevel) {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{Debug, Warn, Error},
		MinLevel: mode,
		Writer:   os.Stdout,
	}

	switch mode {
	case Debug:
		stdlog.SetFlags(stdlog.Ldate | stdlog.Ltime | stdlog.Lmicroseconds | stdlog.Lshortfile)
	default:
		stdlog.SetFlags(stdlog.Ldate | stdlog.Ltime)
	}

	stdlog.SetOutput(filter)
}
