package config

import (
	"fmt"
	"strings"
)

// TestCmdExec
func (t *Test) TestCmdExec() (TestCmds []string) {

	for _, s := range strings.Split(t.Cmds, "\n") {
		if s != "" {
			TestCmd := fmt.Sprintf("docker exec -it %s %s", t.Name, s)
			TestCmds = append(TestCmds, TestCmd)
		}
	}

	return TestCmds
}
