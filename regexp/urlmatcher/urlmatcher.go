// Package urlmatcher match url like "{group:regexp}", extract values
package urlmatcher

import (
	"bytes"
	"strings"

	"github.com/cosiner/golib/regexp"

	. "github.com/cosiner/golib/errors"
)

const (
	// url pattern type
	_LITERAL        = iota // _LITERAL means a normal string without regexp
	_REGEXP_LITERAL        //_REGXP_LITERAL means a regexp string without groups
	_REGEXP                // _REGEXP means a regexp with groups

	_START   // _START means ready to match '{''
	_INPARSE // _INPARSE means already matched '{', ready to match '}'
	_INGROUP // _INGROUP means match group name
	_INREGEX // _INREGEX means match regexp string

	_DEF_REGEXP = "[^/]*" // _DEF_REGEXP is  the default regexp for url section {group}
)

// Matcher is an url matcher that match
// any regular expression must quote by "{}", otherwise, it will be treated
// as literal, regardless you may not need it's value
type Matcher struct {
	isLiteral      bool           // isLiteral means whether url pattern is just literal
	pattern        string         // pattern is url pattern after parse
	literalPattern []string       // literal is url literal split by '/'
	regexpPattern  *regexp.Regexp // pattern is regexp pattern for non-literal url
	// Match match given url, reutrn matched values and if it is match or not
	// if matcher is a literal matcher, no matched value will be returned
	Match func(url string) (map[string]string, bool)
	// MatchOnly only return whether url is match, don't extract url stories
	MatchOnly func(url string) bool
}

// Pattern return url pattern after parse
func (m *Matcher) Pattern() string {
	return m.pattern
}

// IsLiteral check whether it's a literal matcher
func (m *Matcher) IsLiteral() bool {
	return m.isLiteral
}

// literalMatch match literal url
func (m *Matcher) literalMatch(url string) (vals map[string]string, match bool) {
	return regexp.NIL_MAP, m.literalMatchOnly(url)
}

// literalMatchOnly only match literal url, don't extract url variables
func (m *Matcher) literalMatchOnly(url string) bool {
	return literalMatch(m.literalPattern, strings.Split(url, "/"))
}

// literalMatch match two url section, return true only when literal pattern is
// prefix of url pattern
func literalMatch(literalPattern []string, urlPattern []string) (match bool) {
	if match = (len(literalPattern) <= len(urlPattern)); match {
		for i, s := range literalPattern {
			if s != urlPattern[i] {
				match = false
				break
			}
		}
	}
	return
}

// Compile compile string like "{groupname:regexp} {groupname} {:regexp}" to regexp
// for {groupname:regexp}, matched regexp values will be
// stored in map with key groupname
// for {:regexp}, match regexp but don't get matched value
// for {groupname}, use default regexp `[^/]*`, {groupname:} is also treated as {groupname}
// other format like {}, {:}... is all wrong format
func Compile(urlPattern string) (matcher *Matcher, err error) {
	if urlPattern == "" {
		return nil, Err("No content")
	}
	var (
		buf           = bytes.NewBuffer(make([]byte, 0, len(urlPattern)+10))
		state, rstate = _START, _INGROUP
		group, regexp []byte
		// isLiteral
		urlType int = _LITERAL
	)
	for _, c := range urlPattern {
		switch {
		case c == '{':
			if state == _INPARSE {
				break
			}
			urlType, state, group, regexp =
				_REGEXP_LITERAL, _INPARSE, group[:0], regexp[:0]
		case c == '}':
			if state == _START {
				state = _INPARSE
				break
			}
			success, hasGroup := writeRegexp(buf, group, regexp)
			if !success {
				break
			} else if hasGroup {
				urlType = _REGEXP
			}
			state, rstate = _START, _INGROUP
		case state == _START:
			buf.WriteRune(c)
		case c == ':':
			rstate = _INREGEX
		case rstate == _INGROUP:
			group = append(group, byte(c))
		default:
			regexp = append(regexp, byte(c))
		}
	}
	if state == _INPARSE {
		err = Errorf("Wrong format:%s", urlPattern)
	} else {
		matcher, err = newMatcher(urlType, urlPattern, buf.String())
	}
	return
}

// newMatcher create a new matcher, if isLiteral, retuened a literal matcher
// else return a regexp matcher
func newMatcher(urlType int, literalPattern, regexpPattern string) (matcher *Matcher, err error) {
	matcher = &Matcher{isLiteral: urlType == _LITERAL}
	if urlType == _LITERAL {
		matcher.pattern = literalPattern
		matcher.literalPattern = strings.Split(literalPattern, "/")
		matcher.Match = matcher.literalMatch
		matcher.MatchOnly = matcher.literalMatchOnly
	} else if matcher.regexpPattern, err = regexp.Compile(regexpPattern); err == nil {
		matcher.pattern = regexpPattern
		if urlType == _REGEXP {
			matcher.Match = matcher.regexpPattern.SingleSubmatchMap
		} else {
			matcher.Match = func(url string) (map[string]string, bool) {
				return regexp.NIL_MAP, matcher.MatchOnly(url)
			}
		}
		matcher.MatchOnly = matcher.regexpPattern.MatchString
	}
	return
}

// writeRegexp write group and regexp as a grouped regexp to buffer
// if group and regexp is both empty, write failed
func writeRegexp(buf *bytes.Buffer, group, regexp []byte) (success, hasGroup bool) {
	groupLen, regexpLen := len(group), len(regexp)
	if success = (groupLen != 0 || regexpLen != 0); success {
		if groupLen == 0 {
			buf.Write(regexp)
		} else {
			hasGroup = true
			buf.WriteString("(?P<")
			buf.Write(group)
			buf.WriteByte('>')
			if regexpLen == 0 {
				buf.WriteString(_DEF_REGEXP)
			} else {
				buf.Write(regexp)
			}
			buf.WriteByte(')')
		}
	}
	return
}
