package main

import (
	"github.com/jinzhu/configor"
	"github.com/mcuadros/go-defaults"
)

// Config structure
type Config struct {
	Factom struct {
		Endpoint  string `default:"https://api.factomd.net" json:"endpoint" form:"endpoint" query:"endpoint" required:"true"`
		Username  string `default:"" json:"username" form:"username" query:"username" required:"false"`
		Password  string `default:"" json:"password" form:"password" query:"password" required:"false"`
		EsAddress string `default:"" json:"esaddress" form:"esaddress" query:"esaddress" required:"false"`
	}
	Ledger struct {
		Bitcoin     Ledger
		Ethereum    Ledger
		BitcoinCash Ledger
	}
}

// Ledger is a generic sub-structure that reflects configuration for each ledger in the config
type Ledger struct {
	Endpoint   string `default:"" json:"endpoint" form:"endpoint" query:"endpoint" required:"false"`
	PrivateKey string `default:"" json:"privatekey" form:"privatekey" query:"privatekey" required:"false"`
	Frequency  int64  `default:"" json:"frequency" form:"frequency" query:"frequency" required:"false"`
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
