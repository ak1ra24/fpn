package config

import (
	"fmt"
	"strings"
)

func (podconfig *PodConfigs) ParseCmd() (podCmds []string) {
	for _, s := range strings.Split(podconfig.Cmds, "\n") {
		if s != "" {
			podCmd := fmt.Sprintf("docker exec -it %s %s", podconfig.Name, s)
			podCmds = append(podCmds, podCmd)
		}
	}

	return podCmds
}
