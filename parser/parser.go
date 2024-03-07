package parser

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/zSnails/alpha/parser/ast"
	"github.com/zSnails/alpha/tokenizer"
)

// The parser structure implements a recursive descent parser
type Parser struct {
	currentToken int
	tokens       []*tokenizer.Token
	lexer        *tokenizer.Tokenizer
}

// getCurrentoken returns the current token to be worked on
func (p *Parser) getCurrentToken() (*tokenizer.Token, error) {
	if !p.tokensLeft() {
		return nil, io.EOF
	}
	return p.tokens[p.currentToken], nil
}

// mustGetCurrentToken always returns the current token and doesn't do any boundary checks
// useful for cases where you know there must be an available token to be consumed.
func (p *Parser) mustGetCurrentToken() *tokenizer.Token {
	return p.tokens[p.currentToken]
}

// NewParser returns an instance of a brand new parser consuming the tokens in
// the lexer.
func NewParser(lexer *tokenizer.Tokenizer) (*Parser, error) {
	tokens, err := lexer.GetAllTokens()
	if err != nil {
		return nil, err
	}

	return &Parser{
		currentToken: 0,
		lexer:        lexer,
		tokens:       tokens,
	}, nil
}

func (p *Parser) expect(_type tokenizer.TokenType) error {
	token, err := p.getCurrentToken()
	if err != nil || token.Type != _type {
		return p.UnexpectedToken(token, _type)
	}
	p.advance()
	return nil
}

func (p *Parser) advance() {
	p.currentToken++
}

// Program parses the basic program construct
//
//	program ::= singleCommand
func (p *Parser) Program() (ast.Node, error) {
	node, err := p.SingleCommand()
	if err != nil {
		return nil, err
	}

	err = p.expect(tokenizer.EOF)
	if err != nil {
		return nil, err
	}
	return node, nil
}

// SingleCommand parses the basic singleCommand construct
//
//	singleCommand ::=
//	         Identifier (= expression | (expression))
//	        | if expression then singleCommand
//	        | while expression do singleCommand
//	        | let declaration in singleCommand
//	        | begin command end
func (p *Parser) SingleCommand() (ast.Node, error) {
	currentToken, err := p.getCurrentToken() // this error will always be io.EOF
	if err != nil {
		return nil, err
	}

	switch currentToken.Type {
	case tokenizer.Identifier:
		{
			ident := p.mustGetCurrentToken()
			p.advance()
			next, err := p.getCurrentToken()
			if err != nil {
				return nil, p.UnexpectedToken(next, tokenizer.Equals, tokenizer.LeftParenthesis)
			}
			switch next.Type {
			case tokenizer.Equals:
				{
					p.advance()
					expressionNode, err := p.Expression()
					if err != nil {
						return nil, err
					}
					return ast.NewAssignment(ident, expressionNode), nil
				}
			case tokenizer.LeftParenthesis:
				{
					p.advance()
					next, err = p.getCurrentToken()
					if err != nil {
						return nil, p.UnexpectedToken(next,
							tokenizer.RightParenthesis, tokenizer.Integer,
							tokenizer.Float, tokenizer.String)
					}
					if next.Type == tokenizer.RightParenthesis {
						p.advance()
						return ast.NewFunctionCall(ident, nil), nil
					}

					expressionNode, err := p.Expression()
					if err != nil {
						return nil, err
					}
					err = p.expect(tokenizer.RightParenthesis)
					if err != nil {
						return nil, err
					}
					return ast.NewFunctionCall(ident, expressionNode), nil
				}
			}
		}
	case tokenizer.If:
		{
			p.advance()
			expressionNode, err := p.Expression()
			if err != nil {
				return nil, err
			}
			err = p.expect(tokenizer.Then)
			if err != nil {
				return nil, err
			}
			ifBlockSingleCommand, err := p.SingleCommand()
			if err != nil {
				return nil, err
			}
			err = p.expect(tokenizer.Else)
			if err != nil {
				return nil, err
			}
			elseBlockSingleCommand, err := p.SingleCommand()
			if err != nil {
				return nil, err
			}

			return ast.NewIfBlock(expressionNode, ifBlockSingleCommand, elseBlockSingleCommand), nil
		}
	case tokenizer.While:
		{
			p.advance()
			expression, err := p.Expression()
			if err != nil {
				return nil, err
			}
			err = p.expect(tokenizer.Do)
			if err != nil {
				return nil, err
			}
			singleCommand, err := p.SingleCommand()
			if err != nil {
				return nil, err
			}
			return ast.NewWhileBlock(expression, singleCommand), nil
		}

	case tokenizer.Let:
		{
			p.advance()
			declaration, err := p.Declaration()
			if err != nil {
				return nil, err
			}

			err = p.expect(tokenizer.In)
			if err != nil {
				return nil, err
			}

			singleCommand, err := p.SingleCommand()
			if err != nil {
				return nil, err
			}

			return ast.NewLetBlock(declaration, singleCommand), nil
		}
	case tokenizer.Begin:
		{
			p.advance()
			command, err := p.Command()
			if err != nil {
				return nil, err
			}
			err = p.expect(tokenizer.End)
			if err != nil {
				return nil, err
			}
			return ast.NewBeginBlock(command), nil
		}
	}
	return nil, p.UnexpectedToken(currentToken, tokenizer.Begin, tokenizer.Let, tokenizer.While, tokenizer.If, tokenizer.Identifier)
}

