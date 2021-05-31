package monitor

import (
	"github.com/ch629/orca/cmd/config"
	"log"
)

type (
	Remove struct {
		config *config.Config
	}
)

func NewRemove() Remove {
	return Remove{
		config: config.DefaultConfig,
	}
}

func (r Remove) Run(name string) {
	conf, err := r.config.GetMonitorConfig()
	if err != nil {
		log.Fatal(err)
	}

	newMonitors := make([]config.Monitor, 0, len(conf.Monitors))

	for _, monitor := range conf.Monitors {
		if monitor.Name != name {
			newMonitors = append(newMonitors, monitor)
		}
	}
	conf.Monitors = newMonitors
	if err = r.config.WriteMonitors(*conf); err != nil {
		log.Fatal(err)
	}
}
