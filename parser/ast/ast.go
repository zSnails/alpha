package ast

import (
	"fmt"
	"strings"

	"github.com/zSnails/alpha/tokenizer"
)

type Node interface {
	fmt.Stringer

	// AddChild add a child to this node's `Children` slice.
	//
	// The implementation may be left blank
	AddChild(Node)

	// GetChildren returns the node's children.
	//
	// The implementation of this method may return nil for nodes that don't have any children
	GetChildren() []Node

	buildSExp(*strings.Builder, int)
}

type Program struct {
	SingleCommand Node `json:"singleCommand"`
}

// GetChildren implements Node.
func (*Program) GetChildren() []Node {
	return nil
}

// AddChild implements Node.
func (*Program) AddChild(Node) {
	panic("unimplemented")
}

func NewProgram(singleCommand Node) *Program {
	return &Program{
		SingleCommand: singleCommand,
	}
}

type Assignment struct {
	Identifier *tokenizer.Token `json:"identifier"`
	Expression *Expression      `json:"expression"`
}

// GetChildren implements Node.
func (*Assignment) GetChildren() []Node {
	return nil
}

// AddChild implements Node.
func (*Assignment) AddChild(Node) {
	panic("unimplemented")
}

func NewAssignment(identifier *tokenizer.Token, expression *Expression) Node {
	return &Assignment{
		Identifier: identifier,
		Expression: expression,
	}
}

type FunctionCall struct {
	Identifier *tokenizer.Token `json:"identifier"`
	Expression *Expression      `json:"expression"`
}

// GetChildren implements Node.
func (*FunctionCall) GetChildren() []Node {
	return nil
}

// AddChild implements Node.
func (*FunctionCall) AddChild(Node) {
	panic("unimplemented")
}

func NewFunctionCall(identifier *tokenizer.Token, expression *Expression) *FunctionCall {
	return &FunctionCall{
		Identifier: identifier,
		Expression: expression,
	}
}

type Expression struct {
	Children []Node `json:"children"` // TODO: fix this piece of shit
}

// GetChildren implements Node.
func (e *Expression) GetChildren() []Node {
	return e.Children
}

// AddChild implements Node.
func (e *Expression) AddChild(node Node) {
	e.Children = append(e.Children, node)
}

func NewExpression() *Expression {
	return &Expression{
		Children: []Node{},
	}
}

type TypeDenoter struct {
	Name *tokenizer.Token `json:"name"`
}

// GetChildren implements Node.
func (*TypeDenoter) GetChildren() []Node {
	return nil
}

// AddChild implements Node.
func (*TypeDenoter) AddChild(Node) {
	panic("unimplemented")
}

func NewTypeDenoter(identifier *tokenizer.Token) *TypeDenoter {
	return &TypeDenoter{}
}

type PrimaryExpression struct {
	Expression *Expression `json:"expression,omitempty"`
}

// AddChild implements Node.
func (*PrimaryExpression) AddChild(Node) {
	panic("unimplemented")
}

type PrimaryExpressionInteger int

// GetChildren implements Node.
func (PrimaryExpressionInteger) GetChildren() []Node {
	return nil
}

// AddChild implements Node.
func (PrimaryExpressionInteger) AddChild(Node) {
	panic("unimplemented")
}

func NewPrimaryExpressionInteger(value int) PrimaryExpressionInteger {
	return PrimaryExpressionInteger(value)
}

type PrimaryExpressionString string

// GetChildren implements Node.
func (PrimaryExpressionString) GetChildren() []Node {
	return nil
}

// AddChild implements Node.
func (PrimaryExpressionString) AddChild(Node) {
	panic("unimplemented")
}

func NewPrimaryExpressionString(value string) PrimaryExpressionString {
	return PrimaryExpressionString(value)
}

type PrimaryExpressionFloat float64

// GetChildren implements Node.
func (PrimaryExpressionFloat) GetChildren() []Node {
	return nil
}

// AddChild implements Node.
func (PrimaryExpressionFloat) AddChild(Node) {
	panic("unimplemented")
}

