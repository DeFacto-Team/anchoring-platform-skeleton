package main

import (
	"github.com/jinzhu/configor"
	"github.com/mcuadros/go-defaults"
)

// Config structure
type Config struct {
	Factomd string `default:"https://api.factomd.net" json:"factomd" form:"factomd" query:"factomd" required:"true"`
}

// NewConfig creates config from configFile
func NewConfig(configFile string) (*Config, error) {

	config := new(Config)
	defaults.SetDefaults(config)

	if err := configor.Load(config, configFile); err != nil {
		return nil, err
	}
	return config, nil
}
