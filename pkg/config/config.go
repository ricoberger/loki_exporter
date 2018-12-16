package config

import (
	"io/ioutil"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"
)

// Config represents the structur of the configuration file
type Config struct {
	Loki struct {
		ListenAddress string `yaml:"listenAddress"`

		BasicAuth struct {
			Enabled  bool   `yaml:"enabled"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"basicAuth"`
	} `yaml:"loki"`

	Queries []struct {
		Name      string `yaml:"name"`
		Query     string `yaml:"query"`
		Limit     int64  `yaml:"limit"`
		Start     string `yaml:"start"`
		End       string `yaml:"end"`
		Direction string `yaml:"direction"`
		Regexp    string `yaml:"regexp"`
	} `yaml:"queries"`
}

// LoadConfig reads the configuration file and umarshal the data into the config struct
func (c *Config) LoadConfig(file string) error {
	// Set default values
	c.Loki.ListenAddress = "http://localhost:3100"
	c.Loki.BasicAuth.Enabled = false

	// Load configuration file
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return err
	}

	// Set default values for all queries and parse start and end values
	for index := 0; index < len(c.Queries); index++ {
		if c.Queries[index].Limit == 0 {
			c.Queries[index].Limit = -1
		}

		if c.Queries[index].Start == "" {
			c.Queries[index].Start = "-24h"
		}

		if c.Queries[index].End == "" {
			c.Queries[index].End = "0s"
		}

		startDuration, err := time.ParseDuration(c.Queries[index].Start)
		if err != nil {
			return err
		}

		endDuration, err := time.ParseDuration(c.Queries[index].End)
		if err != nil {
			return err
		}

		c.Queries[index].Start = strconv.FormatInt(time.Now().Add(startDuration).UnixNano(), 10)
		c.Queries[index].End = strconv.FormatInt(time.Now().Add(endDuration).UnixNano(), 10)
	}

	return nil
}
