package unpack

import (
	"errors"
	"strings"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	var unpacked strings.Builder
	var prev rune
	var prevSet bool

	for pos := 0; pos < len(s); {
		r, size := utf8.DecodeRuneInString(s[pos:])
		pos += size

		if r == '\\' {
			if pos >= len(s) {
				return "", ErrInvalidString
			}
			esc, sz := utf8.DecodeRuneInString(s[pos:])
			pos += sz

			if prevSet {
				unpacked.WriteRune(prev)
			}
			prev = esc
			prevSet = true
			continue
		}

		if r >= '0' && r <= '9' {
			if !prevSet {
				return "", ErrInvalidString
			}
			num := int(r - '0')
			for pos < len(s) {
				rn, sz := utf8.DecodeRuneInString(s[pos:])
				if rn < '0' || rn > '9' {
					break
				}
				num = num*10 + int(rn-'0')
				pos += sz
			}

			if num > 0 {
				unpacked.WriteString(strings.Repeat(string(prev), num))
			}
			prevSet = false
			continue
		}

		if prevSet {
			unpacked.WriteRune(prev)
		}
		prev = r
		prevSet = true
	}

	if prevSet {
		unpacked.WriteRune(prev)
	}

	return unpacked.String(), nil
}
