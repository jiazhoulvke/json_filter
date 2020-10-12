package json_filter

var keywords = map[string]string{
	KeywordIn:   KeywordIn,
	KeywordNot:  KeywordNot,
	KeywordLike: KeywordLike,
	KeywordAnd:  KeywordAnd,
	KeywordOr:   KeywordOr,
	KeywordIs:   KeywordIs,
	KeywordNULL: KeywordNULL,
}

const (
	KeywordIn   = "in"
	KeywordNot  = "not"
	KeywordLike = "like"
	KeywordAnd  = "and"
	KeywordOr   = "or"
	KeywordIs   = "is"
	KeywordNULL = "null"
)

const (
	OperatorPlus         = "+"
	OperatorMinus        = "-"
	OperatorMult         = "*"
	OperatorDiv          = "/"
	OperatorMod          = "%"
	OperatorEqual        = "="
	OperatorLessThan     = "<"
	OperatorGreaterThan  = ">"
	OperatorLessEqual    = "<="
	OperatorGreaterEqual = ">="
	OperatorNotEqual     = "!="
	OperatorNotEqual2    = "<>"
)
