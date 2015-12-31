package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	plug "go-chat-relay/src/plug/config"
	"path/filepath"
)

var current *Config

type Config struct {
	Plug   map[string]plug.Config
	Router Router
}
type Router struct {
	AllowCircularFlowLoops bool
	Flow                   [][2]string
}

func FromFile(path string) {
	abspath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	cf := new(Config)
	if _, err := toml.DecodeFile(abspath, &cf); err != nil {
		panic(err)
	}
	if err := dump(abspath, cf); err != nil {
		panic(err)
	}
	if err := cf.Validate(); err != nil {
		panic(err)
	}
	set(cf)
}

func Current() *Config {
	if current == nil {
		panic("Current configuration is not set!")
	}
	return current
}

func (c *Config) Validate() error {
	for k, p := range c.Plug {
		if err := p.Validate(k); err != nil {
			return err
		}
	}
	if len(c.Router.Flow) == 0 {
		return errors.New("Router.Flow must contain one or more pairs(righthanded) of plugs to route messages")
	}
	for _, v := range c.Router.Flow {
		for _, vv := range v {
			if _, ok := c.Plug[vv]; !ok {
				return errors.New(fmt.Sprintf("Can't find plug '%s' from Router.Flow pair %+v", vv, v))
			}
		}
		if !c.Router.AllowCircularFlowLoops && v[0] == v[1] {
			return errors.New(fmt.Sprintf("Circular relay Router.Flow %+v", v))
		}
	}
	return nil
}

func set(c *Config) { current = c }

func dump(path string, config interface{}) error {
	prettyConfig, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"path":   path,
		"config": string(prettyConfig),
	}).Debug("Loaded config")

	return nil
}
