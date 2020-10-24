package config

import (
	"github.com/BUGLAN/kit/logutil"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultConfigName = "config"
	defaultConfigType = "yaml"
	defaultConfigPath = "."
)

type KitConfig struct {
	configName string
	configType string
	configPath []string
	logger     zerolog.Logger
	*viper.Viper
	*pflag.FlagSet
}

type Option func(c *KitConfig)

func NewKitConfig(opts ...Option) *KitConfig {
	c := &KitConfig{
		logger:  logutil.NewLogger("component", "config"),
		Viper:   viper.New(),
		FlagSet: &pflag.FlagSet{},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithConfigName(configName string) func(c *KitConfig) {
	return func(c *KitConfig) {
		c.configName = configName
	}
}

func WithConfigType(configType string) func(c *KitConfig) {
	return func(c *KitConfig) {
		c.configType = configType
	}
}

func WithConfigPath(configPath []string) func(c *KitConfig) {
	return func(c *KitConfig) {
		c.configPath = append(c.configPath, configPath...)
	}
}

// LoadConfig load kit config
func (c *KitConfig) LoadFile() {
	if c.configName == "" {
		c.configName = defaultConfigName
	}

	c.logger.Debug().Msgf("LoadFile config file %s", c.configName)
	// set name of config file
	viper.SetConfigFile(c.configName)

	if c.configType == "" {
		c.configType = defaultConfigType
	}

	c.logger.Debug().Msgf("LoadFile config type %s", c.configType)
	// set type of config file
	viper.SetConfigType(c.configType)

	if len(c.configPath) == 0 {
		c.configPath = append(c.configPath, defaultConfigPath)
	}

	// set path of config file
	for _, path := range c.configPath {
		c.logger.Debug().Msgf("LoadFile config path %s", c.configPath)
		viper.AddConfigPath(path)
	}

	if err := viper.ReadInConfig(); err != nil {
		c.logger.Panic().Err(err)
		return
	}
}

// Unmarshal struct
func (c *KitConfig) Unmarshal(v interface{}) (err error) {
	if err = c.Unmarshal(&v); err != nil {
		c.logger.Panic().Err(err)
		return
	}
	return
}

// LoadFlag load command line flag option
func (c *KitConfig) LoadFlag(flag *pflag.FlagSet) (err error) {
	if err = c.BindPFlags(flag); err != nil {
		c.logger.Error().Msg(err.Error())
		return
	}
	return
}
