package parser

import (
	"fmt"
	"lang_test/parser/ast"
	"lang_test/tokenizer"
	"strconv"
)

type Parser struct {
	tokens       []*tokenizer.Token
	currentToken int
}

func (p *Parser) GetCurrentToken() *tokenizer.Token {
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
	}, nil
}

func (p *Parser) expect(_type tokenizer.TokenType) error {
	token := p.GetCurrentToken()
	if token.Type != _type {
		return fmt.Errorf("expected token %s, got %s\n", tokenizer.TokenNames[_type], tokenizer.TokenNames[token.Type])
	}
	p.advance()

	return nil
}

func (p *Parser) advance() {
	p.currentToken++
}

func (p *Parser) Program() (*ast.Node, error) {
	// return p.SingleCommand()
	return p.Command()
}

func (p *Parser) SingleCommand() (*ast.Node, error) {
	currentToken := p.GetCurrentToken()
	node := ast.NewNode(ast.SingleCommand, nil)
	switch currentToken.Type {
	case tokenizer.Identifier:
		{
			node.AddChild(ast.NewNode(ast.Identifier, currentToken.Value))
			p.acceptIt()
			next := p.GetCurrentToken()
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
					next = p.GetCurrentToken()
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
	return nil, fmt.Errorf("unexpected token '%s' while building SingleCommand\n", tokenizer.TokenNames[currentToken.Type])
}

func (p *Parser) Declaration() (*ast.Node, error) {
	node := ast.NewNode(ast.Declaration, nil)
	singleDeclaration, err := p.SingleDeclaration()
	if err != nil {
		return nil, err
	}
	node.AddChild(singleDeclaration)

	for p.tokensLeft() && p.GetCurrentToken().Type == tokenizer.Semicolon {
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
	currentToken := p.GetCurrentToken()
	node := ast.NewNode(ast.SingleDeclaration, nil)
	switch currentToken.Type {
	case tokenizer.Const:
		{
			node.AddChild(ast.NewNode(ast.Const, nil))
			p.acceptIt()
			next := p.GetCurrentToken()
			if next.Type != tokenizer.Identifier {
				return nil, fmt.Errorf("unexpected token '%s' while building SingleDeclaration, expected Identifier\n", tokenizer.TokenNames[next.Type])
			}
			p.acceptIt()
			node.AddChild(ast.NewNode(ast.Identifier, next.Value))

			err := p.expect(tokenizer.Tilde)
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

			next := p.GetCurrentToken()
			if next.Type != tokenizer.Identifier {
				return nil, fmt.Errorf("unexpected token '%s' while building SingleDeclaration, expected Identifier\n", tokenizer.TokenNames[next.Type])
			}
			node.AddChild(ast.NewNode(ast.Identifier, next.Value))
			p.acceptIt()
			err := p.expect(tokenizer.Colon)
			if err != nil {
				return nil, err
			}
			typeDenoter, err := p.TypeDenoter()
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
	return token.Type == tokenizer.PlusOperator ||
		token.Type == tokenizer.MinusOperator ||
		token.Type == tokenizer.DivisionOperator ||
		token.Type == tokenizer.MultiplicationOperator ||
		token.Type == tokenizer.Equals ||
		token.Type == tokenizer.Comparison ||
		token.Type == tokenizer.LessThan ||
		token.Type == tokenizer.GreaterThan ||
		token.Type == tokenizer.LessThanEqual ||
		token.Type == tokenizer.GreaterThanEqual
}

func (p *Parser) TypeDenoter() (*ast.Node, error) {
	currentToken := p.GetCurrentToken()
	if currentToken.Type == tokenizer.Identifier {
		p.acceptIt()
		return ast.NewNode(ast.TypeDenoter, currentToken.Value), nil
	}

	return nil, fmt.Errorf("unexpected token '%s' while building TypeDenoter, expected Identifier\n", tokenizer.TokenNames[currentToken.Type])
}

func (p *Parser) tokensLeft() bool {
	return p.currentToken < len(p.tokens)-1
}

func (p *Parser) Expression() (*ast.Node, error) {
	node := ast.NewNode(ast.Expression, nil)
	primaryExpressionNode, err := p.PrimaryExpression()
	if err != nil {
		return nil, err
	}
	node.AddChild(primaryExpressionNode)

	for p.tokensLeft() && isOperator(p.GetCurrentToken()) {
		operator := p.GetCurrentToken()
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
	currentToken := p.GetCurrentToken()
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
	return nil, fmt.Errorf("unexpected token '%s' while producing PrimaryExpression\n", currentToken)
}

func (p *Parser) acceptIt() {
	if p.currentToken > len(p.tokens)-1 {
		panic("lmao")
	}
	p.advance()
}

func (p *Parser) Command() (*ast.Node, error) {
	node := ast.NewNode(ast.Command, nil)
	singleCommand, err := p.SingleCommand()
	if err != nil {
		return nil, err
	}
	node.AddChild(singleCommand)

	if p.tokensLeft() && p.GetCurrentToken().Type == tokenizer.Semicolon {
		p.acceptIt()
		single, err := p.SingleCommand()
		if err != nil {
			return nil, err
		}
		node.AddChild(single)
	}
	return node, nil
}
