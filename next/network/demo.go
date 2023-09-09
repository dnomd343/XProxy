package network

import (
	"fmt"
	"github.com/coreos/go-iptables/iptables"
)

type ipTables struct {
	v4 *iptables.IPTables
	v6 *iptables.IPTables
}

var tables *ipTables

func init() {
	timeout := iptables.Timeout(8)
	it4, err := iptables.New(iptables.IPFamily(iptables.ProtocolIPv4), timeout)
	if err != nil {
		// TODO: panic here
		fmt.Printf("failed to init iptables -> %v\n", err)
	}
	it6, err := iptables.New(iptables.IPFamily(iptables.ProtocolIPv6), timeout)
	if err != nil {
		fmt.Printf("failed to init ip6tables -> %v\n", err)
	}

	tables = &ipTables{
		v4: it4,
		v6: it6,
	}
}

func Demo() {
	fmt.Println("iptables demo start")

	//it, err := iptables.New(iptables.IPFamily(iptables.ProtocolIPv4), iptables.Timeout(5))
	//it, err := iptables.New(iptables.IPFamily(iptables.ProtocolIPv6), iptables.Timeout(5))

	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(it)

	fmt.Println(tables.v4)
	fmt.Println(tables.v6)

	chains, _ := tables.v4.ListChains("filter")
	fmt.Println(chains)

	rules, _ := tables.v4.List("filter", "DOCKER-ISOLATION-STAGE-2")
	//fmt.Println(rules)
	for _, rule := range rules {
		fmt.Println(rule)

	}
}
