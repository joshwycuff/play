package model

import "github.com/joshwycuff/play/util"

type HistoryData struct {
	Command    string
	Stdin      *string
	Stdout     *string
	Stderr     *string
	ExitStatus int
}

func NewHistoryData() *HistoryData {
	data := HistoryData{
		Command:    "",
		Stdin:      util.P(""),
		Stdout:     util.P(""),
		Stderr:     util.P(""),
		ExitStatus: 0,
	}
	return &data
}

type HistoryNode = Node[HistoryData]

type History struct {
	root    *HistoryNode
	current *HistoryNode
	nodes   []*HistoryNode
}

func NewHistory(rootCommand string, rootStdout *string) History {
	rootEntry := HistoryData{Command: rootCommand, Stdin: util.P(""), Stdout: rootStdout, Stderr: util.P("")}
	rootNode := NewRootNode(&rootEntry)
	childEntry := HistoryData{Command: "", Stdin: rootStdout, Stdout: util.P(""), Stderr: util.P("")}
	childNode := rootNode.AddChild(&childEntry)
	return History{
		root:    &rootNode,
		current: &childNode,
		nodes:   []*HistoryNode{&rootNode, &childNode},
	}
}

func (h *History) GetCurrent() *HistoryNode {
	return h.current
}

func (h *History) SetCurrent(node *HistoryNode) {
	h.current = node
}

func (h *History) AddChild(entry *HistoryData) *HistoryNode {
	child := h.current.AddChild(entry)
	h.nodes = append(h.nodes, &child)
	return &child
}

func (h *History) AddSibling(entry *HistoryData) *HistoryNode {
	sibling := h.current.Parent.AddChild(entry)
	h.nodes = append(h.nodes, &sibling)
	return &sibling
}

func (h *History) GetParent() *HistoryNode {
	if h.current.Equals(h.root) {
		return h.current
	}

	return h.current.Parent
}

func (h *History) GetLatestChild() *HistoryNode {
	if len(h.current.Children) == 0 {
		return h.current
	}

	return h.current.Children[len(h.current.Children)-1]
}

func (h *History) GetPrevSibling() *HistoryNode {
	if h.current.Equals(h.root) {
		return h.root
	}

	parent := h.current.Parent
	index := findIndexOfHistoryNode(parent.Children, h.current)

	if index == 0 {
		return h.current
	}

	return parent.Children[index-1]
}

func (h *History) GetNextSibling() *HistoryNode {
	if h.current.Equals(h.root) {
		return h.root
	}

	parent := h.current.Parent
	index := findIndexOfHistoryNode(parent.Children, h.current)

	if index == len(parent.Children)-1 {
		return h.current
	}

	return parent.Children[index+1]
}

func findIndexOfHistoryNode(nodes []*HistoryNode, node *HistoryNode) int {
	for i, n := range nodes {
		if n.Equals(node) {
			return i
		}
	}
	return -1
}
