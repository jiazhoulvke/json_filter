package json_filter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

var (
	_ BoolNoder = (*NodeAnd)(nil)
	_ BoolNoder = (*NodeOr)(nil)
	_ BoolNoder = (*NodeIn)(nil)
	_ BoolNoder = (*NodeNotIn)(nil)
	_ BoolNoder = (*NodeIsNull)(nil)
	_ BoolNoder = (*NodeIsNotNull)(nil)
	_ BoolNoder = (*NodeLike)(nil)
	_ BoolNoder = (*NodeNotLike)(nil)
	_ BoolNoder = (*NodeEqual)(nil)
	_ BoolNoder = (*NodeNotEqual)(nil)
	_ BoolNoder = (*NodeLessThan)(nil)
	_ BoolNoder = (*NodeLessEqual)(nil)
	_ BoolNoder = (*NodeGreaterThan)(nil)
	_ BoolNoder = (*NodeGreaterEqual)(nil)
	_ BoolNoder = (*NodeTrue)(nil)

	_ StringNoder = (*NodeString)(nil)

	_ FloatNoder = (*NodeField)(nil)
	_ FloatNoder = (*NodeNumber)(nil)
	_ FloatNoder = (*NodePlus)(nil)
	_ FloatNoder = (*NodeMinus)(nil)
	_ FloatNoder = (*NodeMult)(nil)
	_ FloatNoder = (*NodeDiv)(nil)
	_ FloatNoder = (*NodeMod)(nil)

	_ InterfaceNoder = (*NodeField)(nil)
	_ InterfaceNoder = (*NodeString)(nil)
	_ FloatNoder     = (*NodeNumber)(nil)
	_ FloatNoder     = (*NodePlus)(nil)
	_ FloatNoder     = (*NodeMinus)(nil)
	_ FloatNoder     = (*NodeMult)(nil)
	_ FloatNoder     = (*NodeDiv)(nil)
	_ FloatNoder     = (*NodeMod)(nil)
)

type NodeType int8

const (
	NodeTypeString NodeType = iota
	NodeTypeNumber
	NodeTypeField
	NodeTypeAnd
	NodeTypeOr
	NodeTypeIn
	NodeTypeNotIn
	NodeTypeIsNULL
	NodeTypeIsNotNULL
	NodeTypeLike
	NodeTypeNotLike
	NodeTypeEqual
	NodeTypeNotEqual
	NodeTypeLessThan
	NodeTypeLessEqual
	NodeTypeGreaterThan
	NodeTypeGreaterEqual
	NodeTypePlus
	NodeTypeMinus
	NodeTypeMult
	NodeTypeDiv
	NodeTypeMod
	NodeTypeTrue
)

type Noder interface {
	Type() NodeType
}

type InterfaceNoder interface {
	Noder
	Interface(Getter) (interface{}, error)
}

type StringNoder interface {
	Noder
	InterfaceNoder
	Str() string
}

type FloatNoder interface {
	Noder
	InterfaceNoder
	Float(Getter) (float64, error)
}

type BoolNoder interface {
	Noder
	Bool(Getter) (bool, error)
}

type NodeString struct {
	str string
}

func (n NodeString) Type() NodeType {
	return NodeTypeString
}

func (n NodeString) Str() string {
	return n.str
}

func (n NodeString) Interface(getter Getter) (interface{}, error) {
	return n.str, nil
}

type NodeNumber struct {
	f float64
}

func (n NodeNumber) Type() NodeType {
	return NodeTypeNumber
}

func (n NodeNumber) Float(getter Getter) (float64, error) {
	return n.f, nil
}

func (n NodeNumber) Interface(getter Getter) (interface{}, error) {
	return n.f, nil
}

type Getter interface {
	Get(string) (interface{}, error)
}

type NodeField struct {
	key string
}

func (n NodeField) Type() NodeType {
	return NodeTypeField
}

