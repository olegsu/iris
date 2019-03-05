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
	l := log.New(log.Ctx{
		"Command": opt.Command,
	})
	handlers := []log.Handler{}
	lvl := log.LvlError
	if opt.Verbose {
		lvl = log.LvlDebug
	}
	verboseHandler := log.LvlFilterHandler(lvl, log.StdoutHandler)
	handlers = append(handlers, verboseHandler)
	l.SetHandler(log.MultiHandler(handlers...))
	return l
}
