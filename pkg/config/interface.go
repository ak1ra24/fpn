package config

import (
	"fmt"
	"strings"
)

// N2nLink Connect links between nodes
func (inf *Interface) N2nLink(nodeName string) (n2nLinkCmds []string) {

	nodeinf := inf.Name
	peerinfo := strings.Split(inf.Args, "#")
	peerNode := peerinfo[0]
	peerinf := peerinfo[1]
	n2nlinkCmd := fmt.Sprintf("ip link add %s netns %s type veth peer name %s netns %s", nodeinf, nodeName, peerinf, peerNode)
	n2nLinkCmds = append(n2nLinkCmds, n2nlinkCmd)
	n2nLinkCmds = append(n2nLinkCmds, NetnsLinkUp(nodeName, nodeinf))
	n2nLinkCmds = append(n2nLinkCmds, NetnsLinkUp(peerNode, peerinf))

	if len(inf.Addr) != 0 {
		addrSetCmd := fmt.Sprintf("ip netns exec %s ip link set %s address %s", nodeName, inf.Name, inf.Addr)
		n2nLinkCmds = append(n2nLinkCmds, addrSetCmd)
	}

	return n2nLinkCmds
}

// S2nLink Connect links between nodes and switches
func (inf *Interface) S2nLink(nodeName string) (s2nLinkCmds []string) {

	nodeinf := inf.Name
	peerBr := inf.Args
	peerBrInf := fmt.Sprintf("%s-%s", peerBr, nodeName)
	s2nLinkCmd := fmt.Sprintf("ip link add %s netns %s type veth peer name %s", nodeinf, nodeName, peerBrInf)
	s2nLinkCmds = append(s2nLinkCmds, s2nLinkCmd)
	s2nLinkCmds = append(s2nLinkCmds, NetnsLinkUp(nodeName, nodeinf))
	s2nLinkCmds = append(s2nLinkCmds, HostLinkUp(peerBrInf))
	setBrLinkCmd := fmt.Sprintf("ovs-vsctl add-port %s %s", peerBr, peerBrInf)
	s2nLinkCmds = append(s2nLinkCmds, setBrLinkCmd)

	return s2nLinkCmds
}

// V2cLink Connect links between veth and container
func (inf *Interface) V2cLink(nodeName string) (v2cLinkCmds []string) {

	nodeinf := inf.Name
	peerName := inf.Args
	v2cLinkCmd := fmt.Sprintf("ip link add %s type veth peer name %s", nodeinf, peerName)
	v2cLinkCmds = append(v2cLinkCmds, v2cLinkCmd)

	v2cLinkCmds = append(v2cLinkCmds, inf.P2cLink(nodeName)...)
	v2cLinkCmds = append(v2cLinkCmds, HostLinkUp(peerName))

	return v2cLinkCmds
}

// P2cLink Connect links between phys-eth and container
func (inf *Interface) P2cLink(nodeName string) (p2cLinkCmds []string) {

	physInf := inf.Name
	setNsCmd := fmt.Sprintf("ip link set dev %s netns %s", physInf, nodeName)
	p2cLinkCmds = append(p2cLinkCmds, setNsCmd)
	execNsCmd := fmt.Sprintf("ip netns exec %s ip link set %s up", nodeName, physInf)
	p2cLinkCmds = append(p2cLinkCmds, execNsCmd)
	delNsCmd := fmt.Sprintf("ip netns del %s", nodeName)
	p2cLinkCmds = append(p2cLinkCmds, delNsCmd)

	return p2cLinkCmds
}
