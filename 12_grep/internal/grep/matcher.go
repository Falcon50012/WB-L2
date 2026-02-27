package grep

import (
	"bytes"
	"regexp"
)

type Matcher interface {
	Match([]byte) bool
}

type FixedMatcher struct {
	pattern []byte
	lower   []byte
	ignore  bool
}

func NewFixedMatcher(cfg Config) Matcher {
	p := []byte(cfg.Pattern)

	if cfg.IgnoreCase {
		return &FixedMatcher{
			pattern: p,
			lower:   bytes.ToLower(p),
			ignore:  true,
		}
	}

	return &FixedMatcher{pattern: p}
}

func (m *FixedMatcher) Match(line []byte) bool {
	if m.ignore {
		return bytes.Contains(bytes.ToLower(line), m.lower)
	}
	return bytes.Contains(line, m.pattern)
}

type RegexMatcher struct {
	re *regexp.Regexp
}

func NewRegexMatcher(cfg Config) (Matcher, error) {
	p := cfg.Pattern
	if cfg.IgnoreCase {
		p = "(?i)" + p
	}
	re, err := regexp.Compile(p)
	if err != nil {
		return nil, err
	}
	return &RegexMatcher{re: re}, nil
}

func (m *RegexMatcher) Match(line []byte) bool {
	return m.re.Match(line)
}