func (n NodeField) Float(getter Getter) (float64, error) {
	data, err := getter.Get(n.key)
	if err != nil {
		return 0, err
	}
	f, fOK := data.(float64)
	if fOK {
		return f, nil
	}
	s, sOK := data.(string)
	if !sOK {
		return 0, fmt.Errorf("unsupported data type")
	}
	return strconv.ParseFloat(s, 64)
}

func (n NodeField) Interface(getter Getter) (interface{}, error) {
	return getter.Get(n.key)
}

type NodeAnd struct {
	Left  BoolNoder
	Right BoolNoder
}

func (n NodeAnd) Type() NodeType {
	return NodeTypeAnd
}

func (n NodeAnd) Bool(getter Getter) (bool, error) {
	leftOK, leftErr := n.Left.Bool(getter)
	if leftErr != nil {
		return false, leftErr
	}
	if !leftOK {
		return false, nil
	}
	rightOK, rightErr := n.Right.Bool(getter)
	if rightErr != nil {
		return false, rightErr
	}
	if !rightOK {
		return false, nil
	}
	return true, nil
}

type NodeOr struct {
	Left  BoolNoder
	Right BoolNoder
}

func (n NodeOr) Type() NodeType {
	return NodeTypeOr
}

func (n NodeOr) Bool(getter Getter) (bool, error) {
	leftOK, leftErr := n.Left.Bool(getter)
	if leftErr != nil {
		return false, leftErr
	}
	if leftOK {
		return true, nil
	}
	rightOK, rightErr := n.Right.Bool(getter)
	if rightErr != nil {
		return false, rightErr
	}
	return rightOK, nil
}

type NodeIn struct {
	Key   string
	Slice []interface{}
}

func (n NodeIn) Type() NodeType {
	return NodeTypeIn
}

func (n NodeIn) Bool(getter Getter) (bool, error) {
	data, err := getter.Get(n.Key)
	if err != nil {
		return false, err
	}
	for _, item := range n.Slice {
		if item == data {
			return true, nil
		}
	}
	return false, nil
}

type NodeNotIn struct {
	Key   string
	Slice []interface{}
}

func (n NodeNotIn) Type() NodeType {
	return NodeTypeNotIn
}

func (n NodeNotIn) Bool(getter Getter) (bool, error) {
	data, err := getter.Get(n.Key)
	if err != nil {
		return false, err
	}
	for _, item := range n.Slice {
		if item == data {
			return false, nil
		}
	}
	return true, nil
}

type NodeIsNull struct {
	Key string
}

func (n NodeIsNull) Type() NodeType {
	return NodeTypeIsNULL
}

func (n NodeIsNull) Bool(getter Getter) (bool, error) {
	data, err := getter.Get(n.Key)
	if err != nil {
		return false, err
	}
	return data == nil, nil
}

type NodeIsNotNull struct {
	Key string
}

func (n NodeIsNotNull) Type() NodeType {
	return NodeTypeIsNotNULL
}

func (n NodeIsNotNull) Bool(getter Getter) (bool, error) {
	data, err := getter.Get(n.Key)
	if err != nil {
		return false, err
	}
	return data != nil, nil
}

type NodeLike struct {
	Key string
	Str string
}

func (n NodeLike) Type() NodeType {
	return NodeTypeLike
}

func (n NodeLike) Bool(getter Getter) (bool, error) {
	data, err := getter.Get(n.Key)
	if err != nil {
		return false, err
	}
	str, ok := data.(string)
	if !ok {
		return false, nil
	}
	return match(n.Str, str), nil
}

type NodeNotLike struct {
	Key string
	Str string
}

func (n NodeNotLike) Type() NodeType {
	return NodeTypeNotLike
}

func (n NodeNotLike) Bool(getter Getter) (bool, error) {
	data, err := getter.Get(n.Key)
	if err != nil {
		return false, err
	}
	str, ok := data.(string)
	if !ok {
		return false, nil
	}
	return !match(n.Str, str), nil
}

