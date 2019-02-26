package main

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	colorable "github.com/mattn/go-colorable"
)

var (
	log = logrus.WithFields(logrus.Fields{"package": "cmds"})

	// RootCommand is the root of the command tree.
	RootCommand *cobra.Command
)

const dlvCommandLongDesc = "publisher service for Go programs."

// NewCommandRun is
func NewCommandRun(conf *viper.Viper, run func(cmd *cobra.Command, conf *viper.Viper, args []string)) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		run(cmd, conf, args)
	}
}

// New returns an initialized command tree.
func New(docCall bool) *cobra.Command {
	conf := viper.New()

	// Main dlv root command.
	RootCommand = &cobra.Command{
		Use:   "hall-go <deploy> <name>",
		Short: "hall-go service.",
		Long:  dlvCommandLongDesc,
		Run:   NewCommandRun(conf, CommandRun),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			logrus.SetOutput(colorable.NewColorableStdout())
			logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true, TimestampFormat: "2006-01-02 15:04:05"})

			switch cmd.Flag("log").Value.String() {
			case "debug":
				logrus.SetLevel(logrus.DebugLevel)
			case "info":
				logrus.SetLevel(logrus.InfoLevel)
			case "warn":
				logrus.SetLevel(logrus.WarnLevel)
			case "error":
				logrus.SetLevel(logrus.ErrorLevel)
			case "fatal":
				logrus.SetLevel(logrus.FatalLevel)
			case "panic":
				logrus.SetLevel(logrus.PanicLevel)
			default:
				return fmt.Errorf("unrecognized log output level")
			}

			conf.SetConfigFile(cmd.Flag("config").Value.String())
			if err := conf.ReadInConfig(); nil != err {
				return err
			}

			conf.OnConfigChange(func(event fsnotify.Event) {
				err := conf.ReadInConfig()
				log.Info("reload config ", event.Name, ", option ", event.Op, ", error: ", err)
			})

			conf.WatchConfig()

			return nil
		},
	}

	RootCommand.PersistentFlags().String("log", "info", "setting log output level.")
	RootCommand.PersistentFlags().StringP("config", "c", ".deploy.yaml", "specify the config path.")

	RootCommand.DisableAutoGenTag = true

	return RootCommand
}
