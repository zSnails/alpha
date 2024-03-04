package parser

import (
	"errors"
	"fmt"
	"io"
	"lang_test/parser/ast"
	"lang_test/tokenizer"
	"strconv"
	"strings"
)

type Parser struct {
	tokens       []*tokenizer.Token
	currentToken int
	lexer        *tokenizer.Tokenizer
}

func (p *Parser) GetCurrentToken() (*tokenizer.Token, error) {
	if !p.tokensLeft() {
		return nil, io.EOF
	}
	return p.tokens[p.currentToken], nil
}

func (p *Parser) MustGetCurrentToken() *tokenizer.Token {
	return p.tokens[p.currentToken]
}

// EBNF for this piece of shit
// program ::= if expression then expression (; expression)* end
func NewParser(lexer *tokenizer.Tokenizer) (*Parser, error) {
	tokens, err := lexer.GetAllTokens()
	if err != nil {
		return nil, err
	}
	return &Parser{
		tokens: tokens,
		lexer:  lexer,
	}, nil
}

func (p *Parser) expect(_type tokenizer.TokenType) error {
	token, err := p.GetCurrentToken()
	if err != nil {
		return p.UnexpectedTokenExpected(EOF, _type)
	}

	if token.Type != _type {
		return p.UnexpectedTokenExpected(token, _type)
	}
	p.advance()

	return nil
}

func (p *Parser) advance() {
	p.currentToken++
}

func (p *Parser) Program() (*ast.Node, error) {
	node, err := p.Command()
	if err != nil {
		return nil, err
	}

	if p.tokensLeft() {
		return nil, p.UnexpectedToken(p.MustGetCurrentToken())
	}
	return node, nil
}

func (p *Parser) SingleCommand() (*ast.Node, error) {
	node := ast.NewNode(ast.SingleCommand, nil)
	currentToken, err := p.GetCurrentToken() // this error will always be io.EOF
	if err != nil {
		return nil, err
	}

	switch currentToken.Type {
	case tokenizer.Identifier:
		{
			node.AddChild(ast.NewNode(ast.Identifier, currentToken.Value))
			p.acceptIt()
			next, err := p.GetCurrentToken()
			if err != nil {
				return nil, p.UnexpectedTokenExpectedOneOf(EOF, tokenizer.Equals, tokenizer.LeftParenthesis)
			}
			switch next.Type {
			case tokenizer.Equals:
				{
					node.AddChild(ast.NewNode(ast.Equals, next.Value))
					p.acceptIt()
					expressionNode, err := p.Expression()
					if err != nil {
						return nil, err
					}
					node.AddChild(expressionNode)
					return node, nil
				}
			case tokenizer.LeftParenthesis:
				{
					p.acceptIt()
					next, err = p.GetCurrentToken()
					if err != nil {
						return nil, p.UnexpectedTokenExpectedOneOf(EOF, tokenizer.RightParenthesis, tokenizer.Integer, tokenizer.Float, tokenizer.String)
					}
					if next.Type == tokenizer.RightParenthesis {
						p.acceptIt()
						return node, nil
					}

					expressionNode, err := p.Expression()
					if err != nil {
						return nil, err
					}
					err = p.expect(tokenizer.RightParenthesis)
					if err != nil {
						return nil, err
					}
					node.AddChild(expressionNode)
					return node, nil
				}
			}
		}
	case tokenizer.If:
		{
			node.AddChild(ast.NewNode(ast.If, nil))
			p.acceptIt()
			expressionNode, err := p.Expression()
			if err != nil {
				return nil, err
			}
			node.AddChild(expressionNode)
			err = p.expect(tokenizer.Then)
			if err != nil {
				return nil, err
			}
			ifBlockSingleCommand, err := p.SingleCommand()
			if err != nil {
				return nil, err
			}
			node.AddChild(ifBlockSingleCommand)
			err = p.expect(tokenizer.Else)
			if err != nil {
				return nil, err
			}
			elseBlockSingleCommand, err := p.SingleCommand()
			if err != nil {
				return nil, err
			}
			node.AddChild(elseBlockSingleCommand)
			return node, nil
		}
	case tokenizer.While:
		{
			node.AddChild(ast.NewNode(ast.While, nil))
			p.acceptIt()
			while, err := p.Expression()
			if err != nil {
				return nil, err
			}
			node.AddChild(while)
			err = p.expect(tokenizer.Do)
			if err != nil {
				return nil, err
			}
			singleCommand, err := p.SingleCommand()
			if err != nil {
				return nil, err
			}
			node.AddChild(singleCommand)
			return node, nil
		}

	case tokenizer.Let:
		{
			node.AddChild(ast.NewNode(ast.Let, nil))
			p.acceptIt()

			declaration, err := p.Declaration()
			if err != nil {
				return nil, err
			}
			node.AddChild(declaration)

			err = p.expect(tokenizer.In)
			if err != nil {
				return nil, err
			}

			singleCommand, err := p.SingleCommand()
			if err != nil {
				return nil, err
			}

			node.AddChild(singleCommand)
			return node, nil
		}
	case tokenizer.Begin:
		{
			p.acceptIt()
			command, err := p.Command()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil, p.UnexpectedToken(EOF)
				}
				return nil, err
			}
			node.AddChild(command)
			err = p.expect(tokenizer.End)
			if err != nil {
				return nil, err
			}
			return node, nil
		}
	}
	return nil, p.UnexpectedToken(currentToken)
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

