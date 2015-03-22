package config

import (
	"strings"

	"github.com/cosiner/golib/sys"

	. "github.com/cosiner/golib/errors"

	"github.com/cosiner/golib/types"
)

// lineConfig parse format:key1=val1&key2=val2
// there is no different section
type lineConfig struct {
	values map[string]string
}

// NewLineConfig return a line config parser
func NewLineConfig() ConfigParser {
	return &lineConfig{
		values: make(map[string]string),
	}
}

// ParseString parse string
func (lc *lineConfig) ParseString(str string) (err error) {
	fields := strings.Split(str, "&")
	for _, f := range fields {
		pair := types.ParsePair(f, "=")
		if err = pair.TrimQuote(); err == nil {
			if pair.NoKey() {
				if pair.HasValue() {
					err = Errorf("Wrong Format:%s", f)
					break
				}
			} else {
				lc.values[pair.Key] = pair.Value
			}
		}
	}
	return nil
}

// ParseFile parse first line of file
func (lc *lineConfig) ParseFile(confFileName string) (err error) {
	var line string
	if line, err = sys.ReadFirstLine(confFileName); err == nil {
		err = lc.ParseString(line)
	}
	return
}

// SetCurrSec set current section
func (lc *lineConfig) SetCurrSec(section string) {
}

// DefSec return default section name
func (lc *lineConfig) DefSec() string {
	return ""
}

// CurrSec return current section name
func (lc *lineConfig) CurrSec() string {
	return ""
}

// ValFrom return value from section with gived key
func (lc *lineConfig) ValFrom(key, section string) (val string, has bool) {
	val, has = lc.values[key]
	return
}

// HasSection check whether section exist in config
func (lc *lineConfig) HasSection(section string) bool {
	return true
}

// SectionVals return all key-value pairs from section
func (lc *lineConfig) SectionVals(section string) map[string]string {
	return lc.values
}

// SetSectionVals set replace all section values
func (lc *lineConfig) SetSectionVals(section string, values map[string]string) {
	lc.values = values
}
