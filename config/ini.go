package config

import (
	"io"
	"os"

	. "github.com/cosiner/golib/errors"

	"github.com/cosiner/golib/sys"
	"github.com/cosiner/golib/types"
)

// default section, if there is no section,
// all key-value pairs will be placed under this section
const GLOBAL_OPTION string = "global"

// iniConfig implements a ini format parser
type iniConfig struct {
	currSec string // current section
	defSec  string
	values  map[string]map[string]string // values : [section]([key]value)
}

// ValFrom return value from section with gived key
func (ic *iniConfig) ValFrom(key, section string) (val string, has bool) {
	var vals map[string]string
	if vals, has = ic.values[section]; has {
		val, has = vals[key]
	}
	return
}

// ParseString parse from string
func (ic *iniConfig) ParseString(content string) error {
	if content == "" {
		return Err("No Content")
	}

	return ic.parse(types.StringReader(content))
}

// ParseFile parse from file
func (ic *iniConfig) ParseFile(confFileName string) (err error) {
	return sys.OpenForRead(confFileName, func(fd *os.File) error {
		return ic.parse(fd)
	})
}

// newIniConfig return a ini config parser
func newIniConfig() ConfigParser {
	return &iniConfig{
		values: make(map[string]map[string]string),
	}
}

// DefSec return first section of config
func (ic *iniConfig) DefSec() string {
	return ic.defSec
}

// CurrSec return current section
func (ic *iniConfig) CurrSec() string {
	return ic.currSec
}

// SetCurrSec set current section
func (ic *iniConfig) SetCurrSec(section string) {
	if section != "" {
		ic.currSec = section
		if ic.defSec == "" {
			ic.defSec = section
		}
	}
}

// HasSection check whether given section exist in config
func (ic *iniConfig) HasSection(section string) bool {
	return ic.values[section] != nil
}

// addsection prepare space for new section if section not exist
func (ic *iniConfig) addSection(section string) *iniConfig {
	if !ic.HasSection(section) {
		ic.values[section] = make(map[string]string)
	}
	return ic
}

// SectionVals return all key-value pairs from section
func (ic *iniConfig) SectionVals(section string) map[string]string {
	return ic.values[section]
}

// bind bind key-value to current section
func (ic *iniConfig) bind(key, value string) {
	if ic.currSec == "" {
		ic.addSection(GLOBAL_OPTION).SetCurrSec(GLOBAL_OPTION)
	}
	ic.SectionVals(ic.CurrSec())[key] = value
}

// lineType is parse result of per line
type lineType int8

const (
	_SECTION   lineType = iota // section : [section]
	_KV                        // key-value : key=value
	_BLANK                     // blank line or comment : #comment
	_ERRFORMAT                 // error format line : [section, =value ...
)

// parse read and parse from a reader
func (ic *iniConfig) parse(reader io.Reader) (err error) {
	err = sys.FilterLine(reader, func(linenum int, line []byte) error {
		pair, result := parseLine(line)
		switch result {
		case _SECTION:
			sec := pair.Key
			ic.addSection(sec).SetCurrSec(sec)
		case _KV:
			ic.bind(pair.Key, pair.Value)
		case _ERRFORMAT:
			return Errorf("Wrong format of line %d: %s", linenum, string(line))
		}
		return nil
	})
	if err == nil {
		ic.SetCurrSec(ic.DefSec())
	}
	return
}

// parseLine parse ini file line,
// for [section], name returned by key
// for key=value, return key, value
// result represent the line is Section, normal key-value, blank, comment, error
func parseLine(line []byte) (pair *types.Pair, result lineType) {
	line = types.TrimBytesAfter(line, []byte("#"))
	if len(line) == 0 {
		result = _BLANK
		return
	}
	pair = types.ParsePair(string(line), "=").Trim()
	if pair.NoKey() {
		result = _ERRFORMAT
	} else {
		key := pair.Key
		if key[0] == '[' {
			last := len(key) - 1
			if key[last] == ']' {
				pair.Key = key[1:last]
				result = _SECTION
			} else {
				result = _ERRFORMAT
			}
		} else {
			if pair.NoValue() {
				result = _ERRFORMAT
			} else {
				result = _KV
			}
		}
	}
	return
}
