package tokenizer

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
)

type TokenType int8

type Token struct {
	Type     TokenType `json:"type"`
	Value    string    `json:"value"`
	row, col int
}

func (t *Token) GetPosition() (row int, col int) {
	return t.row, t.col
}

func (t *Token) String() string {
	return fmt.Sprintf("[<%s>@%d:%d %s]", TokenNames[t.Type], t.row, t.col, t.Value)
}

const (
	EOF TokenType = iota
	Whitespace
	NewLine
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
		Type: NewLine,
		Spec: `^[\r\n]`,
	},
	{
		Type: Whitespace,
		Spec: `^[ \t]+`,
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
		Spec: `^if\b`,
	},
	{
		Type: End,
		Spec: `^end\b`,
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
		Spec: `^then\b`,
	},
	{
		Type: Else,
		Spec: `^else\b`,
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
		Spec: `^while\b`,
	},
	{
		Type: Do,
		Spec: `^do\b`,
	},
	{
		Type: Let,
		Spec: `^let\b`,
	},
	{
		Type: Var,
		Spec: `^var\b`,
	},
	{
		Type: Const,
		Spec: `^const\b`,
	},
	{
		Type: In,
		Spec: `^in\b`,
	},
	{
		Type: Begin,
		Spec: `^begin\b`,
	},
	{
		Type: Identifier,
		Spec: `^[_a-zA-Z][_a-zA-Z0-9]*`,
	},
}

var TokenNames = map[TokenType]string{
	EOF:                    "EOF",
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
	file    string
	line    int
	char    int
}

func (t *Tokenizer) GetFileName() string {
	return t.file
}

func FromFile(name string) (*Tokenizer, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	return &Tokenizer{
		content: string(data),
		cursor:  0,
		char:    1,
		line:    1,
		file:    path.Base(name),
	}, nil
}

func NewTokenizer(content string) *Tokenizer {
	return &Tokenizer{
		content: content,
		cursor:  0,
		char:    1,
		line:    1,
		file:    "<stdin>",
	}
}

func (t *Tokenizer) match(spec, content string) (string, int) {
	re := regexp.MustCompile(spec)
	matched := re.Find([]byte(content))

	size := len(matched)
	t.cursor += size
	t.char += size
	return string(matched), size
}

func (t *Tokenizer) hasMoreTokens() bool {
	return t.cursor < len(t.content)
}

func (t *Tokenizer) GetAllTokens() ([]*Token, error) {
	out := []*Token{}
	for {
		tok, err := t.GetNextToken()
		if err != nil {
			if errors.Is(err, io.EOF) {
				out = append(out, tok)
				return out, nil
			}
			return nil, err
		}
		out = append(out, tok)
	}
}

// GetNextToken returns the next token recognized from the input stream
//
// returns io.EOF when the last token is reached, when the error is `io.EOF` the last token returned is
// an EOF token indicating that an end of file has been found
func (t *Tokenizer) GetNextToken() (*Token, error) {
	if !t.hasMoreTokens() {
		return &Token{
			Type: EOF,
			row:  t.line,
			col:  t.char,
		}, io.EOF
	}

	for _, spec := range SPECS {
		matched, size := t.match(spec.Spec, t.content[t.cursor:])
		if size == 0 {
			continue
		}

		if spec.Type == NewLine {
			t.line++
			t.char = 1
			return t.GetNextToken()
		}

		if spec.Type == Whitespace {
			return t.GetNextToken()
		}

		if spec.Type == String {
			if !t.hasMoreTokens() || (t.content[t.cursor] != '"' && t.content[t.cursor] != '\'') {
				return nil, fmt.Errorf("%s:%d:%d: syntax error: missing string closing quote", t.file, t.line, t.char)
			}
			t.cursor++ // Skip the closing quote on strings
		}

		return &Token{
			Type:  spec.Type,
			Value: matched,
			col:   t.char - size,
			row:   t.line,
		}, nil
	}

	return nil, fmt.Errorf("%s:%d:%d: syntax error: unexpected token '%c'\n", t.file, t.line, t.char, t.content[t.cursor])
}
