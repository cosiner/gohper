// Package config is a set of config parser
package config

import (
	"mlib/util/types"
	"strconv"
)

// ConfigParser is a interface of actual parser
type ConfigParser interface {
	// ParseString parse from string
	ParseString(content string) error
	// ParseFile parse from file
	ParseFile(confFileName string) (err error)
	// SetDefsec set default section
	SetDefSec(section string)
	// SetCurrSec set current section
	SetCurrSec(section string)
	// CurrSec return current section
	CurrSec() string
	// ValFrom return value from section with gived key
	ValFrom(key, section string) (val string, has bool)
	// SectionVals return all key-value pairs from section
	SectionVals(section string) map[string]string
	// Sections return all sections
	Sections() []string
	// Clear clear parse result
	Clear()
}

// Config implements some common config function
type Config struct {
	ConfigParser
}

// NewConfig return a config parser, default use ini config parser
func NewConfig() *Config {
	return &Config{newIniConfig()}
}

// NewConfigWith use a gived parser
func NewConfigWith(parser ConfigParser) *Config {
	return &Config{parser}
}

// ValFrom return value with gived key in the section
func (c *Config) ValFrom(key string, section string) (string, bool) {
	return c.ConfigParser.ValFrom(key, section)
}

// Val return value from current section
func (c *Config) Val(key string) (string, bool) {
	return c.ValFrom(key, c.CurrSec())
}

// ValDef return value from current section, if no this key, return default value
func (c *Config) ValDef(key string, defVal string) string {
	if val, has := c.Val(key); !has {
		return defVal
	} else {
		return val
	}
}

// BoolValFrom return bool value
func (c *Config) BoolValFrom(key, section string, defaultval bool) bool {
	if v, has := c.ValFrom(key, section); has {
		if b, err := types.Str2Bool(v); err == nil {
			return b
		}
	}
	return defaultval
}

// IntVal return integer value
func (c *Config) IntValFrom(key, section string, defaultval int) int {
	if v, has := c.ValFrom(key, section); has {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultval
}

// BoolVal return bool value from current section
func (c *Config) BoolVal(key string, defaultval bool) bool {
	return c.BoolValFrom(key, c.CurrSec(), defaultval)
}

// IntVal return integer value from current section
func (c *Config) IntVal(key string, defaultval int) int {
	return c.IntValFrom(key, c.CurrSec(), defaultval)
}
