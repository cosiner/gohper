// Package regexp implements some utilitily function for access regexp group values
// via index or name
package regexp

import (
	"regexp"

	"github.com/cosiner/gohper/lib/types"
)

// NilMap is an empty map, it SHOULD NOT BE MODIFYED for common shared
var NilMap = make(map[string]string)

// Regexp is a package type of regexp.Regexp that provide some useful function
// to access grouped value via index or group name in regexp's Find* result
type Regexp struct {
	*regexp.Regexp
}

// WrapRegexp wrap standard regexp
func WrapRegexp(r *regexp.Regexp) *Regexp {
	return &Regexp{r}
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

// SubexpNamesMap return regexp variable names and index map
func (r *Regexp) SubexpNamesMap() map[string]int {
	names := r.SubexpNames()[1:]
	mp := make(map[string]int, len(names))
	for i, name := range names {
		mp[name] = i
	}
	return mp
}

// SingleSubmatch return first matched groups, remove first whole matched string
func (r *Regexp) SingleSubmatch(s string) (m []string, match bool) {
	if m = r.FindStringSubmatch(s); m != nil {
		m, match = m[1:], true
	}
	return
}

// SingleSubmatchAtIndex return single matched string in given index, if index is
// 0, return the whole match string
func (r *Regexp) SingleSubmatchAtIndex(s string, index int) (ms string, match bool) {
	if index >= 0 && index < len(r.SubexpNames()) {
		if m := r.FindStringSubmatch(s); m != nil {
			ms, match = m[index], true
		}
	}
	return
}

// AllSubmatch return all matched groups
func (r *Regexp) AllSubmatch(s string) (res [][]string, match bool) {
	if m := r.FindAllStringSubmatch(s, -1); m != nil {
		for _, sm := range m {
			res = append(res, sm[1:])
		}
		match = true
	}
	return
}

// AllSubmatchAtIndex return all matched group at gived index
// if index is 0, return whole matched string
func (r *Regexp) AllSubmatchAtIndex(s string, index int) (res []string, match bool) {
	if index >= 0 && index < len(r.SubexpNames()) {
		if m := r.FindAllStringSubmatch(s, -1); m != nil {
			for _, sm := range m {
				res = append(res, sm[index])
			}
			match = true
		}
	}
	return
}

// SingleSubmatchMap return first group string in map
func (r *Regexp) SingleSubmatchMap(s string) (matchMap map[string]string, match bool) {
	matchMap = NilMap
	if m := r.FindStringSubmatch(s); m != nil {
		matchNames := r.SubexpNames()
		matchMap = make(map[string]string, len(matchNames))
		for i, val := range m[1:] { // remove first whole matched string
			matchMap[matchNames[i+1]] = val // matchNames start at index 1
		}
		match = true
	}
	return
}

// SingleSubmatchWithName return first match string with the name
func (r *Regexp) SingleSubmatchWithName(s, name string) (res string, match bool) {
	if index := types.StringIn(name, r.SubexpNames()); index >= 0 {
		if m := r.FindStringSubmatch(s); m != nil {
			res, match = m[index], true
		}
	}
	return
}

// AllSubmatchMap return all matched group with group name
func (r *Regexp) AllSubmatchMap(s string) (matchMaps []map[string]string, match bool) {
	if m := r.FindAllStringSubmatch(s, -1); m != nil {
		matchNames := r.SubexpNames()
		l := len(matchNames)
		for _, singleMatch := range m {
			matchMap := make(map[string]string, l-1) // first matchName is null
			for j, val := range singleMatch[1:] {    // remove first whole matched string
				matchMap[matchNames[j+1]] = val // matchnames start at index 1
			}
			matchMaps = append(matchMaps, matchMap)
		}
		match = true
	}
	return
}

// AllSubmatchWithName return all matched group with group name
// if name is "", return the whole matched string
// else return single matched group string with the name
func (r *Regexp) AllSubmatchWithName(s, name string) (matchs []string, match bool) {
	if index := types.StringIn(name, r.SubexpNames()); index >= 0 {
		if m := r.FindAllStringSubmatch(s, -1); m != nil {
			for _, singleMatch := range m {
				matchs = append(matchs, singleMatch[index])
			}
			match = true
		}
	}
	return
}
