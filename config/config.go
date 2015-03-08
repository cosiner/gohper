// Package config is a set of config parser
package config

import (
	"strconv"

	"github.com/cosiner/golib/types"
)

// ConfigParser is a interface of actual parser
type ConfigParser interface {
	// ParseString parse from string
	ParseString(content string) error
	// ParseFile parse from file
	ParseFile(confFileName string) (err error)
	// SetCurrSec set current section
	SetCurrSec(section string)
	// DefSec return default section name
	DefSec() string
	// CurrSec return current section name
	CurrSec() string
	// ValFrom return value from section with gived key
	ValFrom(key, section string) (val string, has bool)
	// HasSection check whether section exist in config
	HasSection(section string) bool
	// SectionVals return all key-value pairs from section
	SectionVals(section string) map[string]string
	// SetsectionVals set section values
	SetSectionVals(section string, values map[string]string)
}

// Config implements some common config function
type Config struct {
	ConfigParser
}

// ConfigType is supported config file type, currently only ini format
type ConfigType int8

const (
	// INI is the ini/conf format config, multi same section will be merged
	INI ConfigType = iota
	// LINE is the format like :k=v&k=v... config, no different sections
	LINE
)

// NewConfig return a config parser, default use ini config parser
func NewConfig(typ ConfigType) (c *Config) {
	switch typ {
	case INI:
		c = &Config{newIniConfig()}
	case LINE:
		c = &Config{NewLineConfig()}
	}
	return
}

// NewConfigWith use a gived parser
func NewConfigWith(parser ConfigParser) *Config {
	return &Config{parser}
}

// Val return value from current section
func (c *Config) Val(key string) (string, bool) {
	return c.ValFrom(key, c.CurrSec())
}

// ValDef return value from current section, if no this key, return default value
func (c *Config) ValDef(key string, defVal string) (val string) {
	val, has := c.Val(key)
	if !has {
		val = defVal
	}
	return
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
func (c *Config) BoolValDef(key string, defaultval bool) bool {
	return c.BoolValFrom(key, c.CurrSec(), defaultval)
}

// IntVal return integer value from current section
func (c *Config) IntValDef(key string, defaultval int) int {
	return c.IntValFrom(key, c.CurrSec(), defaultval)
}

func (c *Config) UnmarshalCurrSec(v interface{}) error {
	return types.UnmarshalToStruct(c.SectionVals(c.CurrSec()), v)
}
