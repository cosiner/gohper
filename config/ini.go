package config

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"mlib/util/types"
	"os"
)

// default section, if there is no section,
// all key-value pairs will be placed under this section
const DEFAULT_SECTION string = "global"

// iniConfig implements a ini format parser
type iniConfig struct {
	defSec   string                       // default section
	currSec  string                       // current section
	values   map[string]map[string]string // values : [section]([key]value)
	sections []string                     // all sections
}

// newIniConfig return a ini config parser
func newIniConfig() ConfigParser {
	ic := &iniConfig{
		DEFAULT_SECTION,
		DEFAULT_SECTION,
		make(map[string]map[string]string),
		[]string{},
	}
	return ic
}

//	Clear clear parse results
func (ic *iniConfig) Clear() {
	ic.values = make(map[string]map[string]string)
	ic.sections = []string{}
}

// ValFrom return value from section with gived key
func (ic *iniConfig) ValFrom(key, section string) (val string, has bool) {
	var vals map[string]string
	if vals, has = ic.values[section]; has {
		val, has = vals[key]
	}
	return
}

// SetDefSec set default section
func (ic *iniConfig) SetDefSec(section string) {
	if section != "" {
		ic.defSec = section
	}
}

// SetCurrSec set current section
func (ic *iniConfig) SetCurrSec(section string) {
	if section != "" {
		ic.currSec = section
	}
}

// CurrSec return current section
func (ic *iniConfig) CurrSec() string {
	return ic.currSec
}

// Sections return all sections
func (ic *iniConfig) Sections() []string {
	return ic.sections
}

// SectionVals return all key-value pairs from section
func (ic *iniConfig) SectionVals(section string) map[string]string {
	return ic.values[section]
}

// bindTo bind key-value to section
func (ic *iniConfig) bindTo(key, value, section string) {
	if ic.values[section] == nil {
		ic.values[section] = make(map[string]string)
	}
	ic.values[section][key] = value
}

// bind bind key-value to current section
func (ic *iniConfig) bind(key, value string) {
	if ic.values[ic.currSec] == nil {
		ic.values[ic.currSec] = make(map[string]string)
	}
	ic.values[ic.currSec][key] = value
}

// ParseString parse from string
func (ic *iniConfig) ParseString(content string) error {
	if content == "" {
		return errors.New("No Content")
	}

	return ic.parse(types.StringReader(content))
}

// ParseFile parse from file
func (ic *iniConfig) ParseFile(confFileName string) (err error) {
	f, err := os.Open(confFileName)
	if err == nil {
		err = ic.parse(f)
		err = f.Close()
	}
	return
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
	r := bufio.NewReader(reader)
	for {
		var line []byte
		line, err = r.ReadBytes('\n')
		if err != nil {
			break
		}
		key, value, result := parseLine(line)
		switch result {
		case _SECTION:
			ic.SetCurrSec(key)
			if len(ic.sections) == 0 { // if the section is the first section, set it to default
				ic.SetDefSec(key)
			}
			ic.sections = append(ic.sections, key)
		case _KV:
			if len(ic.sections) == 0 {
				ic.sections = append(ic.sections, DEFAULT_SECTION) //if the first few value has no section, add a default section
			}
			ic.bind(key, value)
		}
	}
	ic.currSec = ic.defSec // restore to default section
	return
}

// parseLine parse ini file line,
// for [section], name returned by key
// for key=value, return key, value
// result represent the line is Section, normal key-value, blank, comment, error
func parseLine(line []byte) (key, value string, result lineType) {
	line = trimComment(line)
	start, sep := 0, len(line)

	switch {
	case sep == 0:
		result = _BLANK
	case line[0] == '[':
		{
			if line[sep-1] != ']' {
				result = _ERRFORMAT
			} else {
				result = _SECTION
				start++
				sep--
			}
		}
	default:
		{
			result = _KV
			sep = bytes.IndexByte(line, '=')
			if sep <= 0 {
				result = _ERRFORMAT
			}
		}
	}
	if result <= _KV { // _KV or _SECTION
		key = string(bytes.TrimSpace(line[start:sep]))
		if result == _KV {
			value = string(bytes.TrimSpace(line[sep+1 : len(line)]))
		}
	}
	return
}

// trimComment trim line comment
func trimComment(line []byte) []byte {
	line = bytes.TrimSpace(line)
	end := len(line)

	for i := 0; i < end; i++ {
		if line[i] == '#' {
			end = i
			break
		}
	}
	return bytes.TrimSpace(line[:end])
}
