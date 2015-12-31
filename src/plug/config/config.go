package config

import (
	"errors"
	"fmt"
)

type Config struct {
	Type       string
	Connection Connection
	Identity   Identity
	Channels   []string
}
type Connection struct {
	Addr         string
	Secure       bool
	Fingerprints []string
}
type Identity struct {
	Nick     []string
	Password string
	RealName string
}

func (c *Config) Validate(k string) error {
	if c.Connection.Addr == "" {
		return errors.New(fmt.Sprintf("Plug '%s' Connection.Addr is empty", k))
	}
	if c.Type == "" {
		return errors.New(fmt.Sprintf("Plug '%s' Type is empty", k))
	}
	if len(c.Identity.Nick) == 0 {
		return errors.New(fmt.Sprintf("Plug '%s' Identity.Nick must contain one or more nicks to use", k))
	}
	if len(c.Channels) == 0 {
		return errors.New(fmt.Sprintf("Plug '%s' Channels must contain one or more channels to relay", k))
	}
	return nil
}
