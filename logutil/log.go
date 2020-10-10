package logutil

import (
	"flag"
	"github.com/rs/zerolog"
)

var Debug = flag.Bool("debug", false, "sets log level to debug")

func init() {
	flag.Parse()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
