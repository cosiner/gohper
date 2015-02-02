package server

import (
	"bytes"
	"strings"

	"github.com/cosiner/golib/regexp"

	. "github.com/cosiner/golib/errors"
)

type (
	// urlType is type of url, determined by which section url start with
	urlType uint8

	// PatternType is type of pattern, which is LITERAL or REGEXP_LITERAL or REGEXP(with named group)
	// or ERROR_FORMAT
	PatternType uint8

	// UrlCompiler is a compiler which compile simpile url regexp syntax to standard
	// golang regexp syntax

	// matchany character ".", default use "N:"
	UrlCompiler struct {
		RegLeft       byte   // RegLeft: left character of regexp group
		RegRight      byte   // RegRight: right character of regexp group
		NeedEscape    bool   // NeedEscape: escape RegLeft and RegRight exist in group regexp
		NoReplace     bool   // NoReplace: replace "." exist in group regexp with given matchany string
		Matchany      string // if matchany is "", use default "[a-zA-Z0-9_]"
		NoreplaceFlag string // NoreplaceFlag appear at the begin of pattern string means don't replace
		init          bool   // init: whether Compiler has been inited
	}
	// VarIndexer is a indexer for regexp variables and values
	VarIndexer interface {
		// VarIndex return variable index in matcher regexp pattern
		VarIndex(name string) int

		// ValuesOf return values of variable in given values
		ValueOf(values []string, name string) string
		// ScanInto scan given values into variable addresses
		// if address is nil, skip it
		ScanInto(values []string, vars ...*string)
	}

	// Matcher is a url matcher
	// to make match more fast, the compiled url pattern will be divided to
	// five types, for each type, Match and MatchOnly will only match
	// different section of url
	Matcher interface {
		// Pattern return the compiled pattern string
		// for regexp matcher, it return parsed regexp string
		// for literal matcher, it return the original pattern string
		Pattern() string
		// IsLiteral check whether this mathcer is just match literal pattern
		IsLiteral() bool
		// Match regexp matcher
		Match(path string) (values []string, match bool)
		// MatchOnly do only match, don't extract url variable vlaues
		MatchOnly(path string) bool
		// PrefixMatchOnly perform prefix match, only used for literal matcher
		PrefixMatchOnly(path string) bool
		//
		VarIndexer
	}

	// matcher is the actual url matcher
	matcher struct {
		// pattern is url pattern after parse
		pattern string
		// regPattern is compiled regexp pattern, only used when pattern is not literal
		regPattern *regexp.Regexp
		regVars    map[string]int
		// matchFunc
		matchFunc     func(url string) ([]string, bool)
		matchOnlyFunc func(url string) bool
	}
)

const (
	// pattern type
	LITERAL        PatternType = iota // LITERAL means a normal string without regexp
	REGEXP_LITERAL                    //REGXP_LITERAL means a regexp string without groups
	REGEXP                            // REGEXP means a regexp with groups
	ERROR_FORMAT                      // ERROR_FORMAT means pattern format is wrong

	// url type
	_HOST urlType = iota // _HOST means pattern is start with url host, such as google.com
	_PATH                // _PATH means pattern is start with url path, such as /user

	// parse state
	_START      // _START means ready to match '{''
	_INREGGROUP // _INREGGROUP means already matched '{', ready to match '}'

	// regexp group state
	_INGROUP // _INGROUP means match group name
	_INREGEX // _INREGEX means match regexp string
)

var (
	standardCompiler = &UrlCompiler{}
	nonLiteralError  = Err("Not a literal pattern")
)

// NewMatcher create a new matcher use simple url syntax, for more detail
// see package doc
// if want to match fragment, the pattern should start with "#", example:#{section}
// if want to match query, the pattern should start with "?", example:?id={id:\d+}
// if want to match path, the pattern should start with "/", example:/user/{id:\d+}
// in other condition, pattern will auto-match host, example:{site}\.com/user/123
// scheme is not supported
func NewMatcher(urlPattern string) (Matcher, error) {
	return NewMatcherWith(urlPattern, standardCompiler)
}

// NewLiteralMatcher create a literal matcher
func NewLiteralMatcher(urlPattern string) (Matcher, error) {
	return NewLiteralMatcherWith(urlPattern, standardCompiler)
}

