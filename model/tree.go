package model

type NodeId = int

var currentId NodeId = -1

type Node[T any] struct {
	Id       NodeId
	Parent   *Node[T]
	Children []*Node[T]
	Data     *T
}

func NewNode[T any](parent *Node[T], data *T) Node[T] {
	currentId += 1
	return Node[T]{Id: currentId, Parent: parent, Children: []*Node[T]{}, Data: data}
}

func NewRootNode[T any](data *T) Node[T] {
	return NewNode(nil, data)
}

func NewChildNode[T any](parent *Node[T], data *T) Node[T] {
	return NewNode(parent, data)
}

func (n *Node[T]) AddChild(data *T) Node[T] {
	child := NewNode(n, data)
	n.Children = append(n.Children, &child)
	return child
}

func (n *Node[T]) Equals(other *Node[T]) bool {
	return n.Id == other.Id
}
