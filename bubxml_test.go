package main

import (
	"fmt"
	"strings"
	"testing"
)

type NetConf struct {
}

func (n *Netconf) Push(node *Container) error {
	fmt.Println("accept", node.Name)
}

func TestHelloWorld(t *testing.T) {
	r := strings.NewReader(getxml())
	a := NewAdapter(&NetConf{})
	a.DeCoder(r)

}

func getxml() string {
	return `
<?xml version="1.0" encoding="UTF-8"?> 
<rpc message-id="101"  xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
  <get>
   <filter type="subtree">
     <ifm xmlns="http://www.huawei.com/netconf/vrp" content-version="1.0" format-version="1.0">
       <interfaces>
        <interface>
          <ifName>10GE1/0/1</ifName>
        </interface>
       </interfaces>
    </ifm>
	 <ifm xmlns="http://www.huawei.com/netconf/vrp" content-version="1.0" format-version="1.0">
         <interfaces>
          <interface>
            <ifName>20GE1/0/2</ifName>
          </interface>
         </interfaces>
	 </ifm>
  </filter>
 </get>
</rpc>
	`
}
