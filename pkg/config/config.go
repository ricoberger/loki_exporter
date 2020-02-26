package config

import (
	"io/ioutil"
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

	Metrics struct {
		Labels      bool `yaml:"labels"`
		LabelValues bool `yaml:"labelValues"`
		Queries     bool `yaml:"queries"`
	} `yaml:"metrics"`

	Queries []struct {
		Name      string        `yaml:"name"`
		Query     string        `yaml:"query"`
		Limit     int           `yaml:"limit"`
		Start     time.Duration `yaml:"start"`
		End       time.Duration `yaml:"end"`
		Direction string        `yaml:"direction"`
		Regexp    string        `yaml:"regexp"`
	} `yaml:"queries"`
}

// LoadConfig reads the configuration file and umarshal the data into the config struct
func (c *Config) LoadConfig(file string) error {
	// Set default values
	c.Loki.ListenAddress = "http://localhost:3100"
	c.Loki.BasicAuth.Enabled = false

	c.Metrics.Labels = true
	c.Metrics.LabelValues = true
	c.Metrics.Queries = true

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

		if c.Queries[index].Start == 0 {
			c.Queries[index].Start = time.Hour * -24
		}

		if c.Queries[index].End == 0 {
			c.Queries[index].End = time.Second * 0
		}
	}

	return nil
}
