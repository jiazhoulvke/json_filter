package json_filter

import "unicode/utf8"

//match https://github.com/golang/go/blob/master/src/path/filepath/match.go
func match(pattern string, str string) bool {
Pattern:
	for len(pattern) > 0 {
		var percentSymbol bool
		var chunk string
		percentSymbol, chunk, pattern = scanChunk(pattern)
		if percentSymbol && chunk == "" {
			return true
		}
		t, ok := matchChunk(chunk, str)
		if ok && (len(t) == 0 || len(pattern) > 0) {
			str = t
			continue
		}
		if percentSymbol {
			for i := 0; i < len(str); i++ {
				t, ok := matchChunk(chunk, str[i+1:])
				if ok {
					if len(pattern) == 0 && len(t) > 0 {
						continue
					}
					str = t
					continue Pattern
				}
			}
		}
		return false
	}
	return len(str) == 0
}

func scanChunk(pattern string) (percentSymbol bool, chunk, rest string) {
	for len(pattern) > 0 && pattern[0] == '%' {
		pattern = pattern[1:]
		percentSymbol = true
	}
	var i int
Scan:
	for i = 0; i < len(pattern); i++ {
		if pattern[i] == '%' {
			break Scan
		}
	}
	return percentSymbol, pattern[0:i], pattern[i:]
}

func matchChunk(chunk, s string) (rest string, ok bool) {
	for len(chunk) > 0 {
		if len(s) == 0 {
			return
		}
		switch chunk[0] {
		case '_':
			_, n := utf8.DecodeLastRuneInString(s)
			s = s[n:]
			chunk = chunk[1:]
		default:
			if chunk[0] != s[0] {
				return
			}
			s = s[1:]
			chunk = chunk[1:]
		}
	}
	return s, true
}
