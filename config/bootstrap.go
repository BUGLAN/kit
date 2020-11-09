package config

import "github.com/spf13/pflag"

var (
	Debug    = pflag.Bool("debug", false, "debug mode")
	LogLevel = pflag.String("log_level", "INFO", "log level")
)

func init() {
	c := NewKitConfig()
	pflag.Parse()
	if err := c.LoadFlag(pflag.CommandLine); err != nil {
		c.logger.Panic().Err(err)
	}
}
