package config

import (
	"fmt"
)

type Config struct {
	Nodes       []Node       `yaml:"nodes"`
	Pods        []Pod        `yaml:"pods"`
	Switches    []Switch     `yaml:"switches" mapstructure:"switches"`
	PodsConfigs []PodConfigs `yaml:"podsconfigs"`
	Tests       []Test       `yaml:"tests"`
}

// Node
type Node struct {
	Name        string      `yaml:"name" mapstructure:"name"`
	Type        string      `yaml:"type" mapstructure:"type"`
	NetBase     string      `yaml:"net_base" mapstructure:"net_base"`
	VolumeBase  string      `yaml:"volume" mapstructure:"volume"`
	Image       string      `yaml:"image" mapstructure:"image"`
	BuildFile   string      `yaml:"buildfile" mapstructure:"buildfile"`
	Interfaces  []Interface `yaml:"interfaces" mapstructure:"interfaces"`
	Sysctls     []string    `yaml:"sysctls" mapstructure:"sysctls"`
	Mounts      []string    `yaml:"mounts,flow" mapstructure:"mounts,flow"`
	DNS         []string    `yaml:"dns,flow" mapstructure:"dns,flow"`
	DNSSearches []string    `yaml:"dns_search,flow" mapstructure:"dns_search,flow"`
	HostName    string      `yaml:"hostname" mapstructure:"hostname"`
	EntryPoint  string      `yaml:"entrypoint" mapstructure:"entrypoint"`
	ExtraArgs   string      `yaml:"docker_run_extra_args" mapstructure:"docker_run_extra_args"`
}

// Interface
type Interface struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Args string `yaml:"args"`
	Addr string `yaml:"addr"`
}

// Pod
type Pod struct {
	Name        string   `yaml:"name"`
	NodeName    string   `yaml:"node" mapstructure:"node"`
	Image       string   `yaml:"image"`
	VolumeBase  string   `yaml:"volume" mapstructure:"volume"`
	BuildFile   string   `yaml:"buildfile" mapstructure:"buildfile"`
	Sysctls     []string `yaml:"sysctls" mapstructure:"sysctls"`
	Mounts      []string `yaml:"mounts,flow" mapstructure:"mounts,flow"`
	DNS         []string `yaml:"dns,flow" mapstructure:"dns,flow"`
	DNSSearches []string `yaml:"dns_search,flow" mapstructure:"dns_search,flow"`
	EntryPoint  string   `yaml:"entrypoint" mapstructure:"entrypoint"`
	ExtraArgs   string   `yaml:"docker_run_extra_args" mapstructure:"docker_run_extra_args"`
}

// Switch
type Switch struct {
	Name       string      `yaml:"name"`
	Interfaces []Interface `yaml:"interfaces" mapstructure:"interfaces"`
}

// PodConfigs
type PodConfigs struct {
	Name string `yaml:"name"`
	Cmds string `yaml:"cmds"`
}

// Test
type Test struct {
	Name string `yaml:"name" mapstructure:"name"`
	Cmds string `yaml:"cmds" mapstructure:"cmds"`
}

// LinkStatus
type LinkStatus struct {
	LeftNodeName  string
	LeftInfName   string
	LeftIsSet     bool
	LeftNodeType  string
	RightNodeName string
	RightInfName  string
	RightIsSet    bool
	RightNodeType string
}

// GetContainerPid func is Output get Docker PID Command
func GetContainerPid(nodename string) (getpidCmd string) {

	getpidCmd = fmt.Sprintf("PID=`docker inspect %s --format '{{.State.Pid}}'`", nodename)

	return getpidCmd
}

// HostLinkUp Link up link of host
func HostLinkUp(linkName string) (linkUpCmd string) {

	linkUpCmd = fmt.Sprintf("ip link set %s up", linkName)

	return linkUpCmd
}

// NetnsLinkUp Link up link of netns
func NetnsLinkUp(netnsName string, linkName string) (netnsLinkUpCmd string) {

	netnsLinkUpCmd = fmt.Sprintf("ip netns exec %s ip link set %s up", netnsName, linkName)

	return netnsLinkUpCmd
}
