package server

import (
	"bytes"
	"net/url"
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
	UrlCompiler struct {
		regLeft         byte
		regRight        byte
		needEscape      bool
		replaceMatchany bool
		matchany        []byte
		noreplaceFlag   string
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
		// regexp matcher
		Match(url *url.URL) (urlVars map[string]string, match bool)
		// MatchOnly do only match, don't extract url variable vlaues
		MatchOnly(url *url.URL) bool
	}

	// matcher is the actual url matcher
	matcher struct {
		// pattern is url pattern after parse
		pattern string
		// regPattern is compiled regexp pattern, only used when pattern is not literal
		regPattern *regexp.Regexp
		// urlPatternBuilder build a same pattern format with
		// the compiled pattern from a url
		// such as the compiled pattern is google.com/{sub:.*}
		urlPatternBuilder func(url *url.URL) string
		// matchFunc
		matchFunc     func(url string) (map[string]string, bool)
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
	_HOST  urlType = iota // _HOST means pattern is start with url host, such as google.com
	_PATH                 // _PATH means pattern is start with url path, such as /user
	_QUERY                // _QUERY means pattern is start with url query such as ?id=1
	_FRAG                 // _FRAG means pattern is start with page fragment such as #section1

	// parse state
	_START      // _START means ready to match '{''
	_INREGGROUP // _INREGGROUP means already matched '{', ready to match '}'

	// regexp group state
	_INGROUP // _INGROUP means match group name
	_INREGEX // _INREGEX means match regexp string

	// regexp seperator
	_REGLEFT  = '{' // _REGLEFT is the default regexp left letter
	_REGRIGHT = '}' // _REGRIGHT is the default regexp right letter

	// flag that don't replace match-any character exist in regexp
	_NOREPLACE_MATCHANY = "N:"
	_MATCHANY           = "[a-zA-Z0-9_]"
)

var (
	_ORIGIN_MATCHANY = []byte(".")
	// _MATCHANY is defalut match-any letter that replace '.'
	// _ORIGIN_MATCHANY is original match-any letter
	// StandardCompiler is the standard regexp compiler
	// default value: regleft:'{', regright:'}',
	// needEscape:true, replaceMatchany:true
	// matchany:"[a-zA-Z0-9_]", noreplaceFlag:"N:"
	StandardCompiler = NewCompiler('{', '}', true, true, "", "")
)

// NewCompiler create a new compiler
// regLeft: left character of regexp group
// regRight: right character of regexp group
// needEscape: escape regLeft and regRight exist in group regexp
// replaceMatchany: replace "." exist in group regexp with given matchany string
// if matchany is "", use default "[a-zA-Z0-9_]"
// noreplaceFlag appear at the begin of pattern string means don't replace
// matchany character ".", default use "N:"
func NewCompiler(regLeft, regRight byte, needEscape, replaceMatchany bool,
	matchany, noreplaceFlag string) *UrlCompiler {

	if matchany == "" {
		if replaceMatchany {
			matchany = _MATCHANY
		} else {
			matchany = "."
		}
	}
	if noreplaceFlag == "" {
		noreplaceFlag = _NOREPLACE_MATCHANY
	}
	return &UrlCompiler{
		regLeft:         regLeft,
		regRight:        regRight,
		needEscape:      needEscape,
		replaceMatchany: replaceMatchany,
		matchany:        []byte(matchany),
		noreplaceFlag:   noreplaceFlag,
	}
}

// NewMatcher create a new matcher use simple url syntax, for more detail
// see package doc
// if want to match fragment, the pattern should start with "#", example:#{section}
// if want to match query, the pattern should start with "?", example:?id={id:\d+}
// if want to match path, the pattern should start with "/", example:/user/{id:\d+}
// in other condition, pattern will auto-match host, example:{site}\.com/user/123
// scheme is not supported
func NewMatcher(urlPattern string) (Matcher, error) {
	typ, pat := StandardCompiler.Compile(urlPattern)
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
	m := new(matcher).setPatternBuilder(pattern)
	m.pattern = pattern
	if patternType == LITERAL { // literal pattern
		m.matchFunc = m.literalMatch
		m.matchOnlyFunc = m.literalMatchOnly
	} else if m.regPattern, err = regexp.Compile(pattern); err == nil {
		if patternType == REGEXP { // regexp pattern with group
			m.matchFunc = m.regexpMatch
		} else {
			m.matchFunc = m.regexpNoGroupMatch // regexp pattern without group
		}
		m.matchOnlyFunc = m.regexpMatchOnly
	}
	return m, err
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
func (c *UrlCompiler) Compile(urlPattern string) (PatternType, string) {
	if urlPattern == "" {
		return ERROR_FORMAT, ""
	}
	var (
		patternType                      = LITERAL
		buf                              = bytes.NewBuffer(make([]byte, 0, len(urlPattern)+10))
		parseState, groupState           = _START, _INGROUP
		group, regexp                    []byte
		replace, matchany, noreplaceFlag = c.replaceMatchany, c.matchany, c.noreplaceFlag // replace means replace '.' exist in regexp to "[a-zA-Z0-9_]"
		needEscape                       = c.needEscape
		regLeft, regRight                = c.regLeft, c.regRight
		diff                             = (regLeft != regRight)
	)

	if strings.HasPrefix(urlPattern, noreplaceFlag) {
		urlPattern = urlPattern[len(noreplaceFlag):]
		replace = false
	}
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
			regexp = bytes.Replace(regexp, _ORIGIN_MATCHANY, replaceMatchany, -1)
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

// Pattern return url pattern after parse
func (m *matcher) Pattern() string {
	return m.pattern
}

// IsLiteral check whether it's a literal matcher
func (m *matcher) IsLiteral() bool {
	return m.regPattern == nil
}

// Match match given url, reutrn matched values and if it is match or not
// if matcher is a literal matcher, no matched value will be returned
func (m *matcher) Match(url *url.URL) (map[string]string, bool) {
	return m.matchFunc(m.urlPatternBuilder(url))
}

// MatchOnly only return whether url is match, don't extract url stories
func (m *matcher) MatchOnly(url *url.URL) bool {
	return m.matchOnlyFunc(m.urlPatternBuilder(url))
}

// literalMatch match literal pattern
func (m *matcher) literalMatch(pattern string) (vals map[string]string, match bool) {
	return regexp.NIL_MAP, m.literalMatchOnly(pattern)
}

// literalMatchOnly only match literal pattern, don't extract pattern variables
func (m *matcher) literalMatchOnly(pattern string) bool {
	matcherPattern := m.pattern
	return strings.HasPrefix(pattern, matcherPattern) &&
		(len(matcherPattern) == len(pattern) ||
			pattern[len(matcherPattern)] == '/')
}

// regexpNoGroupMatch match no named group regexp
func (m *matcher) regexpNoGroupMatch(pattern string) (map[string]string, bool) {
	return regexp.NIL_MAP, m.regexpMatchOnly(pattern)
}

// regexpMatch match regexp with named group
func (m *matcher) regexpMatch(pattern string) (map[string]string, bool) {
	return m.regPattern.SingleSubmatchMap(pattern)
}

// regexpMatchOnly only match regexp, don't capture values
func (m *matcher) regexpMatchOnly(pattern string) bool {
	return m.regPattern.MatchString(pattern)
}

// setPatternBuilder set up pattern string builder by the start section type of pattern
func (m *matcher) setPatternBuilder(pattern string) *matcher {
	switch checkUrlType(pattern) {
	case _FRAG:
		m.urlPatternBuilder = m.buildWithFrag
	case _QUERY:
		m.urlPatternBuilder = m.buildWithQuery
	case _PATH:
		m.urlPatternBuilder = m.buildWithPath
	case _HOST:
		m.urlPatternBuilder = m.buildWithHost
	}
	return m
}

// checkUrlType check with section is the url start with
func checkUrlType(urlPattern string) (typ urlType) {
	first := urlPattern[0]
	switch {
	case first == '/':
		typ = _PATH
	case first == '?':
		typ = _QUERY
	case first == '#':
		typ = _FRAG
	default:
		typ = _HOST
	}
	return
}

// buildWithFrag use url's fragment to build pattern string
func (*matcher) buildWithFrag(url *url.URL) string {
	return url.Fragment
}

// buildWithQuery use url's querystring + fragment to build pattern string
func (*matcher) buildWithQuery(url_ *url.URL) string {
	buf := bytes.NewBuffer(make([]byte, 0, 20))
	buf.WriteString(buildUrlQuerystring(url_))
	if frag := url_.Fragment; frag != "" {
		buf.WriteByte('#')
		buf.WriteString(frag)
	}
	return buf.String()
}

// buildWithPath use url's path + querystring + fragment to build pattern string
func (*matcher) buildWithPath(url_ *url.URL) string {
	buf := bytes.NewBuffer(make([]byte, 0, 20))
	buf.WriteString(url_.Path)
	if query := buildUrlQuerystring(url_); query != "" {
		buf.WriteByte('?')
		buf.WriteString(query)
	}
	if frag := url_.Fragment; frag != "" {
		buf.WriteByte('#')
		buf.WriteString(frag)
	}
	return buf.String()
}

// buildWithHost use url's host+path+querystring+fragment to build pattern string
func (*matcher) buildWithHost(url_ *url.URL) string {
	buf := bytes.NewBuffer(make([]byte, 0, 20))
	buf.WriteString(url_.Host)
	buf.WriteString(url_.Path)
	if query := buildUrlQuerystring(url_); query != "" {
		buf.WriteByte('?')
		buf.WriteString(query)
	}
	if frag := url_.Fragment; frag != "" {
		buf.WriteByte('#')
		buf.WriteString(frag)
	}
	return buf.String()
}

// buildUrlQuerystring extract query string from url
func buildUrlQuerystring(url_ *url.URL) string {
	query, err := url.QueryUnescape(url_.RawQuery)
	if err != nil {
		query = ""
	}
	return query
}
