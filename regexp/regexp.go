// Package regexp implements some utilitily function for access regexp group values
// via index or name
//
// Index:                 0                           1        2             3
// Name:                  ""                       "name1"    "name2"      "name3"
// FirstMatch []string{"THe Whold Matched String", "Other", "SubMatched", "String"}
// AllMatch   []FirstMatch{}
package regexp

import "regexp"

type Regexp struct {
	*regexp.Regexp
	names map[string]int
}

func Wrap(r *regexp.Regexp) *Regexp {
	return &Regexp{Regexp: r}
}

func Compile(expr string) (*Regexp, error) {
	return compile(expr, regexp.Compile)
}

func CompilePOSIX(expr string) (*Regexp, error) {
	return compile(expr, regexp.CompilePOSIX)
}

func compile(expr string, compFunc func(string) (*regexp.Regexp, error)) (*Regexp, error) {
	r, err := compFunc(expr)
	if err == nil {
		return &Regexp{Regexp: r}, nil
	}
	return nil, err
}

func MustCompile(expr string) *Regexp {
	return &Regexp{Regexp: regexp.MustCompile(expr)}
}

func MustCompilePOSIX(expr string) *Regexp {
	return &Regexp{Regexp: regexp.MustCompilePOSIX(expr)}
}

func (r *Regexp) Names() map[string]int {
	if r.names == nil {
		names := r.SubexpNames()
		r.names = make(map[string]int, len(names))
		for i, name := range names {
			r.names[name] = i
		}
	}
	return r.names
}

func (r *Regexp) First(s string) []string {
	return r.FindStringSubmatch(s)
}

func (r *Regexp) ByIndex(s string, index int) string {
	if index >= 0 && index < len(r.Names()) {
		if vals := r.First(s); vals != nil {
			return vals[index]
		}
	}
	return ""
}

func (r *Regexp) ByName(s, name string) string {
	names := r.Names()
	if index, has := names[name]; has {
		return r.ByIndex(s, index)
	}
	return ""
}

func (r *Regexp) All(s string) [][]string {
	return r.FindAllStringSubmatch(s, -1)
}

func (r *Regexp) AllByIndex(s string, index int) []string {
	if index >= 0 && index < len(r.Names()) {
		if vals := r.All(s); vals != nil {
			res := make([]string, len(vals))
			for i := range res {
				res[i] = vals[i][index]
			}
			return res
		}
	}
	return nil
}

func (r *Regexp) AllByName(s, name string) []string {
	names := r.Names()
	if index, has := names[name]; has {
		return r.AllByIndex(s, index)
	}
	return nil
}