// NewLiteralMatcherWith create a literal matcher
func NewLiteralMatcherWith(urlPattern string, compiler *UrlCompiler) (Matcher, error) {
	typ, pat := compiler.Compile(urlPattern)
	if typ != LITERAL {
		return nil, nonLiteralError
	}
	return NewStandardMatcher(typ, pat)
}

// NewMatcherWith use customed compiler to create matcher
func NewMatcherWith(urlPattern string, compiler *UrlCompiler) (Matcher, error) {
	typ, pat := compiler.Compile(urlPattern)
	return NewStandardMatcher(typ, pat)
}

// NewMatcher create a new matcher, the regexpPattern must be standard golang regexp format
// if want another simple syntax, please see "Compile" function for detail
// patternType can be one of three type, each type's match algorighm is different
// LITERAL, a literal pattern without regexp
// REGEXP_LITERAL, a regexp pattern without named group, no need to capture values
// REGEXP, a regexp pattern with named group, need to capture values
// see NewMatcher for url match rules
func NewStandardMatcher(patternType PatternType, pattern string) (
	Matcher, error) {
	if pattern == "" {
		return nil, Err("No content")
	} else if patternType == ERROR_FORMAT {
		return nil, Errorf("Wrong format:", pattern)
	}
	var err error
	m := &matcher{
		pattern: pattern,
	}
	if patternType == LITERAL { // literal pattern
		m.matchFunc = m.literalMatch
		m.matchOnlyFunc = m.literalMatchOnly
	} else if m.regPattern, err = regexp.Compile(pattern); err == nil {
		if patternType == REGEXP { // regexp pattern with group
			m.matchFunc = m.regexpMatch
			m.regVars = m.regPattern.SubexpNamesMap()
		} else {
			m.matchFunc = m.regexpNoGroupMatch // regexp pattern without group
		}
		m.matchOnlyFunc = m.regexpMatchOnly
	}
	return m, err
}

// Pattern return url pattern after parse
func (m *matcher) Pattern() string {
	return m.pattern
}

// IsLiteral check whether it's a literal matcher
func (m *matcher) IsLiteral() bool {
	return m.regPattern == nil
}

// VarIndex return variable index in matcher regexp pattern
func (m *matcher) VarIndex(name string) (index int) {
	index = -1
	if vars := m.regVars; vars != nil {
		if i, has := vars[name]; has {
			index = i
		}
	}
	return
}

// ValuesOf return values of variable in given values
func (m *matcher) ValueOf(values []string, name string) string {
	if index := m.VarIndex(name); index != -1 {
		return values[index]
	}
	return ""
}

// ScanInto scan given values into variable addresses
// if address is nil, skip it
func (*matcher) ScanInto(values []string, vars ...*string) {
	l1, l2 := len(values), len(vars)
	for i := 0; i < l1 && i < l2; i++ {
		if vars[i] != nil {
			*vars[i] = values[i]
		}
	}
}

// Match match given url path, reutrn matched values and if it is match or not
// if matcher is a literal matcher, no matched value will be returned
func (m *matcher) Match(path string) ([]string, bool) {
	return m.matchFunc(path)
}

// MatchOnly only return whether url is match, don't extract url variable values
func (m *matcher) MatchOnly(path string) bool {
	return m.matchOnlyFunc(path)
}

// literalMatch match literal pattern
func (m *matcher) literalMatch(path string) ([]string, bool) {
	return nil, m.literalMatchOnly(path)
}

// literalMatchOnly match literal patterh
func (m *matcher) literalMatchOnly(path string) bool {
	return m.pattern == path
}

// PrefixMatchOnly perform prefix match, only for literal matcher
func (m *matcher) PrefixMatchOnly(path string) bool {
	matcherPattern := m.pattern
	l1, l2 := len(matcherPattern), len(m.pattern)
	return l1 == l2 || (l1 < l2 && path[l1] == '/')
}

// regexpNoGroupMatch match no named group regexp
func (m *matcher) regexpNoGroupMatch(path string) ([]string, bool) {
	return nil, m.regexpMatchOnly(path)
}

// regexpMatch match regexp with named group
func (m *matcher) regexpMatch(path string) ([]string, bool) {
	return m.regPattern.SingleSubmatch(path)
}

// regexpMatchOnly only match regexp, don't capture values
func (m *matcher) regexpMatchOnly(path string) bool {
	return m.regPattern.MatchString(path)
}

