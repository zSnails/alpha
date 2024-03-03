package tokenizer

import (
	"errors"
	"fmt"
	"io"
	"regexp"
)

type TokenType int8

type Token struct {
	Type     TokenType `json:"type"`
	Value    string    `json:"value"`
	row, col int
}

func (t *Token) GetPosition() (int, int) {
	return t.row, t.col
}

func (t *Token) String() string {
	// TODO: report line and characer instead of this piece of shit.
	// return fmt.Sprintf("[<%s>@%d:%d %s]", TokenNames[t.Type], t.row, t.col, t.Value)
	return fmt.Sprintf("<%s>", TokenNames[t.Type])
}

const (
	Whitespace TokenType = iota
	If
	Then
	Else
	While
	Do
	Let
	Var
	Const
	Tilde
	In
	Begin
	End
	Identifier
	Float
	Integer
	PlusOperator
	MinusOperator
	DivisionOperator
	MultiplicationOperator
	Equals
	Comparison
	LessThan
	GreaterThan
	LessThanEqual
	GreaterThanEqual
	LeftParenthesis
	RightParenthesis
	Colon
	Semicolon
	String
)

type Spec struct {
	Type TokenType `json:"type"`
	Spec string    `json:"expression"`
}

var SPECS = []Spec{
	{
		Type: Whitespace,
		Spec: `^[ \r\n\t]+`,
	},
	{
		Type: Whitespace,
		Spec: `^\/\/.*`,
	},
	{
		Type: String,
		Spec: `^"[^"]*`,
	},
	{
		Type: String,
		Spec: `^'[^']*`,
	},
	{
		Type: If,
		Spec: `^if`,
	},
	{
		Type: End,
		Spec: `^end`,
	},
	{
		Type: Tilde,
		Spec: `^\~`,
	},
	{
		Type: Colon,
		Spec: `^\:`,
	},
	{
		Type: Semicolon,
		Spec: `^\;`,
	},
	{
		Type: Then,
		Spec: `^then`,
	},
	{
		Type: Else,
		Spec: `^else`,
	},

	// NOTE: The order is important
	{
		Type: Float,
		Spec: `^\d+\.\d+`,
	},
	{
		Type: Integer,
		Spec: `^\d+`,
	},
	{
		Type: PlusOperator,
		Spec: `^\+`,
	},
	{
		Type: MinusOperator,
		Spec: `^\-`,
	},
	{
		Type: DivisionOperator,
		Spec: `^\/`,
	},
	{
		Type: MultiplicationOperator,
		Spec: `^\*`,
	},
	{
		Type: LeftParenthesis,
		Spec: `^\(`,
	},
	{
		Type: RightParenthesis,
		Spec: `^\)`,
	},
	{
		Type: Comparison,
		Spec: `^==`,
	},
	{
		Type: Equals,
		Spec: `^=`,
	},
	{
		Type: LessThan,
		Spec: `^<`,
	},
	{
		Type: GreaterThan,
		Spec: `^>`,
	},
	{
		Type: LessThanEqual,
		Spec: `^<=`,
	},
	{
		Type: GreaterThanEqual,
		Spec: `^>=`,
	},
	{
		Type: While,
		Spec: `^while`,
	},
	{
		Type: Do,
		Spec: `^do`,
	},
	{
		Type: Let,
		Spec: `^let`,
	},
	{
		Type: Var,
		Spec: `^var`,
	},
	{
		Type: Const,
		Spec: `^const`,
	},
	{
		Type: In,
		Spec: `^in`,
	},
	{
		Type: Begin,
		Spec: `^begin`,
	},
	{
		Type: Identifier,
		Spec: `^[_a-zA-Z][_a-zA-Z0-9]*`,
	},
}

var TokenNames = map[TokenType]string{
	Whitespace:             "whitespace",
	If:                     "if",
	Then:                   "then",
	Else:                   "else",
	End:                    "end",
	Identifier:             "identifier",
	Float:                  "float",
	Integer:                "integer",
	PlusOperator:           "+",
	MinusOperator:          "-",
	DivisionOperator:       "/",
	MultiplicationOperator: "*",
	LeftParenthesis:        "(",
	RightParenthesis:       ")",
	Comparison:             "comparison",
	Equals:                 "=",
	LessThan:               "<",
	GreaterThan:            ">",
	LessThanEqual:          "<=",
	GreaterThanEqual:       ">=",
	While:                  "while",
	Do:                     "do",
	Let:                    "let",
	Var:                    "var",
	Const:                  "const",
	In:                     "in",
	Begin:                  "begin",
	String:                 "string",
	Tilde:                  "~",
	Colon:                  ":",
	Semicolon:              ";",
}

type Tokenizer struct {
	content string
	cursor  int
}

func NewTokenizer(content string) *Tokenizer {
	return &Tokenizer{
		content: content,
		cursor:  0,
	}
}

func (t *Tokenizer) match(spec, content string) (string, int) {
	re := regexp.MustCompile(spec)
	matched := re.Find([]byte(content))

	// TODO: check if this actually does something
	// if matched == nil || len(matched) == 0 {
	// 	return "", 0
	// }

	size := len(matched)
	t.cursor += size
	return string(matched), size
}

func (t *Tokenizer) hasMoreTokens() bool {
	return t.cursor < len(t.content)
}

func (t *Tokenizer) GetAllTokens() ([]*Token, error) {
	out := []*Token{}
	for {
		tok, err := t.GetNextToken()
		if err != nil && errors.Is(err, io.EOF) {
			return out, nil
		} else if err != nil {
			return nil, err
		}
		out = append(out, tok)
	}
}

func (t *Tokenizer) GetNextToken() (*Token, error) {
	if !t.hasMoreTokens() {
		return nil, io.EOF
	}

	for _, spec := range SPECS {
		matched, size := t.match(spec.Spec, t.content[t.cursor:])
		if size == 0 {
			continue
		}

		if spec.Type == Whitespace {
			return t.GetNextToken()
		}

		if spec.Type == String {
			t.cursor++ // Skip the closing quote on strings
		}

		return &Token{
			Type:  spec.Type,
			Value: matched,
		}, nil
	}

	return nil, fmt.Errorf("unexpected token %c at position %d\n", t.content[t.cursor], t.cursor)
}