type NodeEqual struct {
	Left  InterfaceNoder
	Right InterfaceNoder
}

func (n NodeEqual) Type() NodeType {
	return NodeTypeEqual
}

func (n NodeEqual) Bool(getter Getter) (bool, error) {
	leftI, err := n.Left.Interface(getter)
	if err != nil {
		return false, err
	}
	rightI, err := n.Right.Interface(getter)
	if err != nil {
		return false, err
	}
	str1, ok1 := leftI.(string)
	str2, ok2 := rightI.(string)
	if ok1 != ok2 {
		return false, nil
	}
	if ok1 {
		return str1 == str2, nil
	}
	num1, ok1 := leftI.(float64)
	num2, ok2 := rightI.(float64)
	if ok1 != ok2 {
		return false, nil
	}
	return num1 == num2, nil
}

type NodeNotEqual struct {
	Left  InterfaceNoder
	Right InterfaceNoder
}

func (n NodeNotEqual) Type() NodeType {
	return NodeTypeNotEqual
}

func (n NodeNotEqual) Bool(getter Getter) (bool, error) {
	leftI, err := n.Left.Interface(getter)
	if err != nil {
		return false, err
	}
	rightI, err := n.Right.Interface(getter)
	if err != nil {
		return false, err
	}
	str1, ok1 := leftI.(string)
	str2, ok2 := rightI.(string)
	if ok1 != ok2 {
		return true, nil
	}
	if ok1 {
		return str1 != str2, nil
	}
	num1, ok1 := leftI.(float64)
	num2, ok2 := rightI.(float64)
	if ok1 != ok2 {
		return true, nil
	}
	return num1 != num2, nil
}

type NodeLessThan struct {
	Left  FloatNoder
	Right FloatNoder
}

func (n NodeLessThan) Type() NodeType {
	return NodeTypeLessThan
}

func (n NodeLessThan) Bool(getter Getter) (bool, error) {
	leftF, err := n.Left.Float(getter)
	if err != nil {
		return false, err
	}
	rightF, err := n.Right.Float(getter)
	if err != nil {
		return false, err
	}
	return leftF < rightF, nil
}

type NodeLessEqual struct {
	Left  FloatNoder
	Right FloatNoder
}

func (n NodeLessEqual) Type() NodeType {
	return NodeTypeLessEqual
}

func (n NodeLessEqual) Bool(getter Getter) (bool, error) {
	leftF, err := n.Left.Float(getter)
	if err != nil {
		return false, err
	}
	rightF, err := n.Right.Float(getter)
	if err != nil {
		return false, err
	}
	return leftF <= rightF, nil
}

type NodeGreaterThan struct {
	Left  FloatNoder
	Right FloatNoder
}

func (n NodeGreaterThan) Type() NodeType {
	return NodeTypeGreaterThan
}

func (n NodeGreaterThan) Bool(getter Getter) (bool, error) {
	leftF, err := n.Left.Float(getter)
	if err != nil {
		return false, err
	}
	rightF, err := n.Right.Float(getter)
	if err != nil {
		return false, err
	}
	return leftF > rightF, nil
}

type NodeGreaterEqual struct {
	Left  FloatNoder
	Right FloatNoder
}

func (n NodeGreaterEqual) Type() NodeType {
	return NodeTypeGreaterEqual
}

func (n NodeGreaterEqual) Bool(getter Getter) (bool, error) {
	leftF, err := n.Left.Float(getter)
	if err != nil {
		return false, err
	}
	rightF, err := n.Right.Float(getter)
	if err != nil {
		return false, err
	}
	return leftF >= rightF, nil
}

type NodePlus struct {
	Left  FloatNoder
	Right FloatNoder
}

func (n NodePlus) Type() NodeType {
	return NodeTypePlus
}

