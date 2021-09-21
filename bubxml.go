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

func (a *Adapter) init() {
	a.depth = 0
	a.breadth = 0
}

func (a *Adapter) depthPlus() {
	a.depth++
}

func (a *Adapter) depthMinus() {
	a.depth--
}

func (a *Adapter) breadthPlus() {
	a.breadth++
}

func (a *Adapter) breadthMinus() {
	a.breadth--
}

func (a *Adapter) DoParse(r io.Reader) {
}
