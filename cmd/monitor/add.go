package monitor

import (
	"github.com/ch629/orca/cmd/config"
	"github.com/spf13/cobra"
	"log"
)

type (
	Add struct {
		config *config.Config
	}
	AddFlags struct {
		Name          string
		Url           string
		Interval      int
		Retries       int
		RetryInterval int
	}
)

func (flags *AddFlags) SetupFlags(command *cobra.Command) {
	command.Flags().StringVarP(&flags.Name, "name", "n", "", "Name of the monitor")
	command.Flags().StringVarP(&flags.Url, "url", "u", "", "Url of status endpoint")
	command.Flags().IntVarP(&flags.Interval, "interval", "i", 10, "Interval between status checks in seconds")
	command.Flags().IntVarP(&flags.Retries, "retries", "r", 3, "Amount of retries before considering a failure")
	command.Flags().IntVarP(&flags.RetryInterval, "retry-interval", "R", 10, "Interval between retries in seconds")

	command.MarkFlagRequired("name")
	command.MarkFlagRequired("url")
}

func NewAdd() Add {
	return Add{
		config: config.DefaultConfig,
	}
}

func (a Add) Run(flags AddFlags) {
	conf, err := a.config.GetMonitorConfig()
	if err != nil {
		log.Fatal(err)
	}

	conf.Monitors = append(conf.Monitors, config.Monitor{
		Name:          flags.Name,
		Interval:      flags.Interval,
		Url:           flags.Url,
		Retries:       flags.Retries,
		RetryInterval: flags.RetryInterval,
	})

	if err = a.config.WriteMonitors(*conf); err != nil {
		log.Fatal(err)
	}
}
