package config

import "fmt"

// CreateSwitch Create bridge set in config
func (bridge *Switch) CreateSwitch() (createSwitchCmds []string) {

	addSwitchCmd := fmt.Sprintf("ovs-vsctl add-br %s", bridge.Name)
	createSwitchCmds = append(createSwitchCmds, addSwitchCmd)

	bridgeUpCmd := HostLinkUp(bridge.Name)
	createSwitchCmds = append(createSwitchCmds, bridgeUpCmd)

	return createSwitchCmds
}

// DeleteSwitch Delete bridge
func (sw *Switch) DeleteSwitch() (deleteSwCmd string) {

	deleteSwCmd = fmt.Sprintf("ovs-vsctl del-br %s", sw.Name)

	return deleteSwCmd
}
