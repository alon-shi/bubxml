package bubxml

import "io"

type Branch struct {
	Name       string
	Attributes map[string]string
	Leafs      []Leaf
}

type Leaf struct {
	Name, Text string
}

type Revicer interface {
	Push(node *Branch) error
}

type tnode struct {
	name, value string
}

type Adapter struct {
	depth, breadth int
	revicer        Revicer
}

func NewAdapter(revi Revicer) {
	return &Adapter{depth: -1, breadth: -1, revicer: revi}
}

func doParse(r io.Reader) {
}