func NewPrimaryExpressionFloat(float float64) PrimaryExpressionFloat {
	return PrimaryExpressionFloat(float)
}

type PrimaryExpressionIdentifier tokenizer.Token

// GetChildren implements Node.
func (PrimaryExpressionIdentifier) GetChildren() []Node {
	return nil
}

// AddChild implements Node.
func (PrimaryExpressionIdentifier) AddChild(Node) {
	panic("unimplemented")
}

func NewPrimaryExpressionIdentifier(identifier tokenizer.Token) PrimaryExpressionIdentifier {
	return PrimaryExpressionIdentifier(identifier)
}

type Operator struct {
	Type tokenizer.TokenType
}

// GetChildren implements Node.
func (*Operator) GetChildren() []Node {
	return nil
}

// AddChild implements Node.
func (*Operator) AddChild(Node) {
	panic("unimplemented")
}

func NewOperator(_type tokenizer.TokenType) Node {
	return &Operator{
		Type: _type,
	}
}

type ConstBlock struct {
	Identifier *tokenizer.Token `json:"identifier"`
	Expression *Expression      `json:"expression"`
}

// GetChildren implements Node.
func (*ConstBlock) GetChildren() []Node {
	return nil
}

// AddChild implements Node.
func (*ConstBlock) AddChild(Node) {
	panic("unimplemented")
}

func NewConstBlock(identifier *tokenizer.Token, expression *Expression) *ConstBlock {
	return &ConstBlock{
		Identifier: identifier,
		Expression: expression,
	}
}

type VarBlock struct {
	Identifier *tokenizer.Token
	Type       *TypeDenoter
}

// GetChildren implements Node.
func (*VarBlock) GetChildren() []Node {
	return nil
}

// AddChild implements Node.
func (*VarBlock) AddChild(Node) {
	panic("unimplemented")
}

func NewVarBlock(identifier *tokenizer.Token, _type *TypeDenoter) *VarBlock {
	return &VarBlock{
		Identifier: identifier,
		Type:       _type,
	}
}

type Declaration struct {
	Children []Node `json:"children"`
}

// GetChildren implements Node.
func (d *Declaration) GetChildren() []Node {
	return d.Children
}

// AddChild implements Node.
func (d *Declaration) AddChild(node Node) {
	d.Children = append(d.Children, node)
}

func NewDeclaration() *Declaration {
	return &Declaration{
		Children: []Node{},
	}
}

type IfBlock struct {
	Condition *Expression `json:"condition"`
	IfBody    Node        `json:"ifBody"`
	ElseBody  Node        `json:"elseBody"`
}

// GetChildren implements Node.
func (*IfBlock) GetChildren() []Node {
	return nil
}

// AddChild implements Node.
func (*IfBlock) AddChild(Node) {
	panic("unimplemented")
}

func NewIfBlock(condition *Expression, ifBody Node, elseBody Node) *IfBlock {
	return &IfBlock{
		Condition: condition,
		IfBody:    ifBody,
		ElseBody:  elseBody,
	}
}

type WhileBlock struct {
	Condition *Expression `json:"condition"`
	Body      Node        `json:"body"`
}

// GetChildren implements Node.
func (*WhileBlock) GetChildren() []Node {
	return nil
}

// AddChild implements Node.
func (*WhileBlock) AddChild(Node) {
	panic("unimplemented")
}

func NewWhileBlock(condition *Expression, body Node) *WhileBlock {
	return &WhileBlock{
		Condition: condition,
		Body:      body,
	}
}

type LetBlock struct {
	Declaration *Declaration `json:"declaration"`
	Body        Node         `json:"body"`
}

// GetChildren implements Node.
func (*LetBlock) GetChildren() []Node {
	return nil
}

// AddChild implements Node.
func (*LetBlock) AddChild(Node) {
	panic("unimplemented")
}

func NewLetBlock(declaration *Declaration, body Node) *LetBlock {
	return &LetBlock{
		Declaration: declaration,
		Body:        body,
	}
}

type BeginBlock struct {
	Body *Command `json:"beginBody"`
}

