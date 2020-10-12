package json_filter

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	json "github.com/json-iterator/go"
)

type JSONFilter struct {
	reader    *bufio.Reader
	errWriter io.Writer
	Line      []byte
	fields    []string
	checker   BoolNoder
}

func (f *JSONFilter) Next() bool {
	var err error
	f.Line, err = f.reader.ReadBytes('\n')
	if err != nil {
		if !errors.Is(err, io.EOF) {
			fmt.Fprintf(f.errWriter, "read line error: %v\n", err)
		}
		return false
	}
	f.Line = bytes.TrimSpace(f.Line)
	ok, err := f.checker.Bool(f)
	if err != nil {
		fmt.Fprintf(f.errWriter, "check line error: %s\n", err.Error())
		return false
	}
	if !ok {
		return f.Next()
	}
	return true
}

func (f *JSONFilter) Get(key string) (interface{}, error) {
	return GetDataFromJSON(f.Line, key)
}

func (f *JSONFilter) GetData() ([]byte, error) {
	if len(f.fields) == 1 && f.fields[0] == "*" {
		return f.Line, nil
	}
	m := make(map[string]interface{})
	for _, field := range f.fields {
		if field == "*" {
			continue
		}
		data, err := f.Get(field)
		if err != nil {
			return nil, err
		}
		m[field] = data
	}
	bs, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func GetDataFromJSON(data []byte, key string) (interface{}, error) {
	if key == "[keys]" {
		m := make(map[string]interface{})
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, err
		}
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Sort(sort.StringSlice(keys))
		return strings.Join(keys, ","), nil
	}
	if !strings.Contains(key, ".") {
		return json.Get(data, key).GetInterface(), nil
	}
	keys := strings.Split(key, ".")
	j := json.Get(data, keys[0])
	for i := 1; i < len(keys); i++ {
		j = j.Get(keys[i])
	}
	return j.GetInterface(), nil
}

func GetFieldsAndChecker(sql string) ([]string, BoolNoder, error) {
	fields := make([]string, 0)
	tokens, err := Parse(sql)
	if err != nil {
		return nil, nil, fmt.Errorf("parse token error: %w", err)
	}
	if strings.ToLower(tokens[0].Str) != "select" {
		return nil, nil, fmt.Errorf("sql syntax error[1]")
	}
	foundExpression := false
	foundFrom := false
	expressionTokens := make([]*Token, 0)
	for i := 1; i < len(tokens); i++ {
		if strings.ToLower(tokens[i].Str) == "from" {
			foundFrom = true
			continue
		}
		if strings.ToLower(tokens[i].Str) == "t" {
			if foundFrom && strings.ToLower(tokens[i-1].Str) == "from" {
				if i-2 < 1 {
					return nil, nil, fmt.Errorf("sql syntax error[2]")
				}
				for j := 1; j < i-1; j++ {
					if tokens[j].Str == "," {
						continue
					}
					fields = append(fields, tokens[j].Str)
				}
				if len(fields) == 0 {
					return nil, nil, fmt.Errorf("sql syntax error[3]")
				}
				if i == len(tokens)-1 {
					foundExpression = false
					continue
				}
				if strings.ToLower(tokens[i+1].Str) != "where" {
					return nil, nil, fmt.Errorf("sql syntax error[4]")
				}
				expressionTokens = tokens[i+2:]
				if len(expressionTokens) == 0 {
					return nil, nil, fmt.Errorf("sql syntax error[5]")
				}
				foundExpression = true
			}
		}
	}
	if !foundExpression {
		return fields, NodeTrue{}, nil
	}
	node, err := parseTokens(expressionTokens)
	if err != nil {
		return nil, nil, fmt.Errorf("parse tokens to node error: %w", err)
	}
	bNode, ok := node.(BoolNoder)
	if !ok {
		return nil, nil, fmt.Errorf("node is not BoolNoder")
	}
	return fields, bNode, nil
}

func NewJSONFilterWithConfig(cfg FilterConfig) (*JSONFilter, error) {
	fields, checker, err := GetFieldsAndChecker(cfg.SQL)
	if err != nil {
		return nil, err
	}
	return &JSONFilter{
		reader:    bufio.NewReader(cfg.Reader),
		errWriter: cfg.ErrWriter,
		fields:    fields,
		checker:   checker,
	}, nil
}

type FilterConfig struct {
	Reader    io.Reader
	ErrWriter io.Writer
	SQL       string
}