func (p *Parser) UnexpectedTokenExpectedOneOf(got *tokenizer.Token, exptected ...tokenizer.TokenType) error {
	tokens := Map(exptected, func(token tokenizer.TokenType) string {
		return fmt.Sprintf("'%s'", tokenizer.TokenNames[token])
	})
	expectedTokens := strings.Join(tokens, ", ")

	row, col := got.GetPosition()
	return fmt.Errorf("%s:%d:%d: unexpected token '%s' expected one of %s", p.lexer.GetFileName(), row, col, tokenizer.TokenNames[got.Type], expectedTokens)
}

func (p *Parser) UnexpectedTokenExpected(got *tokenizer.Token, expected tokenizer.TokenType) error {
	row, col := got.GetPosition()
	return fmt.Errorf("%s:%d:%d: unexpected token '%s' expected '%s'\n", p.lexer.GetFileName(), row, col, tokenizer.TokenNames[got.Type], tokenizer.TokenNames[expected])
}

func (p *Parser) UnexpectedToken(token *tokenizer.Token) error {
	row, col := token.GetPosition()
	return fmt.Errorf("%s:%d:%d: unexpected token '%s'\n", p.lexer.GetFileName(), row, col, tokenizer.TokenNames[token.Type])
}

func (p *Parser) Declaration() (*ast.Node, error) {
	node := ast.NewNode(ast.Declaration, nil)
	singleDeclaration, err := p.SingleDeclaration()
	if err != nil {
		return nil, err
	}
	node.AddChild(singleDeclaration)

	for p.tokensLeft() && p.MustGetCurrentToken().Type == tokenizer.Semicolon {
		p.acceptIt()
		single, err := p.SingleDeclaration()
		if err != nil {
			return nil, err
		}
		node.AddChild(single)
	}

	return node, nil
}

func (p *Parser) SingleDeclaration() (*ast.Node, error) {
	currentToken, err := p.GetCurrentToken()
	if err != nil {
		return nil, err
	}
	node := ast.NewNode(ast.SingleDeclaration, nil)
	switch currentToken.Type {
	case tokenizer.Const:
		{
			node.AddChild(ast.NewNode(ast.Const, nil))
			p.acceptIt()
			next, err := p.GetCurrentToken()
			if err != nil {
				return nil, err
			}
			if next.Type != tokenizer.Identifier {
				return nil, p.UnexpectedTokenExpected(currentToken, tokenizer.Identifier)
			}
			p.acceptIt()
			node.AddChild(ast.NewNode(ast.Identifier, next.Value))

			err = p.expect(tokenizer.Tilde)
			if err != nil {
				return nil, err
			}

			expression, err := p.Expression()
			if err != nil {
				return nil, err
			}

			node.AddChild(expression)
			return node, nil
		}
	case tokenizer.Var:
		{
			node.AddChild(ast.NewNode(ast.Var, nil))
			p.acceptIt()

			next, err := p.GetCurrentToken()
			if err != nil {
				return nil, err
			}
			if next.Type != tokenizer.Identifier {
				return nil, p.UnexpectedTokenExpected(currentToken, tokenizer.Identifier)
			}
			node.AddChild(ast.NewNode(ast.Identifier, next.Value))
			p.acceptIt()
			err = p.expect(tokenizer.Colon)
			if err != nil {
				return nil, err
			}
			typeDenoter, err := p.TypeHint()
			if err != nil {
				return nil, err
			}
			node.AddChild(typeDenoter)
			return node, nil
		}
	}
	return nil, nil
}