// GetChildren implements Node.
func (*BeginBlock) GetChildren() []Node {
	return nil
}

// AddChild implements Node.
func (*BeginBlock) AddChild(Node) {
	panic("unimplemented")
}

func NewBeginBlock(body *Command) *BeginBlock {
	return &BeginBlock{
		Body: body,
	}
}

type Command struct {
	Children []Node `json:"children"`
}

// AddChild implements Node.
func (c *Command) AddChild(node Node) {
	c.Children = append(c.Children, node)
}

func (c *Command) GetChildren() []Node {
	return c.Children
}

func NewCommand() *Command {
	return &Command{
		Children: []Node{},
	}
}

func (p *Program) String() string {
	if p == nil {
		return ""
	}

	var sb strings.Builder
	p.buildSExp(&sb, 0)
	return sb.String()
}

func (p *Program) buildSExp(sb *strings.Builder, level int) {
	if p == nil {
		return
	}
	indent := strings.Repeat("    ", level)
	fmt.Fprintf(sb, "%sProgram", indent)
	fmt.Fprint(sb, " (")
	fmt.Fprintln(sb)
	p.SingleCommand.buildSExp(sb, level+1)
	fmt.Fprint(sb, ")")
}

func (d *Declaration) String() string {
	if d == nil {
		return ""
	}

	var sb strings.Builder
	d.buildSExp(&sb, 0)
	return sb.String()
}

func (d *Declaration) buildSExp(sb *strings.Builder, level int) {
	if d == nil {
		return
	}
	indent := strings.Repeat("    ", level)
	fmt.Fprintf(sb, "%sDeclaration", indent)
	if len(d.Children) > 0 {
		fmt.Fprint(sb, " (")
		for _, child := range d.Children {
			fmt.Fprintln(sb)
			child.buildSExp(sb, level+1)
		}
		fmt.Fprint(sb, ")")
	}
}

func (c *Command) String() string {
	if c == nil {
		return ""
	}

	var sb strings.Builder
	c.buildSExp(&sb, 0)
	return sb.String()
}

func (c *Command) buildSExp(sb *strings.Builder, level int) {
	if c == nil {
		return
	}
	indent := strings.Repeat("    ", level)
	fmt.Fprintf(sb, "%sCommand", indent)
	if len(c.Children) > 0 {
		fmt.Fprint(sb, " (")
		for _, child := range c.Children {
			fmt.Fprintln(sb)
			child.buildSExp(sb, level+1)
		}
		fmt.Fprint(sb, ")")
	}
}

func (c *BeginBlock) String() string {
	if c == nil {
		return ""
	}

	var sb strings.Builder
	c.buildSExp(&sb, 0)
	return sb.String()
}

func (c *BeginBlock) buildSExp(sb *strings.Builder, level int) {
	if c == nil {
		return
	}
	indent := strings.Repeat("    ", level)
	fmt.Fprintf(sb, "%sBeginBlock", indent)

	fmt.Fprint(sb, " (")
	fmt.Fprintln(sb)
	c.Body.buildSExp(sb, level+1)
	fmt.Fprint(sb, ")")
}

func (c *FunctionCall) String() string {
	if c == nil {
		return ""
	}

	var sb strings.Builder
	c.buildSExp(&sb, 0)
	return sb.String()
}

func (c *FunctionCall) buildSExp(sb *strings.Builder, level int) {
	if c == nil {
		return
	}
	indent := strings.Repeat("    ", level)
	inner := strings.Repeat("    ", level+1)
	fmt.Fprintf(sb, "%sFunctionCall", indent)

	fmt.Fprint(sb, " (")
	fmt.Fprintln(sb)
	fmt.Fprintf(sb, "%s%s", inner, c.Identifier.Value)
	fmt.Fprintln(sb)
	c.Expression.buildSExp(sb, level+1)
	fmt.Fprint(sb, ")")

}

func (c *Expression) String() string {
	if c == nil {
		return ""
	}

	var sb strings.Builder
	c.buildSExp(&sb, 0)
	return sb.String()
}

