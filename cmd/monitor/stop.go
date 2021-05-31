package monitor

import (
	"github.com/ch629/orca/cmd/config"
	"log"
	"os"
)

type Stop struct {
	config *config.Config
}

func NewStop() Stop {
	return Stop{
		config: config.DefaultConfig,
	}
}

func (s Stop) Run() {
	pid := s.getPid()
	if pid < 0 {
		return
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		log.Fatal("failed to find process", err)
	}
	if err = proc.Kill(); err != nil {
		log.Fatal("failed to kill process", err)
	}

	if _, err = proc.Wait(); err != nil {
		log.Fatal("failed while waiting for process to end", err)
	}
	if err = s.config.DeleteMonitorPid(); err != nil {
		log.Fatal(err)
	}
}

func (s Stop) getPid() int {
	return s.config.GetMonitorPid()
}
