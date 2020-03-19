package config

import "fmt"

// CreatePod
func (pod *Pod) CreatePod() (createPodCmd string) {

	createPodCmd = fmt.Sprintf("docker run -td --net=container:%s --name %s --rm --privileged ", pod.NodeName, pod.Name)

	if len(pod.Sysctls) != 0 {
		for _, sysctl := range pod.Sysctls {
			createPodCmd += fmt.Sprintf("--sysctl %s ", sysctl)
		}
	}

	if pod.EntryPoint != "" {
		createPodCmd += fmt.Sprintf("--entrypoint %s ", pod.EntryPoint)
	}

	if pod.VolumeBase == "" {
		createPodCmd += "-v /tmp/tinet:/tinet "
	} else {
		createPodCmd += fmt.Sprintf("-v %s:/tinet ", pod.VolumeBase)
	}

	if len(pod.Mounts) != 0 {
		for _, mount := range pod.Mounts {
			createPodCmd += fmt.Sprintf("-v %s ", mount)
		}
	}

	if len(pod.DNS) != 0 {
		for _, dns := range pod.DNS {
			createPodCmd += fmt.Sprintf("--dns=%s ", dns)
		}
	}

	if len(pod.DNSSearches) != 0 {
		for _, dns_search := range pod.DNSSearches {
			createPodCmd += fmt.Sprintf("--dns-search=%s ", dns_search)
		}
	}

	if pod.ExtraArgs != "" {
		createPodCmd += fmt.Sprintf("%s ", pod.ExtraArgs)
	}

	createPodCmd += pod.Image

	return createPodCmd
}

// DeleteNode Delete docker and netns
func (pod *Pod) DeletePod() (deletePodCmds []string) {

	deletePodCmd := fmt.Sprintf("docker rm -f %s", pod.Name)

	deleteNsCmd := fmt.Sprintf("rm -rf /var/run/netns/%s", pod.Name)

	deletePodCmds = []string{deletePodCmd, deleteNsCmd}

	return deletePodCmds
}
