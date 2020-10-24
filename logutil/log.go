package logutil

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewLogger(key, val string) zerolog.Logger {
	return log.With().Str(key, val).Caller().Logger()
}

func SetGlobalLevel(l string) {
	level, err := zerolog.ParseLevel(l)
	if err != nil {
		panic(err)
	}
	zerolog.SetGlobalLevel(level)
}