func (c *Expression) buildSExp(sb *strings.Builder, level int) {
	if c == nil {
		return
	}
	indent := strings.Repeat("    ", level)
	fmt.Fprintf(sb, "%sExpression", indent)
	if len(c.Children) > 0 {
		fmt.Fprint(sb, " (")
		for _, child := range c.Children {
			fmt.Fprintln(sb)
			child.buildSExp(sb, level+1)
		}
		fmt.Fprint(sb, ")")
	}
}

func (c PrimaryExpressionString) String() string {
	var sb strings.Builder
	c.buildSExp(&sb, 0)
	return sb.String()
}

func (c PrimaryExpressionString) buildSExp(sb *strings.Builder, level int) {
	indent := strings.Repeat("    ", level)
	inner := strings.Repeat("    ", level+1)
	fmt.Fprintf(sb, "%sPrimaryExpressionString", indent)

	fmt.Fprint(sb, " (")
	fmt.Fprintln(sb)
	fmt.Fprintf(sb, "%s'%s'", inner, c)
	fmt.Fprint(sb, ")")
}

func (c *LetBlock) String() string {
	if c == nil {
		return ""
	}

	var sb strings.Builder
	c.buildSExp(&sb, 0)
	return sb.String()
}

func (c *LetBlock) buildSExp(sb *strings.Builder, level int) {
	indent := strings.Repeat("    ", level)
	fmt.Fprintf(sb, "%sLetBlock", indent)

	fmt.Fprint(sb, " (")
	fmt.Fprintln(sb)
	c.Declaration.buildSExp(sb, level+1)
	fmt.Fprintln(sb)
	c.Body.buildSExp(sb, level+1)
	fmt.Fprint(sb, ")")
}

func (c *ConstBlock) String() string {
	if c == nil {
		return ""
	}

	var sb strings.Builder
	c.buildSExp(&sb, 0)
	return sb.String()
}

func (c *ConstBlock) buildSExp(sb *strings.Builder, level int) {
	indent := strings.Repeat("    ", level)
	inner := strings.Repeat("    ", level+1)
	fmt.Fprintf(sb, "%sConstBlock", indent)

	fmt.Fprint(sb, " (")
	// c.Declaration.buildSExp(sb, level+1)
	fmt.Fprintf(sb, "%s%s", inner, c.Identifier.Value)
	fmt.Fprintln(sb)
	c.Expression.buildSExp(sb, level+1)
	// c.Body.buildSExp(sb, level+1)
	fmt.Fprint(sb, ")")
}

func (c *VarBlock) String() string {
	if c == nil {
		return ""
	}

	var sb strings.Builder
	c.buildSExp(&sb, 0)
	return sb.String()
}

func (c *VarBlock) buildSExp(sb *strings.Builder, level int) {
	indent := strings.Repeat("    ", level)
	inner := strings.Repeat("    ", level+1)
	fmt.Fprintf(sb, "%sVarBlock", indent)

	fmt.Fprint(sb, " (")
	fmt.Fprintln(sb)
	fmt.Fprintf(sb, "%s%s", inner, c.Identifier.Value)
	fmt.Fprintln(sb)
	c.Type.buildSExp(sb, level+1)
	fmt.Fprint(sb, ")")
}

func (c *TypeDenoter) String() string {
	if c == nil {
		return ""
	}

	var sb strings.Builder
	c.buildSExp(&sb, 0)
	return sb.String()
}

func (c *TypeDenoter) buildSExp(sb *strings.Builder, level int) {
	indent := strings.Repeat("    ", level)
	fmt.Fprintf(sb, "%sTypeDenoter %s", indent, c.Name)
}

func (c *IfBlock) String() string {
	if c == nil {
		return ""
	}

	var sb strings.Builder
	c.buildSExp(&sb, 0)
	return sb.String()
}

func (c *IfBlock) buildSExp(sb *strings.Builder, level int) {
	indent := strings.Repeat("    ", level)
	fmt.Fprintf(sb, "%sIfBlock", indent)
	fmt.Fprint(sb, " (")
	fmt.Fprintln(sb)
	c.Condition.buildSExp(sb, level)
	fmt.Fprintln(sb)
	c.IfBody.buildSExp(sb, level)
	fmt.Fprintln(sb)
	c.ElseBody.buildSExp(sb, level)
	fmt.Fprint(sb, ")")
}