func (n NodePlus) Float(getter Getter) (float64, error) {
	leftFloat, err := n.Left.Float(getter)
	if err != nil {
		return 0, err
	}
	rightFloat, err := n.Right.Float(getter)
	if err != nil {
		return 0, err
	}
	num, _ := decimal.NewFromFloat(leftFloat).Add(decimal.NewFromFloat(rightFloat)).Float64()
	return num, nil
}

func (n NodePlus) Interface(getter Getter) (interface{}, error) {
	return n.Float(getter)
}

type NodeMinus struct {
	Left  FloatNoder
	Right FloatNoder
}

func (n NodeMinus) Type() NodeType {
	return NodeTypeMinus
}

func (n NodeMinus) Float(getter Getter) (float64, error) {
	leftFloat, err := n.Left.Float(getter)
	if err != nil {
		return 0, err
	}
	rightFloat, err := n.Right.Float(getter)
	if err != nil {
		return 0, err
	}
	num, _ := decimal.NewFromFloat(leftFloat).Sub(decimal.NewFromFloat(rightFloat)).Float64()
	return num, nil
}

func (n NodeMinus) Interface(getter Getter) (interface{}, error) {
	return n.Float(getter)
}

type NodeMult struct {
	Left  FloatNoder
	Right FloatNoder
}

func (n NodeMult) Type() NodeType {
	return NodeTypeMult
}

func (n NodeMult) Float(getter Getter) (float64, error) {
	leftFloat, err := n.Left.Float(getter)
	if err != nil {
		return 0, err
	}
	rightFloat, err := n.Right.Float(getter)
	if err != nil {
		return 0, err
	}
	num, _ := decimal.NewFromFloat(leftFloat).Mul(decimal.NewFromFloat(rightFloat)).Float64()
	return num, nil
}

func (n NodeMult) Interface(getter Getter) (interface{}, error) {
	return n.Float(getter)
}

type NodeDiv struct {
	Left  FloatNoder
	Right FloatNoder
}

func (n NodeDiv) Type() NodeType {
	return NodeTypeDiv
}

func (n NodeDiv) Float(getter Getter) (float64, error) {
	leftFloat, err := n.Left.Float(getter)
	if err != nil {
		return 0, err
	}
	rightFloat, err := n.Right.Float(getter)
	if err != nil {
		return 0, err
	}
	num, _ := decimal.NewFromFloat(leftFloat).Div(decimal.NewFromFloat(rightFloat)).Float64()
	return num, nil
}

func (n NodeDiv) Interface(getter Getter) (interface{}, error) {
	return n.Float(getter)
}

type NodeMod struct {
	Left  FloatNoder
	Right FloatNoder
}

func (n NodeMod) Type() NodeType {
	return NodeTypeMod
}

func (n NodeMod) Float(getter Getter) (float64, error) {
	leftFloat, err := n.Left.Float(getter)
	if err != nil {
		return 0, err
	}
	rightFloat, err := n.Right.Float(getter)
	if err != nil {
		return 0, err
	}
	num, _ := decimal.NewFromFloat(leftFloat).Mod(decimal.NewFromFloat(rightFloat)).Float64()
	return num, nil
}

func (n NodeMod) Interface(getter Getter) (interface{}, error) {
	return n.Float(getter)
}

type NodeTrue struct {
}

func (n NodeTrue) Type() NodeType {
	return NodeTypeTrue
}

func (n NodeTrue) Bool(getter Getter) (bool, error) {
	return true, nil
}

