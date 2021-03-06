package cmd

import (
	"os"
	"runtime/pprof"

	"path/filepath"

	"github.com/Unknwon/com"
	"github.com/fatih/color"
	"github.com/rai-project/cmd"
	"github.com/rai-project/config"
	"github.com/rai-project/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	isColor    bool
	isVerbose  bool
	isDebug    bool
	inShutdown int32
	appSecret  string
	configFile string
)

type prof struct{}

func (c prof) Close() error {
	pprof.StopCPUProfile()
	return nil
}

func serverOptions() []server.Option {
	return []server.Option{
		server.Stdout(os.Stdout),
		server.Stderr(os.Stderr),
	}
}

var RootCmd = &cobra.Command{
	Use:          "rai-docker-volume",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cudaCmd.Run(cmd, args)
		return nil
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initColor)

	RootCmd.AddCommand(cmd.VersionCmd)
	RootCmd.AddCommand(cmd.LicenseCmd)
	RootCmd.AddCommand(cmd.EnvCmd)
	RootCmd.AddCommand(cmd.GendocCmd)
	RootCmd.AddCommand(cmd.CompletionCmd)
	RootCmd.AddCommand(cmd.BuildTimeCmd)

	RootCmd.PersistentFlags().StringVar(&configFile, "config", "", "The absolute path to the server configuration. If not set, then the configuration file is searched.")
	RootCmd.PersistentFlags().StringVarP(&appSecret, "secret", "s", "", "The application secret.")
	RootCmd.PersistentFlags().BoolVarP(&isColor, "color", "c", true, "Toggle color output.")
	RootCmd.PersistentFlags().BoolVarP(&isVerbose, "verbose", "v", false, "Toggle verbose mode.")
	RootCmd.PersistentFlags().BoolVarP(&isDebug, "debug", "d", false, "Toggle debug mode.")

	// mark secret flag hidden
	RootCmd.PersistentFlags().MarkHidden("secret")

	viper.BindPFlag("app.secret", RootCmd.PersistentFlags().Lookup("secret"))
	viper.BindPFlag("app.debug", RootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("app.verbose", RootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("app.color", RootCmd.PersistentFlags().Lookup("color"))
}

func initConfig() {
	opts := []config.Option{
		config.AppName("rai-docker-volume"),
		config.ColorMode(isColor),
	}
	if configFile != "" && com.IsFile(configFile) {
		if c, err := filepath.Abs(configFile); err == nil {
			configFile = c
		}
		opts = append(opts, config.ConfigFileAbsolutePath(configFile))
	} else {
		opts = append(opts, config.ConfigFileBaseName(".rai_config"))
	}
	if appSecret != "" {
		opts = append(opts, config.AppSecret(appSecret))
	}
	config.Init(opts...)
}

func initColor() {
	color.NoColor = !isColor
}
