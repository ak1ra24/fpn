package config

import (
	"fmt"
)

// CreateNode
func (node *Node) CreateNode() (createNodeCmd string) {

	if node.NetBase == "" {
		node.NetBase = "none"
	}

	if node.Type == "docker" || node.Type == "" {
		createNodeCmd = fmt.Sprintf("docker run -td --net %s --name %s --rm --privileged ", node.NetBase, node.Name)

		if len(node.HostName) != 0 {
			createNodeCmd += fmt.Sprintf("--hostname %s ", node.HostName)
		} else {
			createNodeCmd += fmt.Sprintf("--hostname %s ", node.Name)
		}

		if len(node.Sysctls) != 0 {
			for _, sysctl := range node.Sysctls {
				createNodeCmd += fmt.Sprintf("--sysctl %s ", sysctl)
			}
		}

		if node.EntryPoint != "" {
			createNodeCmd += fmt.Sprintf("--entrypoint %s ", node.EntryPoint)
		}

		if node.VolumeBase == "" {
			createNodeCmd += "-v /tmp/tinet:/tinet "
		} else {
			createNodeCmd += fmt.Sprintf("-v %s:/tinet ", node.VolumeBase)
		}

		if len(node.Mounts) != 0 {
			for _, mount := range node.Mounts {
				createNodeCmd += fmt.Sprintf("-v %s ", mount)
			}
		}

		if len(node.DNS) != 0 {
			for _, dns := range node.DNS {
				createNodeCmd += fmt.Sprintf("--dns=%s ", dns)
			}
		}

		if len(node.DNSSearches) != 0 {
			for _, dns_search := range node.DNSSearches {
				createNodeCmd += fmt.Sprintf("--dns-search=%s ", dns_search)
			}
		}

		if node.ExtraArgs != "" {
			createNodeCmd += fmt.Sprintf("%s ", node.ExtraArgs)
		}

		createNodeCmd += "k8s.gcr.io/pause"
	}

	return createNodeCmd
}

// Mount_docker_netns Mount docker netns to ip netns<Paste>
func (node *Node) Mount_docker_netns() (mountDockerNetnsCmds []string) {

	netnsDir := "/var/run/netns"
	mkdirCmd := fmt.Sprintf("mkdir -p %s", netnsDir)
	mountDockerNetnsCmds = append(mountDockerNetnsCmds, mkdirCmd)
	dockerPid := GetContainerPid(node.Name)
	mountDockerNetnsCmds = append(mountDockerNetnsCmds, dockerPid)
	mountDockerNetnsCmd := fmt.Sprintf("ln -s /proc/$PID/ns/net /var/run/netns/%s", node.Name)
	mountDockerNetnsCmds = append(mountDockerNetnsCmds, mountDockerNetnsCmd)

	return mountDockerNetnsCmds
}

// DelNsCmds
func (node *Node) DelNsCmd() (delNsCmd string) {

	delNsCmd = fmt.Sprintf("ip netns del %s", node.Name)

	return delNsCmd
}

// DeleteNode Delete docker and netns
func (node *Node) DeleteNode() (deleteNodeCmds []string) {

	var deleteCmd string

	if node.Type == "docker" || node.Type == "" {
		deleteCmd = fmt.Sprintf("docker rm -f %s", node.Name)
	} else if node.Type == "netns" {
		deleteCmd = fmt.Sprintf("ip netns del %s", node.Name)
	} else {
		return []string{""}
	}

	deleteNsCmd := fmt.Sprintf("rm -rf /var/run/netns/%s", node.Name)

	deleteNodeCmds = []string{deleteCmd, deleteNsCmd}

	return deleteNodeCmds
}