// init init compiler
func (c *UrlCompiler) Init() {
	if c.init {
		return
	}
	if c.RegLeft == 0 {
		c.RegLeft = '{'
	}
	if c.RegRight == 0 {
		c.RegRight = '}'
	}
	if c.Matchany == "" {
		if c.NoReplace {
			c.Matchany = "."
		} else {
			c.Matchany = "[a-zA-Z0-9_]"
		}
	}
	if c.NoreplaceFlag == "" {
		c.NoreplaceFlag = "N:"
	}
}

// Compile convert a simple url regexp syntax to standard golang regexp syntax
//
// regexp is all quoted by "{}", otherwise, it will be treated as literal
// if regexp contains "{" or "}", just use "\{", "\}" to make a escape
// note:only '{' and '}' will be escaped
// if don't like "{" or "}", just make your own two different character
// and it's needEscape if it's conflict with standard regexp character
//
// if need capture value, just use "{name:regexp}", if not, use {:regexp},
// {} or {:} is wrong format
// {name} is default means capture a section of url seperated by "/", it's default
// regexp will be [a-zA-Z0-9_]*
//
// all "." exist in regexp will be replaced with "[a-zA-Z0-9_]" to limit this regexp
// only used in one url section, if don't need this replacement, just place an
// "N:" before your regexp, or use custom own UrlCompiler
// because of Compile prepend a "^" and append a "$" to regexp result
// so it only can do full match
func (c *UrlCompiler) Compile(urlPattern string) (PatternType, string) {
	c.Init()
	if urlPattern == "" {
		return ERROR_FORMAT, ""
	}
	var (
		patternType            = LITERAL
		buf                    = bytes.NewBuffer(make([]byte, 0, len(urlPattern)+10))
		parseState, groupState = _START, _INGROUP
		group, regexp          []byte
		replace, matchany      = !c.NoReplace, []byte(c.Matchany)
		needEscape             = c.NeedEscape
		regLeft, regRight      = c.RegLeft, c.RegRight
		diff                   = (regLeft != regRight)
	)
	if flag := c.NoreplaceFlag; strings.HasPrefix(urlPattern, flag) {
		urlPattern = urlPattern[len(flag):]
		replace = false
	}
	buf.WriteByte('^')
	for i, l := 0, len(urlPattern); i < l; i++ {
		c := byte(urlPattern[i])
		if parseState == _START {
			if diff && c == regRight {
				goto ERROR
			} else if c == regLeft {
				patternType, parseState, group, regexp =
					REGEXP_LITERAL, _INREGGROUP, group[:0], regexp[:0]
			} else {
				buf.WriteByte(c)
			}
		} else {
			if diff && c == regLeft {
				goto ERROR
			} else if c == regRight {
				success, hasGroup := writeRegexp(buf, replace, group, regexp, matchany)
				if !success {
					goto ERROR
				} else if hasGroup {
					patternType = REGEXP
				}
				parseState, groupState = _START, _INGROUP
			} else if c == ':' {
				groupState = _INREGEX
			} else if groupState == _INGROUP {
				group = append(group, c)
			} else {
				if needEscape && c == '\\' { // escape
					i++
					if i == l {
						goto ERROR
					}
					c = byte(urlPattern[i])
					if c != regLeft && c != regRight {
						regexp = append(regexp, '\\')
					}
				}
				regexp = append(regexp, c)
			}
		}
	}
	if patternType != LITERAL {
		buf.WriteByte('$')
		urlPattern = buf.String()
	}
	goto END
ERROR:
	patternType = ERROR_FORMAT
END:
	return patternType, urlPattern
}

// writeRegexp write group and regexp as a grouped regexp to buffer
// if group and regexp is both empty, write failed
func writeRegexp(buf *bytes.Buffer, replace bool,
	group, regexp, replaceMatchany []byte) (success, hasGroup bool) {

	groupLen, regexpLen := len(group), len(regexp)
	if success = (groupLen != 0 || regexpLen != 0); success {
		if replace {
			regexp = bytes.Replace(regexp, []byte("."), replaceMatchany, -1)
		}
		if groupLen == 0 {
			buf.Write(regexp)
		} else {
			hasGroup = true
			buf.WriteString("(?P<")
			buf.Write(group)
			buf.WriteByte('>')
			if regexpLen == 0 {
				buf.Write(replaceMatchany)
				buf.WriteByte('*')
			} else {
				buf.Write(regexp)
			}
			buf.WriteByte(')')
		}
	}
	return
}
