package bubxml

import (
	"container/list"
	"encoding/xml"
	"io"
	"log"
)

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
	attrs       map[string]string
}

type Adapter struct {
	depth, breadth int
	revicer        Revicer
	stack          list.List
	leafMap        map[string]string
	body           chan Branch
	done           chan string
}

func NewAdapter(revi Revicer) {
	return &Adapter{depth: -1, breadth: -1, revicer: revi}
}

func (a *Adapter) init() {
	a.depth = 0
	a.breadth = 0
	a.stack = list.New()
	a.body = make(chan Branch)
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

func (a *Adapter) DeCode(r io.Reader) {
	go doParse(r)
	defer a.finally()
outer:
	for {
		select {
		case msg <- a.body:
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
	for tk, err = decoder.Token(); err == nil; tk, err = decoder.Token() {
		switch token := t.(type) {
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
}

func (a *Adapter) prefix(token xml.StartElement) {

}

func (a *Adapter) infix(token xml.CharData) {

}

func (a *Adapter) subfix(token xml.EndElement) {
}
