package json_filter

import (
	"strconv"
)

func isWhitespace(c byte) bool {
	return c == ' ' ||
		c == '\t' ||
		c == '\n' ||
		c == '\v' ||
		c == '\f' ||
		c == '\r'
}

func isDigitNum(c byte) bool {
	return '0' <= c && c <= '9' || c == '.'
}

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return true
	}
	return false
}

func isKeyword(s string) bool {
	_, ok := keywords[s]
	return ok
}
