package ast

import (
	"github.com/m1gwings/treedrawer/tree"

	"github.com/arxeiss/go-expression-calculator/lexer"
)

func ToTreeDrawer(rootNode Node) *tree.Tree {
	t := tree.NewTree(nil)
	rootNode.toTreeDrawer(t)
	return t
}

type Node interface {
	toTreeDrawer(*tree.Tree)
	GetToken() *lexer.Token
}

// Just make sure all node types implements Node interface
var _ Node = &NumericNode{}
var _ Node = &VariableNode{}
var _ Node = &UnaryNode{}
var _ Node = &BinaryNode{}
var _ Node = &AssignNode{}
var _ Node = &FunctionNode{}

type NumericNode struct {
	val   float64
	token *lexer.Token
}

func NewNumericNode(val float64, token *lexer.Token) *NumericNode {
	return &NumericNode{
		val:   val,
		token: token,
	}
}

func (n *NumericNode) toTreeDrawer(t *tree.Tree) {
	t.SetVal(tree.NodeFloat64(n.val))
}
func (n *NumericNode) GetToken() *lexer.Token {
	return n.token
}
func (n *NumericNode) Value() float64 {
	return n.val
}

type VariableNode struct {
	name  string
	token *lexer.Token
}

func NewVariableNode(name string, token *lexer.Token) *VariableNode {
	return &VariableNode{
		name:  name,
		token: token,
	}
}

func (n *VariableNode) toTreeDrawer(t *tree.Tree) {
	t.SetVal(tree.NodeString(n.name))
}
func (n *VariableNode) GetToken() *lexer.Token {
	return n.token
}
func (n *VariableNode) Name() string {
	return n.name
}

type UnaryNode struct {
	next     Node
	operator Operation
	token    *lexer.Token
}

func NewUnaryNode(operator Operation, next Node, token *lexer.Token) *UnaryNode {
	return &UnaryNode{
		operator: operator,
		next:     next,
		token:    token,
	}
}

func (n *UnaryNode) Operator() Operation {
	return n.operator
}
func (n *UnaryNode) Next() Node {
	return n.next
}
func (n *UnaryNode) toTreeDrawer(t *tree.Tree) {
	t.SetVal(tree.NodeString(n.operator.String()))
	n.next.toTreeDrawer(t.AddChild(nil))
}
func (n *UnaryNode) GetToken() *lexer.Token {
	return n.token
}

type BinaryNode struct {
	operator Operation
	left     Node
	right    Node
	token    *lexer.Token
}

func NewBinaryNode(operator Operation, left, right Node, token *lexer.Token) *BinaryNode {
	return &BinaryNode{
		operator: operator,
		left:     left,
		right:    right,
		token:    token,
	}
}

func (n *BinaryNode) Left() Node {
	return n.left
}
func (n *BinaryNode) Right() Node {
	return n.right
}
func (n *BinaryNode) Operator() Operation {
	return n.operator
}
func (n *BinaryNode) toTreeDrawer(t *tree.Tree) {
	t.SetVal(tree.NodeString(n.operator.String()))
	n.left.toTreeDrawer(t.AddChild(nil))
	n.right.toTreeDrawer(t.AddChild(nil))
}
func (n *BinaryNode) GetToken() *lexer.Token {
	return n.token
}

type AssignNode struct {
	left  *VariableNode
	right Node
	token *lexer.Token
}

func NewAssignNode(left *VariableNode, right Node, token *lexer.Token) *AssignNode {
	return &AssignNode{
		left:  left,
		right: right,
		token: token,
	}
}

func (n *AssignNode) Left() *VariableNode {
	return n.left
}
func (n *AssignNode) Right() Node {
	return n.right
}
func (n *AssignNode) toTreeDrawer(t *tree.Tree) {
	t.SetVal(tree.NodeString(lexer.Equal))
	n.left.toTreeDrawer(t.AddChild(nil))
	n.right.toTreeDrawer(t.AddChild(nil))
}
func (n *AssignNode) GetToken() *lexer.Token {
	return n.token
}

type FunctionNode struct {
	name  string
	param Node
	token *lexer.Token
}

func NewFunctionNode(name string, param Node, token *lexer.Token) *FunctionNode {
	return &FunctionNode{
		name:  name,
		param: param,
		token: token,
	}
}

func (n *FunctionNode) Param() Node {
	return n.param
}
func (n *FunctionNode) Name() string {
	return n.name
}
func (n *FunctionNode) toTreeDrawer(t *tree.Tree) {
	t.SetVal(tree.NodeString(n.name + "()"))
	n.param.toTreeDrawer(t.AddChild(nil))
}
func (n *FunctionNode) GetToken() *lexer.Token {
	return n.token
}