func Map[T, B any](slice []T, f func(T) B) []B {
	out := []B{}
	for _, v := range slice {
		out = append(out, f(v))
	}
	return out
}

var EOF = &tokenizer.Token{
	Type: tokenizer.EOF,
}

func (p *Parser) UnexpectedToken(got *tokenizer.Token, expected ...tokenizer.TokenType) error {
	row, col := got.GetPosition()
	if len(expected) > 1 {
		tokens := Map(expected, func(token tokenizer.TokenType) string {
			return fmt.Sprintf("'%s'", tokenizer.TokenNames[token])
		})
		expectedTokens := strings.Join(tokens, ", ")

		return fmt.Errorf("%s:%d:%d: unexpected token '%s' expected one of %s", p.lexer.GetFileName(), row, col, tokenizer.TokenNames[got.Type], expectedTokens)
	} else if len(expected) == 1 {
		return fmt.Errorf("%s:%d:%d: unexpected token '%s' expected '%s'\n", p.lexer.GetFileName(), row, col, tokenizer.TokenNames[got.Type], tokenizer.TokenNames[expected[0]])
	}
	return fmt.Errorf("%s:%d:%d: unexpected token '%s'\n", p.lexer.GetFileName(), row, col, tokenizer.TokenNames[got.Type])
}

// Declaration parses the basic declaration construct
//
// declaration ::= singleDeclaration (; singleDeclaration)*
func (p *Parser) Declaration() (*ast.Declaration, error) {
	node := ast.NewDeclaration()
	singleDeclaration, err := p.SingleDeclaration()
	if err != nil {
		return nil, err
	}
	node.AddChild(singleDeclaration)

	for p.tokensLeft() && p.mustGetCurrentToken().Type == tokenizer.Semicolon {
		p.advance()
		single, err := p.SingleDeclaration()
		if err != nil {
			return nil, err
		}
		node.AddChild(single)
	}

	return node, nil
}

// SingleDeclaration parses the basic singleDeclaration construct
//
//	singleDeclaration ::=
//	         const Identifier ~ expression
//	       | var identifier : typeDenoter
func (p *Parser) SingleDeclaration() (ast.Node, error) {
	currentToken, err := p.getCurrentToken()
	if err != nil {
		return nil, err
	}
	switch currentToken.Type {
	case tokenizer.Const:
		{
			p.advance()
			next, err := p.getCurrentToken()
			if err != nil {
				return nil, err
			}
			if next.Type != tokenizer.Identifier {
				return nil, p.UnexpectedToken(currentToken, tokenizer.Identifier)
			}
			p.advance()
			err = p.expect(tokenizer.Tilde)
			if err != nil {
				return nil, err
			}

			expression, err := p.Expression()
			if err != nil {
				return nil, err
			}

			return ast.NewConstBlock(next, expression), nil
		}
	case tokenizer.Var:
		{
			p.advance()
			next, err := p.getCurrentToken()
			if err != nil {
				return nil, err
			}
			if next.Type != tokenizer.Identifier {
				return nil, p.UnexpectedToken(currentToken, tokenizer.Identifier)
			}
			p.advance()
			err = p.expect(tokenizer.Colon)
			if err != nil {
				return nil, err
			}
			typeDenoter, err := p.TypeDenoter()
			if err != nil {
				return nil, err
			}
			return ast.NewVarBlock(next, typeDenoter), nil
		}
	}
	return nil, nil
}