func parseTokens(tokens []*Token) (Noder, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("tokens is empty")
	}
	var parenDepth int
	if len(tokens) == 1 {
		switch tokens[0].Type {
		case TokenTypeString:
			return &NodeString{
				str: tokens[0].Str,
			}, nil
		case TokenTypeNumber:
			f, _ := strconv.ParseFloat(tokens[0].Str, 64)
			return &NodeNumber{
				f: f,
			}, nil
		default:
			return &NodeField{
				key: tokens[0].Str,
			}, nil
		}
	}
	pStartIndex := 0
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		if isLeftParen(t) {
			if parenDepth == 0 {
				pStartIndex = i
			}
			parenDepth++
		}
		if isRightParen(t) {
			parenDepth--
			if parenDepth == 0 {
				//如果整个列表的开头和结束都是括号，则去掉后重新解析
				if i == len(tokens)-1 && pStartIndex == 0 {
					return parseTokens(tokens[1:i])
				}
			}
		}
		if parenDepth == 0 && (isOr(t) || isAnd(t)) {
			left, err := parseTokens(tokens[0:i])
			if err != nil {
				return nil, err
			}
			leftB, leftOK := left.(BoolNoder)
			if !leftOK {
				return nil, fmt.Errorf("left is not BoolNoder")
			}
			right, err := parseTokens(tokens[i+1:])
			if err != nil {
				return nil, err
			}
			rightB, rightOK := right.(BoolNoder)
			if !rightOK {
				return nil, fmt.Errorf("right is not BoolNoder")
			}
			if isOr(t) {
				return NodeOr{
					Left:  leftB,
					Right: rightB,
				}, nil
			} else {
				return &NodeAnd{
					Left:  leftB,
					Right: rightB,
				}, nil
			}
		}
	}
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		if isLeftParen(t) {
			if parenDepth == 0 {
				pStartIndex = i
			}
			parenDepth++
		}
		if isRightParen(t) {
			parenDepth--
			if parenDepth == 0 {
				//如果整个列表的开头和结束都是括号，则去掉后重新解析
				if i == len(tokens)-1 && pStartIndex == 0 {
					return parseTokens(tokens[1:i])
				}
			}
		}
		if parenDepth == 0 && isLogicalOperator(t) {
			left, err := parseTokens(tokens[0:i])
			if err != nil {
				return nil, err
			}
			right, err := parseTokens(tokens[i+1:])
			if err != nil {
				return nil, err
			}
			if t.Str == OperatorEqual || t.Str == OperatorNotEqual || t.Str == OperatorNotEqual2 {
				leftI, leftOK := left.(InterfaceNoder)
				rightI, rightOK := right.(InterfaceNoder)
				if !leftOK {
					return nil, fmt.Errorf("left is not InterfaceNoder")
				}
				if !rightOK {
					return nil, fmt.Errorf("right is not InterfaceNoder")
				}
				switch t.Str {
				case OperatorEqual:
					return NodeEqual{
						Left:  leftI,
						Right: rightI,
					}, nil
				case OperatorNotEqual, OperatorNotEqual2:
					return NodeNotEqual{
						Left:  leftI,
						Right: rightI,
					}, nil
				}
			}
			leftF, leftOK := left.(FloatNoder)
			if !leftOK {
				return nil, fmt.Errorf("left is not FloatNoder")
			}
			rightF, rightOK := right.(FloatNoder)
			if !rightOK {
				return nil, fmt.Errorf("right is not FloatNoder")
			}
			switch t.Str {
			case OperatorLessThan:
				return NodeLessThan{
					Left:  leftF,
					Right: rightF,
				}, nil
			case OperatorLessEqual:
				return NodeLessEqual{
					Left:  leftF,
					Right: rightF,
				}, nil
			case OperatorGreaterThan:
				return NodeGreaterThan{
					Left:  leftF,
					Right: rightF,
				}, nil
			case OperatorGreaterEqual:
				return NodeGreaterEqual{
					Left:  leftF,
					Right: rightF,
				}, nil
			}
			return nil, fmt.Errorf("unknow node type")
		}
	}
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		if isArithmeticOperator(t) {
			left, err := parseTokens(tokens[0:i])
			if err != nil {
				return nil, err
			}
			right, err := parseTokens(tokens[i+1:])
			if err != nil {
				return nil, err
			}
			leftF, leftOK := left.(FloatNoder)
			if !leftOK {
				return nil, fmt.Errorf("left is not FloatNoder")
			}
			rightF, rightOK := right.(FloatNoder)
			if !rightOK {
				return nil, fmt.Errorf("right is not FloatNoder")
			}
			switch t.Str {
			case OperatorPlus:
				return &NodePlus{
					Left:  leftF,
					Right: rightF,
				}, nil
			case OperatorMinus:
				return &NodeMinus{
					Left:  leftF,
					Right: rightF,
				}, nil
			case OperatorMult:
				return &NodeMult{
					Left:  leftF,
					Right: rightF,
				}, nil
			case OperatorDiv:
				return &NodeDiv{
					Left:  leftF,
					Right: rightF,
				}, nil
			case OperatorMod:
				return &NodeMod{
					Left:  leftF,
					Right: rightF,
				}, nil
			}
			return nil, fmt.Errorf("unknow node type")
		}
	}
	//is null
	if len(tokens) == 3 && strings.ToLower(tokens[1].Str) == KeywordIs && strings.ToLower(tokens[2].Str) == KeywordNULL {
		return &NodeIsNull{
			Key: tokens[0].Str,
		}, nil
	}
	//is not null
	if len(tokens) == 4 && strings.ToLower(tokens[1].Str) == KeywordIs && strings.ToLower(tokens[2].Str) == KeywordNot && strings.ToLower(tokens[3].Str) == KeywordNULL {
		return &NodeIsNotNull{
			Key: tokens[0].Str,
		}, nil
	}
	//like
	if len(tokens) == 3 && strings.ToLower(tokens[1].Str) == KeywordLike && tokens[2].Type == TokenTypeString {
		return &NodeLike{
			Key: tokens[0].Str,
			Str: tokens[2].Str,
		}, nil
	}
	//not like
	if len(tokens) == 4 && strings.ToLower(tokens[1].Str) == KeywordNot && strings.ToLower(tokens[2].Str) == KeywordLike && tokens[3].Type == TokenTypeString {
		return &NodeNotLike{
			Key: tokens[0].Str,
			Str: tokens[3].Str,
		}, nil
	}
	// in
	if len(tokens) >= 5 && strings.ToLower(tokens[1].Str) == KeywordIn && isLeftParen(tokens[2]) && isRightParen(tokens[len(tokens)-1]) {
		data := make([]interface{}, 0)
		for j := 3; j < len(tokens)-1; j++ {
			data = append(data, tokens[j].Str)
		}
		return &NodeIn{
			Key:   tokens[0].Str,
			Slice: data,
		}, nil
	}
	// not in
	if len(tokens) >= 6 && strings.ToLower(tokens[1].Str) == KeywordNot && strings.ToLower(tokens[2].Str) == KeywordIn && isLeftParen(tokens[3]) && isRightParen(tokens[len(tokens)-1]) {
		data := make([]interface{}, 0)
		for j := 4; j < len(tokens)-1; j++ {
			data = append(data, tokens[j].Str)
		}
		return &NodeNotIn{
			Key:   tokens[0].Str,
			Slice: data,
		}, nil
	}
	return nil, nil
}

func isAnd(t *Token) bool {
	return t.Type == TokenTypeUnknow && strings.ToLower(t.Str) == "and"
}

func isOr(t *Token) bool {
	return t.Type == TokenTypeUnknow && strings.ToLower(t.Str) == "or"
}

func isArithmeticOperator(t *Token) bool {
	if t.Type != TokenTypeOperator {
		return false
	}
	switch t.Str {
	case OperatorPlus, OperatorMinus, OperatorMult, OperatorDiv, OperatorMod:
		return true
	}
	return false
}

func isLogicalOperator(t *Token) bool {
	if t.Type != TokenTypeOperator {
		return false
	}
	switch t.Str {
	case "=", "<>", "!=", "<", ">", "<=", ">=":
		return true
	}
	return false
}
