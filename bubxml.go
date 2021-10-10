package bubxml

import (
	"container/list"
	"encoding/xml"
	"io"
	"log"
)

type Container struct {
	Name       string
	Attributes map[string]string
	Leafs      []Leaf
}

type Leaf struct {
	Name, Text string
}

type Revicer interface {
	Push(node Container) error
}

type tnode struct {
	name, value string
	attrs       map[string]string
}

type Adapter struct {
	depth, breadth int
	revicer        Revicer
	stack          *list.List
	leafMap        map[int][]Leaf
	body           chan Container
	done           chan string
}

func NewAdapter(revi Revicer) *Adapter {
	return &Adapter{depth: -1, breadth: -1, revicer: revi}
}

func (a *Adapter) init() {
	a.depth = 0
	a.breadth = 0
	a.stack = list.New()
	a.leafMap = make(map[int][]Leaf, 0)
	a.body = make(chan Container)
	a.done = make(chan string)
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

func (a *Adapter) breadthReset() {
	a.breadth = 1
}

func (a *Adapter) isLeaf() bool {
	if a.breadth == 2 {
		a.breadth = 3
	}
	return a.breadth == 3
}

func (a *Adapter) isBreadthReseted() bool {
	return a.breadth == 1
}

func (a *Adapter) pushNode(t *tnode) {
	a.stack.PushBack(t)
}

func (a *Adapter) getTopNode() *tnode {
	last := a.stack.Back()
	node, _ := last.Value.(*tnode)
	return node
}

func (a *Adapter) popNode() *tnode {
	last := a.stack.Back()
	node, ok := last.Value.(*tnode)
	if ok {
		a.stack.Remove(last)
		return node
	}
	return nil
}

func (a *Adapter) DeCoder(r io.Reader) {
	a.init()
	go a.doParse(r)
	defer a.finally()
outer:
	for {
		select {
		case msg := <-a.body:
			a.revicer.Push(msg)
		case <-a.done:
			break outer
		}
	}
}

func (a *Adapter) finally() {
	a.depth = -1
	a.breadth = -1
	a.stack = nil
	close(a.body)
	close(a.done)
}

func (a *Adapter) doParse(r io.Reader) {
	decoder := xml.NewDecoder(r)
	for tk, err := decoder.Token(); err == nil; tk, err = decoder.Token() {
		switch token := tk.(type) {
		case xml.StartElement:
			a.prefix(token)
		case xml.CharData:
			a.infix(token)
		case xml.EndElement:
			a.subfix(token)
		default:
			log.Println("No support right now.")

		}
	}
	a.done <- "done"
}

func (a *Adapter) prefix(token xml.StartElement) {
	a.depthPlus()
	a.breadthReset()
	a.pushNode(&tnode{name: token.Name.Local, attrs: getAttributes(token.Attr)})
}

func (a *Adapter) infix(token xml.CharData) {
	if a.isBreadthReseted() {
		a.breadthPlus()
		node := a.getTopNode()
		node.value = string([]byte(token))
	}
}

func (a *Adapter) subfix(token xml.EndElement) {
	a.breadthPlus()
	if a.isLeaf() {
		a.depthMinus()
		a.appendLeaf(a.popNode())
	} else {
		n := a.popNode()
		c := Container{Name: n.name, Attributes: n.attrs, Leafs: a.leafMap[a.depth]}
		a.body <- c
		delete(a.leafMap, a.depth)
		a.depthMinus()
	}

}

func (a *Adapter) appendLeaf(n *tnode) {
	leaf := Leaf{Name: n.name, Text: n.value}
	if l, ok := a.leafMap[a.depth]; ok {
		l = append(l, leaf)
		a.leafMap[a.depth] = l
	} else {
		ls := make([]Leaf, 0)
		ls = append(ls, leaf)
		a.leafMap[a.depth] = ls
	}
}

func getAttributes(attr []xml.Attr) map[string]string {
	attrs := make(map[string]string, 0)
	for _, a := range attr {
		attrs[a.Name.Local] = a.Value
	}
	return attrs
}
