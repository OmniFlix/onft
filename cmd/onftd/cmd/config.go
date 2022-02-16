package cmd

import (
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
)

type SimAppConfig struct {
	serverconfig.Config
}

func initAppConfig() (string, interface{}) {
	srvCfg := serverconfig.DefaultConfig()

	srvCfg.MinGasPrices = "0.001uflix"

	simAppConfig := SimAppConfig{
		Config: *srvCfg,
	}

	simAppTemplate := serverconfig.DefaultConfigTemplate

	return simAppTemplate, simAppConfig
}