func (c PrimaryExpressionIdentifier) String() string {
	var sb strings.Builder
	c.buildSExp(&sb, 0)
	return sb.String()
}

func (c PrimaryExpressionIdentifier) buildSExp(sb *strings.Builder, level int) {
	indent := strings.Repeat("    ", level)
	inner := strings.Repeat("    ", level+1)
	fmt.Fprintf(sb, "%sPrimaryExpressionIdentifier", indent)
	fmt.Fprint(sb, " (")
	fmt.Fprintln(sb)
	fmt.Fprintf(sb, "%s%s", inner, c.Value)
	fmt.Fprint(sb, ")")
}

func (c *Operator) String() string {
	var sb strings.Builder
	c.buildSExp(&sb, 0)
	return sb.String()
}

func (c *Operator) buildSExp(sb *strings.Builder, level int) {
	indent := strings.Repeat("    ", level)
	inner := strings.Repeat("    ", level+1)
	fmt.Fprintf(sb, "%sOperator", indent)
	fmt.Fprint(sb, " (")
	fmt.Fprintln(sb)
	fmt.Fprintf(sb, "%s%s", inner, tokenizer.TokenNames[c.Type])
	fmt.Fprint(sb, ")")
}

func (c PrimaryExpressionInteger) String() string {
	var sb strings.Builder
	c.buildSExp(&sb, 0)
	return sb.String()
}

func (c PrimaryExpressionInteger) buildSExp(sb *strings.Builder, level int) {
	indent := strings.Repeat("    ", level)
	inner := strings.Repeat("    ", level+1)
	fmt.Fprintf(sb, "%sPrimaryExpressionInteger", indent)
	fmt.Fprint(sb, " (")
	fmt.Fprintln(sb)
	fmt.Fprintf(sb, "%s%d", inner, c)
	fmt.Fprint(sb, ")")
}

func (c PrimaryExpressionFloat) String() string {
	var sb strings.Builder
	c.buildSExp(&sb, 0)
	return sb.String()
}

func (c PrimaryExpressionFloat) buildSExp(sb *strings.Builder, level int) {
	indent := strings.Repeat("    ", level)
	inner := strings.Repeat("    ", level+1)
	fmt.Fprintf(sb, "%sPrimaryExpressionFloat", indent)
	fmt.Fprint(sb, " (")
	fmt.Fprintln(sb)
	fmt.Fprintf(sb, "%s%f", inner, c)
	fmt.Fprint(sb, ")")
}

func (c *WhileBlock) String() string {
	if c == nil {
		return ""
	}

	var sb strings.Builder
	c.buildSExp(&sb, 0)
	return sb.String()
}

func (c *WhileBlock) buildSExp(sb *strings.Builder, level int) {
	if c == nil {
		return
	}
	indent := strings.Repeat("    ", level)
	fmt.Fprintf(sb, "%sWhileBlock", indent)

	fmt.Fprint(sb, " (")
	fmt.Fprintln(sb)
	c.Condition.buildSExp(sb, level+1)
	fmt.Fprintln(sb)
	c.Body.buildSExp(sb, level+1)
	fmt.Fprint(sb, ")")
}

func (c *Assignment) String() string {
	if c == nil {
		return ""
	}

	var sb strings.Builder
	c.buildSExp(&sb, 0)
	return sb.String()
}

func (c *Assignment) buildSExp(sb *strings.Builder, level int) {
	indent := strings.Repeat("    ", level)
	inner := strings.Repeat("    ", level+1)

	fmt.Fprintf(sb, "%sAssignment", indent)

	fmt.Fprint(sb, " (")
	fmt.Fprintln(sb)
	fmt.Fprintf(sb, "%s%s", inner, c.Identifier.Value)
	fmt.Fprintln(sb)
	c.Expression.buildSExp(sb, level+1)
	fmt.Fprint(sb, ")")
}
