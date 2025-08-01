package cmd

import (
	"github.com/noizwaves/grab/pkg"
	"github.com/spf13/viper"
)

func newContext() (*pkg.GrabContext, error) {
	configPath := viper.GetString("config-path")
	binPath := viper.GetString("bin-path")

	return pkg.NewGrabContext(configPath, binPath) //nolint:wrapcheck
}
