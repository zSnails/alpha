package ast

type NodeType int8

const (
	Program NodeType = iota
	Command
	SingleCommand
	Declaration
	SingleDeclaration
	TypeDenoter
	Expression
	PrimaryExpression
	Operator
	Integer
	Float
	Identifier
	String

	Equals
	If
	Then
	Else
	While
	Do
	Let
	In
	Begin
	End
)

type Node struct {
	Type     NodeType `json:"type"`
	Value    any      `json:"value,omitempty"`
	Children []*Node  `json:"children,omitempty"`
}

func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
}

func NewNode(_type NodeType, value any) *Node {
	return &Node{
		Type:     _type,
		Value:    value,
		Children: []*Node{},
	}
}
