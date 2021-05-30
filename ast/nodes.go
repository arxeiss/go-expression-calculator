package ast

import (
	"github.com/m1gwings/treedrawer/tree"
)

func ToTreeDrawer(rootNode Node) *tree.Tree {
	t := tree.NewTree(nil)
	rootNode.toTreeDrawer(t)
	return t
}

type Node interface {
	toTreeDrawer(*tree.Tree)
}

// Just make sure all node types implements Node interface
var _ Node = NumericNode(123)
var _ Node = VariableNode("varName")
var _ Node = &UnaryNode{}
var _ Node = &BinaryNode{}
var _ Node = &FunctionNode{}

type NumericNode float64

func (n NumericNode) toTreeDrawer(t *tree.Tree) {
	t.SetVal(tree.NodeFloat64(n))
}

type VariableNode string

func (n VariableNode) toTreeDrawer(t *tree.Tree) {
	t.SetVal(tree.NodeString(n))
}

type UnaryNode struct {
	next     Node
	operator Operation
}

func NewUnaryNode(operator Operation, next Node) *UnaryNode {
	return &UnaryNode{
		operator: operator,
		next:     next,
	}
}

func (n *UnaryNode) Operator() Operation {
	return n.operator
}
func (n *UnaryNode) Next() Node {
	return n.next
}
func (n UnaryNode) toTreeDrawer(t *tree.Tree) {
	t.SetVal(tree.NodeString(n.operator.String()))
	n.next.toTreeDrawer(t.AddChild(nil))
}

type BinaryNode struct {
	operator Operation
	left     Node
	right    Node
}

func NewBinaryNode(operator Operation, left, right Node) *BinaryNode {
	return &BinaryNode{
		operator: operator,
		left:     left,
		right:    right,
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

type FunctionNode struct {
	name  string
	param Node
}

func NewFunctionNode(name string, param Node) *FunctionNode {
	return &FunctionNode{
		name:  name,
		param: param,
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
