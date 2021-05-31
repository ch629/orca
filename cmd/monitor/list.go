package monitor

import (
	"github.com/ch629/orca/cmd/config"
	"log"
)

type List struct {
	config *config.Config
}

func NewList() List {
	return List{
		config: config.DefaultConfig,
	}
}

func (l List) Run() {
	conf, err := l.config.GetMonitorConfig()
	if err != nil {
		log.Fatal(err)
	}

	for _, monitor := range conf.Monitors {
		log.Println(monitor.Name)
	}
}
