package ast

import (
	"fmt"
	"strings"
)

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
	Const
	Var
	Tilde
	In
	Begin
	End
)

var ConstructNames = map[NodeType]string{

	Program:           "Program",
	Command:           "Command",
	SingleCommand:     "SingleCommand",
	Declaration:       "Declaration",
	SingleDeclaration: "SingleDeclaration",
	TypeDenoter:       "TypeDenoter",
	Expression:        "Expression",
	PrimaryExpression: "PrimaryExpression",
	Operator:          "Operator",
	Integer:           "Integer",
	Float:             "Float",
	Identifier:        "Identifier",
	String:            "String",

	Equals: "=",
	If:     "if",
	Then:   "then",
	Else:   "else",
	While:  "while",
	Do:     "do",
	Let:    "let",
	Const:  "const",
	Var:    "var",
	Tilde:  "~",
	In:     "in",
	Begin:  "begin",
	End:    "end",
}

type Node struct {
	Type     NodeType `json:"type"`
	Value    any      `json:"value,omitempty"`
	Children []*Node  `json:"children,omitempty"`
}

func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
}

func (n *Node) String() string {
	if n == nil {
		return ""
	}
	var sb strings.Builder
	n.buildSExp(&sb, 0)
	return sb.String()
}

func (n *Node) buildSExp(sb *strings.Builder, level int) {
    if n == nil {
        return
    }
    indent := strings.Repeat("    ", level)

	switch v := n.Value.(type) {
	case string:
        fmt.Fprintf(sb, "%s\"%s\"", indent, v)
	case nil:
		fmt.Fprintf(sb, "%s%s", indent, ConstructNames[n.Type])
	default:
		fmt.Fprintf(sb, "%s%v", indent, n.Value)
	}

    if len(n.Children) > 0 {
        fmt.Fprint(sb, " (")
        for _, child := range n.Children {
            fmt.Fprintln(sb)
            child.buildSExp(sb, level+1)
        }
        fmt.Fprint(sb, ")")
    }
}

func NewNode(_type NodeType, value any) *Node {
	return &Node{
		Type:     _type,
		Value:    value,
		Children: []*Node{},
	}
}
