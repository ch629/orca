package config

import (
	"errors"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"path"
	"strconv"
	"sync"
)

type (
	Config struct {
		fs            afero.Fs
		configDirInit sync.Once
		configDir     string
		monitorViper  *viper.Viper
	}

	MonitorConfig struct {
		Monitors []Monitor `mapstructure:"monitors"`
	}

	Monitor struct {
		Name          string `mapstructure:"name" yaml:"name"`
		Url           string `mapstructure:"url" yaml:"url"`
		Interval      int    `mapstructure:"interval" yaml:"interval"`
		Retries       int    `mapstructure:"retries" yaml:"retries"`
		RetryInterval int    `mapstructure:"retry-interval" yaml:"retry-interval"`
	}
)

var DefaultConfig = &Config{
	fs:           afero.NewOsFs(),
	monitorViper: viper.New(),
}

const (
	monitorPidFile    = ".monitor"
	monitorConfigFile = "monitor.yml"
)

// GetMonitorConfig returns the data in the monitoring config file
func (c *Config) GetMonitorConfig() (*MonitorConfig, error) {
	dir := c.ConfigDir()
	c.monitorViper.SetFs(c.fs)
	c.monitorViper.AddConfigPath(dir)
	c.monitorViper.SetConfigName(monitorConfigFile)
	c.monitorViper.SetConfigType("yml")
	if err := c.monitorViper.ReadInConfig(); err != nil {
		return nil, err
	}
	var config MonitorConfig
	if err := c.monitorViper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

// DeleteMonitorPid deletes the file storing the monitoring PID
func (c *Config) DeleteMonitorPid() error {
	return c.fs.Remove(path.Join(c.ConfigDir(), monitorPidFile))
}

// GetMonitorPid returns the Process ID of the monitor
// returns -1 if none exists
// TODO: Should this return an err instead of -1?
func (c *Config) GetMonitorPid() int {
	dir := c.ConfigDir()
	f, err := c.fs.Open(path.Join(dir, monitorPidFile))
	if err != nil {
		return -1
	}
	defer f.Close()
	bs, err := io.ReadAll(f)
	if err != nil {
		return -1
	}

	i, err := strconv.Atoi(string(bs))
	if err != nil {
		return -1
	}
	return i
}

// WriteMonitorPid attempts to write the monitoring Process ID to a file for persistence
func (c *Config) WriteMonitorPid(pid int) (err error) {
	monitorFilePath := path.Join(c.ConfigDir(), monitorPidFile)
	var f afero.File
	// Try to create the file
	if f, err = c.fs.Create(monitorFilePath); errors.Is(err, os.ErrExist) {
		// If it already exists, just open it
		if f, err = c.fs.OpenFile(monitorFilePath, os.O_TRUNC, 0); err != nil {
			return
		}
	} else if err != nil {
		return
	}
	defer f.Close()
	// Write the PID to the file
	_, err = io.WriteString(f, strconv.Itoa(pid))
	return
}

// WriteMonitors will write the config to a file, overwriting if it already exists
func (c *Config) WriteMonitors(config MonitorConfig) (err error) {
	monitorConfigPath := path.Join(c.ConfigDir(), monitorConfigFile)
	var f afero.File
	// Try to create the file
	if f, err = c.fs.Create(monitorConfigPath); errors.Is(err, os.ErrExist) {
		// If it already exists, just open it
		if f, err = c.fs.OpenFile(monitorConfigPath, os.O_TRUNC, 0); err != nil {
			return
		}
	} else if err != nil {
		return
	}
	defer f.Close()
	var bs []byte
	if bs, err = yaml.Marshal(config); err != nil {
		return
	}
	_, err = f.Write(bs)
	return
}
