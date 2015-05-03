package pair

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosiner/gohper/strings2"
)

// Pair represent a key-value pair
type Pair struct {
	Key   string
	Value string
}

// Parse seperate string use first seperator string
func Parse(str, sep string) *Pair {
	return parse(str, strings.Index(str, sep))
}

// Rparse seperate string use last seperator string
func Rparse(str, sep string) *Pair {
	return parse(str, strings.LastIndex(str, sep))
}

// ParsePairWith seperate string use given function
func ParsePairWith(str, sep string, sepIndexFn func(string, string) int) *Pair {
	return parse(str, sepIndexFn(str, sep))
}

// parse seperate string at given index
// if key or value is nil, set to ""
func parse(str string, index int) *Pair {
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
	return fmt.Sprintf("(%s:%s)", p.Key, p.Value)
}

// Trim trim all space of pair's key and value
func (p *Pair) Trim() *Pair {
	p.Key = strings.TrimSpace(p.Key)
	p.Value = strings.TrimSpace(p.Value)
	return p
}

// TrimQuote trim quote for pair's key and value
func (p *Pair) TrimQuote() (err error) {
	var key, value string
	key, err = strings2.TrimQuote(p.Key)
	if err == nil {
		value, err = strings2.TrimQuote(p.Value)
		if err == nil {
			p.Key, p.Value = key, value
		}
	}
	return
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
	return strconv.Atoi(p.Value)
}

// BoolValue convert pair's value to bool
func (p *Pair) BoolValue() (bool, error) {
	return strconv.ParseBool(p.Value)
}
