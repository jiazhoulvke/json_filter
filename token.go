package json_filter

import (
	"fmt"
	"strings"
)

type Token struct {
	Str  string
	Type TokenType
}

func (t Token) String() string {
	s := t.Str
	if t.Type == TokenTypeString {
		s = "'" + s + "'"
	}
	return fmt.Sprintf("%s", s)
}

type TokenType int8

const (
	TokenTypeUnknow TokenType = iota
	TokenTypeString
	TokenTypeNumber
	TokenTypeOperator
	TokenTypeKeyword
	TokenTypeLeftParen
	TokenTypeRightParen
)

func (nt TokenType) String() string {
	switch nt {
	case TokenTypeUnknow:
		return "unknow"
	case TokenTypeString:
		return "string"
	case TokenTypeNumber:
		return "number"
	case TokenTypeOperator:
		return "operator"
	case TokenTypeKeyword:
		return "keyword"
	case TokenTypeLeftParen:
		return "("
	case TokenTypeRightParen:
		return ")"
	default:
		panic("unsupported type")
	}
}

type ParseTokenState int8

const (
	ParseStateNormal ParseTokenState = iota
	ParseStateInOperator
)

//Parse 将sql解析成token
func Parse(sql string) ([]*Token, error) {
	str := sql
	tokens := make([]*Token, 0)
	var state ParseTokenState = ParseStateNormal
	var stringLeftChar byte
	var i int
	for i = 0; i < len(str); {
		//默认状态
		if state == ParseStateNormal {
			if isWhitespace(str[i]) {
				if i > 0 {
					t := &Token{
						Str:  str[0:i],
						Type: TokenTypeUnknow,
					}
					if isNumber(t.Str) {
						t.Type = TokenTypeNumber
					}
					tokens = append(tokens, t)
				}
				str = str[i+1:]
				i = 0
				continue
			}
			switch str[i] {
			case '\'', '"':
				stringLeftChar = str[i]
				if i > 0 {
					t := &Token{
						Str:  str[0:i],
						Type: TokenTypeUnknow,
					}
					tokens = append(tokens, t)
				}
				str = str[i+1:]
				i = 0
				stringRightCharIndex := strings.Index(str, string(stringLeftChar))
				if stringRightCharIndex == -1 {
					return nil, fmt.Errorf("sql error: %s", sql)
				}
				t := &Token{
					Str:  str[0:stringRightCharIndex],
					Type: TokenTypeString,
				}
				tokens = append(tokens, t)
				str = str[stringRightCharIndex+1:]
				i = 0
			case '(', ')':
				if i > 0 {
					t := &Token{
						Str:  str[0:i],
						Type: TokenTypeUnknow,
					}
					if isNumber(t.Str) {
						t.Type = TokenTypeNumber
					}
					tokens = append(tokens, t)
				}
				t := &Token{
					Type: TokenTypeOperator,
					Str:  string(str[i]),
				}
				if str[i] == '(' {
					t.Type = TokenTypeLeftParen
				} else {
					t.Type = TokenTypeRightParen
				}
				tokens = append(tokens, t)
				str = str[i+1:]
				i = 0
			case ',':
				if i > 0 {
					t := &Token{
						Str:  str[0:i],
						Type: TokenTypeUnknow,
					}
					if isNumber(t.Str) {
						t.Type = TokenTypeNumber
					}
					tokens = append(tokens, t)
				}
				tokens = append(tokens, &Token{
					Str:  ",",
					Type: TokenTypeKeyword,
				})
				str = str[i+1:]
				i = 0
			case '+', '-', '*', '/', '%', '=',
				'>', '!', '<': //'<=','>=','!=','<>'
				if i > 0 {
					tokens = append(tokens, &Token{
						Str:  str[0:i],
						Type: TokenTypeUnknow,
					})
				}
				state = ParseStateInOperator
				i++
			default:
				if i == len(str)-1 {
					t := Token{
						Str:  str[0:],
						Type: TokenTypeUnknow,
					}
					if isNumber(t.Str) {
						t.Type = TokenTypeNumber
					}
					tokens = append(tokens, &t)
				}
				i++
			}
		} else if state == ParseStateInOperator { // 操作符中
			if isWhitespace(str[i]) && str[i-1] == '*' {
				tokens = append(tokens, &Token{
					Str:  string('*'),
					Type: TokenTypeUnknow,
				})
				str = str[i:]
			} else if str[i] == '=' || str[i] == '>' {
				tokens = append(tokens, &Token{
					Str:  str[i-1 : i+1],
					Type: TokenTypeOperator,
				})
				str = str[i+1:]
			} else {
				tokens = append(tokens, &Token{
					Str:  string(str[i-1]),
					Type: TokenTypeOperator,
				})
				str = str[i:]
			}
			state = ParseStateNormal
			i = 0
		} else { //未知状态
			return nil, fmt.Errorf("unknow error: sql:%s, char:%c", str, str[0])
		}
	}
	return tokens, nil
}

func isLeftParen(t *Token) bool {
	return t.Type == TokenTypeLeftParen
}

func isRightParen(t *Token) bool {
	return t.Type == TokenTypeRightParen
}
