package types

import (
	"strings"
)

// Pair represent a key-value pair
type Pair struct {
	Key   string
	Value string
}

// ParisePair seperate string use first seperator string
func ParsePair(str, sep string) Pair {
	return parsePair(str, strings.Index(str, sep))
}

// RParsePair seperate string use last seperator string
func RParsePair(str, sep string) Pair {
	return parsePair(str, strings.LastIndex(str, sep))
}

// ParsePairWith seperate string use given function
func ParsePairWith(str, sep string, sepIndexFn func(string, string) int) Pair {
	return parsePair(str, sepIndexFn(str, sep))
}

// parsePair seperate string at given index
// if key or value is nil, set to ""
func parsePair(str string, index int) Pair {
	var key, value string
	if index > 0 {
		key, value = str[:index], str[index+1:]
	} else if index == 0 {
		key, value = "", str
	} else if index < 0 {
		key, value = str, ""
	}
	return Pair{Key: key, Value: value}
}

// NoKey check whether pair has key or not
func (p Pair) NoKey() bool {
	return p.Key == ""
}

// NoValue check whether pair has value or not
func (p Pair) NoValue() bool {
	return p.Value == ""
}
