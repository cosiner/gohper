package types

import (
	"fmt"
	"strings"
)

// Pair represent a key-value pair
type Pair struct {
	Key   string
	Value string
}

// ParisePair seperate string use first seperator string
func ParsePair(str, sep string) *Pair {
	return parsePair(str, strings.Index(str, sep))
}

// RParsePair seperate string use last seperator string
func RParsePair(str, sep string) *Pair {
	return parsePair(str, strings.LastIndex(str, sep))
}

// ParsePairWith seperate string use given function
func ParsePairWith(str, sep string, sepIndexFn func(string, string) int) *Pair {
	return parsePair(str, sepIndexFn(str, sep))
}

// parsePair seperate string at given index
// if key or value is nil, set to ""
func parsePair(str string, index int) *Pair {
	var key, value string
	if index > 0 {
		key, value = str[:index], str[index+1:]
	} else if index == 0 {
		key, value = "", str
	} else if index < 0 {
		key, value = str, ""
	}
	return &Pair{Key: key, Value: value}
}

// String return pair's display string, format as key=value
func (p *Pair) String() string {
	return fmt.Sprint("(%s:%s)", p.Key, p.Value)
}

// Trim trim all space of pair's key and value
func (p *Pair) Trim() *Pair {
	p.Key = strings.TrimSpace(p.Key)
	p.Value = strings.TrimSpace(p.Value)
	return p
}

// NoKey check whether pair has key or not
func (p *Pair) NoKey() bool {
	return p.Key == ""
}

// NoValue check whether pair has value or not
func (p *Pair) NoValue() bool {
	return p.Value == ""
}

// HasKey check whether pair has key or not
func (p *Pair) HasKey() bool {
	return p.Key != ""
}

// HasValue check whether pair has value or not
func (p *Pair) HasValue() bool {
	return p.Value != ""
}

// ValueOrKey return value when value is not "", otherwise return key
func (p *Pair) ValueOrKey() string {
	if p.HasValue() {
		return p.Value
	}
	return p.Key
}

// IntValue convert pair's value to int
func (p *Pair) IntValue() (int, error) {
	return Str2Int(p.Value)
}

// MustIntValue convert pair's value to int, on error panic
func (p *Pair) MustIntValue() int {
	return MustStr2Int(p.Value)
}

// IntValueDef convert pair's value to int, on error use default value
func (p *Pair) IntValueDef(def int) int {
	return Str2IntDef(p.Value, def)
}

// BoolValue convert pair's value to bool
func (p *Pair) BoolValue() (bool, error) {
	return Str2Bool(p.Value)
}

// MustBoolValue convert pair's value to bool, on error panic
func (p *Pair) MustBoolValue() bool {
	return MustStr2Bool(p.Value)
}

// MustBoolValue convert pair's value to bool, on error use default value
func (p *Pair) BoolValueDef(def bool) bool {
	return Str2BoolDef(p.Value, def)
}
