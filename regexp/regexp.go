// Package rgx implements some utilitily function for access regexp group values
// via index or name
package regexp

import (
	"regexp"

	"github.com/cosiner/golib/types"
)

// Regexp is a package type of regexp.Regexp that provide some useful function
// to access grouped value via index or group name in regexp's Find* result
type Regexp struct {
	*regexp.Regexp
}

// Compile is a package function of regexp.Compile
func Compile(pattern string) (r *Regexp, err error) {
	rgx, err := regexp.Compile(pattern)
	if err == nil {
		r = &Regexp{rgx}
	}
	return
}

// CompilePOSIX is a package function of regexp.CompilePOSIX
func CompilePOSIX(pattern string) (r *Regexp, err error) {
	rgx, err := regexp.CompilePOSIX(pattern)
	if err == nil {
		r = &Regexp{rgx}
	}
	return
}

// MustCompile is a package function of regexp.MustCompile
func MustCompile(pattern string) *Regexp {
	return &Regexp{regexp.MustCompile(pattern)}
}

// SingleSubmatch return first matched groups
func (r *Regexp) SingleSubmatch(s string) []string {
	return r.FindStringSubmatch(s)[1:]
}

// SingleSubmatch return first matched groups
func (r *Regexp) SingleSubmatchAtIndex(s string, index int) string {
	return r.FindStringSubmatch(s)[index]
}

// AllSubmatch return all matched group
func (r *Regexp) AllSubmatch(s string) (res [][]string) {
	for _, match := range r.FindAllStringSubmatch(s, -1) {
		res = append(res, match[1:])
	}
	return
}

// AllSubmatchAtIndex return all matched group at gived index
// if index is 0, return whole matched string
// else return group string at the index
func (r *Regexp) AllSubmatchAtIndex(s string, index int) (res []string) {
	for _, match := range r.FindAllStringSubmatch(s, -1) {
		res = append(res, match[index])
	}
	return
}

// SingleSubmatchMap return first group string in map
func (r *Regexp) SingleSubmatchMap(s string) map[string]string {
	matchNames := r.SubexpNames()
	matchMap := make(map[string]string, len(matchNames))
	for i, val := range r.FindStringSubmatch(s)[1:] { // remove first whole matched string
		matchMap[matchNames[i+1]] = val // matchNames start at index 1
	}
	return matchMap
}

// SingleSubmatchWithName return first match string with the name
func (r *Regexp) SingleSubmatchWithName(s, name string) (res string) {
	if index := types.StringIn(name, r.SubexpNames()); index >= 0 {
		res = r.FindStringSubmatch(s)[index]
	}
	return
}

// AllSubmatchMap return all matched group with group name
func (r *Regexp) AllSubmatchMap(s string) (matchMaps []map[string]string) {
	matchNames := r.SubexpNames()
	l := len(matchNames)
	for _, singleMatch := range r.FindAllStringSubmatch(s, -1) {
		matchMap := make(map[string]string, l-1) // first matchName is null
		for j, val := range singleMatch[1:] {    // remove first whole matched string
			matchMap[matchNames[j+1]] = val // matchnames start at index 1
		}
		matchMaps = append(matchMaps, matchMap)
	}
	return
}

// AllSubmatchMap return all matched group with group name
// if name is "", return the whole matched string
// else return single matched group string with the name
func (r *Regexp) AllSubmatchWithName(s, name string) (matchs []string) {
	if index := types.StringIn(name, r.SubexpNames()); index >= 0 {
		for _, singleMatch := range r.FindAllStringSubmatch(s, -1) {
			matchs = append(matchs, singleMatch[index])
		}
	}
	return
}
