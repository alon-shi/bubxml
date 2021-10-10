package main

import (
	"bubxml"
	"bufio"
	"fmt"
	"os"
)

type NetConf struct {
}

func (n *NetConf) Push(node bubxml.Container) error {
	fmt.Println("accept", node.Name, node.Attributes)
	for _, n := range node.Leafs {
		fmt.Println(n.Name, n.Text)
	}
	return nil
}

func main() {
	f, err := os.Open("./netconf.xml")
	if err == nil {
		defer f.Close()
		r := bufio.NewReader(f)
		a := bubxml.NewAdapter(&NetConf{})
		a.DeCoder(r)
	}

}
