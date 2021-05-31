package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ch629/orca/cmd/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Start struct {
	config *config.Config
}

func NewStart() Start {
	return Start{
		config: config.DefaultConfig,
	}
}

func (s Start) Run(ctx context.Context) {
	// TODO: Listen for config updates
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, os.Kill)

	// Stop process on interrupt or kill
	go func() {
		<-signals
		cancel()
	}()

	// Don't run a monitor if one is already running
	if s.isAlreadyRunning() {
		log.Println("Already running")
		return
	}

	pid := os.Getpid()
	log.Println("Started monitoring on PID:", pid)

	if err := s.config.WriteMonitorPid(pid); err != nil {
		log.Fatal("Failed to write monitor pid to file", err)
	}

	defer s.config.DeleteMonitorPid()

	var wg sync.WaitGroup
	monitors := s.getMonitors()
	wg.Add(len(monitors))

	// Setup goroutines to run each monitor
	for _, monitor := range monitors {
		go func(monitor Monitor) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case <-time.Tick(monitor.Interval):
				}

				if err := monitor.Check(ctx); err != nil {
					log.Println("Status check failed", err)
					monitor.onFailure()
				}
			}
		}(monitor)
	}

	wg.Wait()
}

func (s Start) isAlreadyRunning() bool {
	pid := s.config.GetMonitorPid()
	if pid == -1 {
		return false
	}
	// TODO: interface to test this
	proc, err := os.FindProcess(pid)
	return err == nil && proc != nil
}

func (s Start) getMonitors() []Monitor {
	conf, err := s.config.GetMonitorConfig()
	if err != nil {
		// TODO: Return instead
		log.Fatal("error getting monitors", err)
	}
	monitors := make([]Monitor, len(conf.Monitors))
	for i, monitor := range conf.Monitors {
		monitors[i] = Monitor{
			Name:          monitor.Name,
			Interval:      time.Duration(monitor.Interval) * time.Second,
			URL:           monitor.Url,
			statusFunc:    DefaultStatusFunc,
			Retries:       monitor.Retries,
			RetryInterval: time.Duration(monitor.RetryInterval) * time.Second,
			Timeout:       5 * time.Second,
			onFailure: func() {
				log.Printf("Monitor %v Failed\n", monitor.Name)
			},
		}
	}
	return monitors[:]
}

// CheckOnce checks the monitor once without retries
func (monitor Monitor) CheckOnce(ctx context.Context) (err error) {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(monitor.Timeout))
	defer cancel()

	// Send request to the monitor endpoint
	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, monitor.URL, nil); err != nil {
		return
	}
	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	// Read the JSON response as a map
	var body map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return
	}

	// Check the response
	if err = monitor.statusFunc(resp.StatusCode, body); err != nil {
		err = fmt.Errorf("status check failed due to %w", err)
	}
	return
}

// Check checks the monitor with retries
func (monitor Monitor) Check(ctx context.Context) (err error) {
	for i := 0; i < monitor.Retries; i++ {
		//  Status is UP
		if err = monitor.CheckOnce(ctx); err == nil {
			return
		}

		// Wait between retries
		if i < monitor.Retries {
			select {
			// Don't retry if we've been cancelled
			case <-ctx.Done():
				return
			case <-time.Tick(monitor.RetryInterval):
			}
		}
	}
	return
}