func isOperator(token *tokenizer.Token) bool {
	return isOneOf(token, tokenizer.PlusOperator, tokenizer.MinusOperator,
		tokenizer.DivisionOperator, tokenizer.MultiplicationOperator, tokenizer.Equals, tokenizer.Comparison,
		tokenizer.LessThan, tokenizer.GreaterThan, tokenizer.LessThanEqual, tokenizer.GreaterThanEqual)
}

// TypeDenoter parses the basic typeDenoter construct
//
// typeDenoter ::= Identifier
func (p *Parser) TypeDenoter() (*ast.TypeDenoter, error) {
	currentToken, err := p.getCurrentToken()
	if err != nil {
		return nil, err
	}
	if currentToken.Type == tokenizer.Identifier {
		p.advance()
		return ast.NewTypeDenoter(currentToken), nil
	}

	return nil, p.UnexpectedToken(currentToken, tokenizer.Identifier)
}

func (p *Parser) tokensLeft() bool {
	return p.currentToken < len(p.tokens)
}

// Expression parses the expression construct
//
// expression ::= primaryExpression (operator primaryExpression)*
func (p *Parser) Expression() (*ast.Expression, error) {
	node := ast.NewExpression()
	primaryExpressionNode, err := p.PrimaryExpression()
	if err != nil {
		return nil, err
	}
	node.AddChild(primaryExpressionNode)

	for p.tokensLeft() && isOperator(p.mustGetCurrentToken()) {
		operator, err := p.getCurrentToken()
		if err != nil {
			return nil, err
		}
		operatorNode := ast.NewOperator(operator.Type)
		node.AddChild(operatorNode)
		p.advance()
		primaryExpressionNode, err = p.PrimaryExpression()
		if err != nil {
			return nil, err
		}
		node.AddChild(primaryExpressionNode)
	}

	return node, nil
}

func isOneOf(token *tokenizer.Token, types ...tokenizer.TokenType) bool {
	for _, _type := range types {
		if token.Type == _type {
			return true
		}
	}
	return false
}

// PrimaryExpression parses the basic primaryExpression construct
//
// primaryExpression ::= Literal | Identifier | ( expression )
func (p *Parser) PrimaryExpression() (ast.Node, error) {
	currentToken, err := p.getCurrentToken()
	if err != nil {
		return nil, err
	}
	switch currentToken.Type {
	case tokenizer.Integer:
		{
			p.advance()
			value, err := strconv.Atoi(currentToken.Value)
			if err != nil {
				return nil, err
			}
			return ast.NewPrimaryExpressionInteger(value), nil
		}
	case tokenizer.Float:
		{
			p.advance()
			value, err := strconv.ParseFloat(currentToken.Value, 64)
			if err != nil {
				return nil, err
			}
			return ast.NewPrimaryExpressionFloat(value), nil
		}
	case tokenizer.LeftParenthesis:
		{
			p.advance()
			res, err := p.Expression()
			if err != nil {
				return nil, err
			}
			err = p.expect(tokenizer.RightParenthesis)
			return res, err
		}

	case tokenizer.Identifier:
		{
			p.advance()
			return ast.NewPrimaryExpressionIdentifier(*currentToken), nil
		}
	case tokenizer.String:
		{
			p.advance()
			return ast.NewPrimaryExpressionString(currentToken.Value[1:]), nil
		}
	}
	return nil, p.UnexpectedToken(currentToken, tokenizer.Identifier, tokenizer.String, tokenizer.Integer, tokenizer.Float, tokenizer.LeftParenthesis)
}

// Command parses the basic command construct
//
//	command ::= singleCommand (; singleCommand)*
func (p *Parser) Command() (*ast.Command, error) {
	node := ast.NewCommand()
	singleCommand, err := p.SingleCommand()
	if err != nil {
		return nil, err
	}

	node.AddChild(singleCommand)

	for p.tokensLeft() && p.mustGetCurrentToken().Type == tokenizer.Semicolon {
		p.advance()
		single, err := p.SingleCommand()
		if err != nil {
			return nil, err
		}

		node.AddChild(single)
	}
	return node, nil
}