func isOperator(token *tokenizer.Token) bool {
	return isOneOf(token, tokenizer.PlusOperator, tokenizer.MinusOperator,
		tokenizer.DivisionOperator, tokenizer.MultiplicationOperator, tokenizer.Equals, tokenizer.Comparison,
		tokenizer.LessThan, tokenizer.GreaterThan, tokenizer.LessThanEqual, tokenizer.GreaterThanEqual)
}

func (p *Parser) TypeHint() (*ast.Node, error) {
	currentToken, err := p.GetCurrentToken()
	if err != nil {
		return nil, err
	}
	if currentToken.Type == tokenizer.Identifier {
		p.acceptIt()
		return ast.NewNode(ast.TypeDenoter, currentToken.Value), nil
	}

	return nil, p.UnexpectedTokenExpected(currentToken, tokenizer.Identifier)
}

func (p *Parser) tokensLeft() bool {
	return p.currentToken < len(p.tokens)
}

func (p *Parser) Expression() (*ast.Node, error) {
	node := ast.NewNode(ast.Expression, nil)
	primaryExpressionNode, err := p.PrimaryExpression()
	if err != nil {
		return nil, err
	}
	node.AddChild(primaryExpressionNode)

	for p.tokensLeft() && isOperator(p.MustGetCurrentToken()) {
		operator, err := p.GetCurrentToken()
		if err != nil {
			return nil, err
		}
		operatorNode := ast.NewNode(ast.Operator, operator)
		node.AddChild(operatorNode)
		p.acceptIt()
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

func (p *Parser) PrimaryExpression() (*ast.Node, error) {
	currentToken, err := p.GetCurrentToken()
	if err != nil {
		return nil, err
	}
	switch currentToken.Type {
	case tokenizer.Integer:
		{
			p.acceptIt()
			value, err := strconv.Atoi(currentToken.Value)
			if err != nil {
				return nil, err
			}
			return ast.NewNode(ast.Integer, value), nil
		}
	case tokenizer.Float:
		{
			p.acceptIt()
			value, err := strconv.ParseFloat(currentToken.Value, 64)
			if err != nil {
				return nil, err
			}
			return ast.NewNode(ast.Integer, value), nil
		}
	case tokenizer.LeftParenthesis:
		{
			p.acceptIt()
			res, err := p.Expression()
			if err != nil {
				return nil, err
			}
			err = p.expect(tokenizer.RightParenthesis)
			return res, err
		}

	case tokenizer.Identifier:
		{
			p.acceptIt()
			return ast.NewNode(ast.Identifier, currentToken.Value), nil
		}
	case tokenizer.String:
		{
			p.acceptIt()
			return ast.NewNode(ast.String, currentToken.Value[1:]), nil
		}
	}
	return nil, p.UnexpectedToken(currentToken)
}

func (p *Parser) acceptIt() {
	p.advance()
}

func (p *Parser) Command() (*ast.Node, error) {
	node := ast.NewNode(ast.Command, nil)
	singleCommand, err := p.SingleCommand()
	if err != nil {
		return nil, err
	}

	node.AddChild(singleCommand)

	for p.tokensLeft() && p.MustGetCurrentToken().Type == tokenizer.Semicolon {
		p.acceptIt()
		single, err := p.SingleCommand()
		if err != nil {
			return nil, err
		}

		node.AddChild(single)
	}
	return node, nil
}
