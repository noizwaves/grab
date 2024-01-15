package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func configureLogging() error {
	opts := slog.HandlerOptions{}

	logLevel := viper.GetString("log-level")
	switch strings.ToLower(logLevel) {
	case "":
		opts.Level = slog.LevelWarn
	case "debug":
		opts.Level = slog.LevelDebug
	case "info":
		opts.Level = slog.LevelInfo
	case "warn":
		opts.Level = slog.LevelWarn
	case "error":
		opts.Level = slog.LevelError
	default:
		return fmt.Errorf("invalid log level %q", logLevel)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &opts))
	slog.SetDefault(logger)

	return nil
}

func makeRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "grab",
		Short: "A fast, sudo-less package manager for your terminal programs",
		Long: `
A fast, sudo-less package manager for your terminal programs. Downloads directly form GitHub.
Supports macOS, Linux, and containers.
`,
	}

	rootCmd.PersistentFlags().String("log-level", "warn", "Logging level (i.e. debug, info, warn, error) (GRAB_LOG_LEVEL)")
	rootCmd.PersistentFlags().String("home-dir", "", "Directory to treat as home directory (GRAB_HOME_DIR)")

	viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level")) //nolint:errcheck
	viper.SetDefault("log-level", "warn")

	viper.BindPFlag("home-dir", rootCmd.PersistentFlags().Lookup("home-dir")) //nolint:errcheck
	viper.SetDefault("home-dir", "")

	viper.SetEnvPrefix("grab")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	rootCmd.AddCommand(makeInstallCommand())
	rootCmd.AddCommand(makeUpdateCommand())

	return rootCmd
}

func Execute() {
	rootCmd := makeRootCommand()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
