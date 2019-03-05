package logger

import (
	log "github.com/inconshreveable/log15"
)

type (
	Logger interface {
		log.Logger
	}

	Options struct {
		Command string
		Verbose bool
	}
)

func New(opt *Options) Logger {
	var l Logger
	lvl := log.LvlError
	if opt != nil {
		l = log.New(log.Ctx{
			"Command": opt.Command,
		})
		if opt.Verbose {
			lvl = log.LvlDebug
		}
	} else {
		l = log.New()
	}
	handlers := []log.Handler{}
	verboseHandler := log.LvlFilterHandler(lvl, log.StdoutHandler)
	handlers = append(handlers, verboseHandler)
	l.SetHandler(log.MultiHandler(handlers...))
	return l
}
