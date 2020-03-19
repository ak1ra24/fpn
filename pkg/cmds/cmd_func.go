package cmds

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ak1ra24/fpn/pkg/config"
	"github.com/ak1ra24/fpn/pkg/utils"

	"github.com/k0kubun/pp"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

func LoadCfg(c *cli.Context) (config config.Config, err error) {
	cfgFile := c.String("config")
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		viper.SetConfigType("yaml")

		if err = viper.ReadInConfig(); err != nil {
			return config, err
		}

		if err = viper.Unmarshal(&config); err != nil {
			return config, err
		}
	} else {
		err = fmt.Errorf("not set config file.")
		return config, err
	}

	return config, nil

}

func CmdPrint(c *cli.Context) error {
	config, err := LoadCfg(c)
	if err != nil {
		return err
	}

	pp.Print(config)

	return nil
}

func CmdCheck(c *cli.Context) error {

	fpnconfig, err := LoadCfg(c)
	if err != nil {
		return err
	}

	nodes := fpnconfig.Nodes
	swiches := fpnconfig.Switches
	confmap := map[string]string{}

	for _, node := range nodes {
		for _, inf := range node.Interfaces {
			if inf.Type == "direct" {
				host := node.Name + ":" + inf.Name
				peer := strings.Split(inf.Args, "#")
				target := peer[0] + ":" + peer[1]
				confmap[host] = target
			} else if inf.Type == "bridge" {
				host := node.Name + ":" + inf.Name
				target := inf.Args + ":" + node.Name
				confmap[host] = target
			}
		}
	}

	for _, sw := range swiches {
		for _, inf := range sw.Interfaces {
			host := sw.Name + ":" + inf.Args
			target := inf.Args + ":" + inf.Name
			confmap[host] = target
		}
	}

	var matchNum int
	falseConfigMap := map[string]string{}

	for key, value := range confmap {
		if confmap[key] == value && confmap[value] == key {
			matchNum++
		} else {
			falseConfigMap[key] = value
		}
	}

	if len(confmap) == matchNum {
		return nil
	} else {
		var errMsg string
		for key, value := range falseConfigMap {
			errMsg += fmt.Sprintf("%s<->%s\n", key, value)
		}
		return fmt.Errorf(errMsg)
	}
}

func CmdUp(c *cli.Context) error {
	fpnconfig, err := LoadCfg(c)
	if err != nil {
		return err
	}

	for _, n := range fpnconfig.Nodes {
		cmd := n.CreateNode()
		fmt.Println(cmd)
		mountDockerNsCmds := n.Mount_docker_netns()
		for _, mountDockerNsCmd := range mountDockerNsCmds {
			fmt.Println(mountDockerNsCmd)
		}
	}

	if len(fpnconfig.Switches) != 0 {
		for _, sw := range fpnconfig.Switches {
			createSwCmds := sw.CreateSwitch()
			utils.PrintCmds(os.Stdout, createSwCmds)
		}
	}

	var links []config.LinkStatus

	for _, node := range fpnconfig.Nodes {
		for _, inf := range node.Interfaces {
			if inf.Type == "direct" {
				rNodeArgs := strings.Split(inf.Args, "#")
				rNodeName := rNodeArgs[0]
				rInfName := rNodeArgs[1]
				peerFound := false
				for _, link := range links {
					if !link.RightIsSet {
						nodecheck := link.LeftNodeName == rNodeName
						infcheck := link.LeftInfName == rInfName
						if nodecheck && infcheck {
							link.RightNodeName = node.Name
							link.RightInfName = inf.Name
							link.RightNodeType = node.Type
							link.RightIsSet = true
							peerFound = true
						}
					}
				}
				if !peerFound {
					link := config.LinkStatus{LeftNodeName: node.Name, LeftInfName: inf.Name, RightNodeName: rNodeName, RightInfName: rInfName}
					link.LeftNodeType = node.Type
					link.LeftIsSet = true
					links = append(links, link)
					n2nLinkCmds := inf.N2nLink(node.Name)
					utils.PrintCmds(os.Stdout, n2nLinkCmds)
				}
			} else if inf.Type == "bridge" {
				s2nLinkCmds := inf.S2nLink(node.Name)
				utils.PrintCmds(os.Stdout, s2nLinkCmds)
			} else if inf.Type == "veth" {
				v2cLinkCmds := inf.V2cLink(node.Name)
				utils.PrintCmds(os.Stdout, v2cLinkCmds)
			} else if inf.Type == "phys" {
				p2cLinkCmds := inf.P2cLink(node.Name)
				utils.PrintCmds(os.Stdout, p2cLinkCmds)
			} else {
				err := fmt.Errorf("not supported interface type: %s", inf.Type)
				log.Fatal(err)
			}
		}
	}

	// check
	err = CmdCheck(c)
	if err != nil {
		return err
	}

	for _, node := range fpnconfig.Nodes {
		delNsCmd := node.DelNsCmd()
		utils.PrintCmd(os.Stdout, delNsCmd)
	}

	for _, p := range fpnconfig.Pods {
		cmd := p.CreatePod()
		fmt.Println(cmd)
	}

	return nil
}

func CmdConf(c *cli.Context) error {
	fpnconfig, err := LoadCfg(c)
	if err != nil {
		return err
	}

	for _, podConfig := range fpnconfig.PodsConfigs {
		podCmds := podConfig.ParseCmd()
		utils.PrintCmds(os.Stdout, podCmds)
	}

	return nil
}

func CmdUpConf(c *cli.Context) error {
	fpnconfig, err := LoadCfg(c)
	if err != nil {
		return err
	}

	// Up
	for _, n := range fpnconfig.Nodes {
		cmd := n.CreateNode()
		fmt.Println(cmd)
		mountDockerNsCmds := n.Mount_docker_netns()
		for _, mountDockerNsCmd := range mountDockerNsCmds {
			fmt.Println(mountDockerNsCmd)
		}
	}

	if len(fpnconfig.Switches) != 0 {
		for _, sw := range fpnconfig.Switches {
			createSwCmds := sw.CreateSwitch()
			utils.PrintCmds(os.Stdout, createSwCmds)
		}
	}

	var links []config.LinkStatus

	for _, node := range fpnconfig.Nodes {
		for _, inf := range node.Interfaces {
			if inf.Type == "direct" {
				rNodeArgs := strings.Split(inf.Args, "#")
				rNodeName := rNodeArgs[0]
				rInfName := rNodeArgs[1]
				peerFound := false
				for _, link := range links {
					if !link.RightIsSet {
						nodecheck := link.LeftNodeName == rNodeName
						infcheck := link.LeftInfName == rInfName
						if nodecheck && infcheck {
							link.RightNodeName = node.Name
							link.RightInfName = inf.Name
							link.RightNodeType = node.Type
							link.RightIsSet = true
							peerFound = true
						}
					}
				}
				if !peerFound {
					link := config.LinkStatus{LeftNodeName: node.Name, LeftInfName: inf.Name, RightNodeName: rNodeName, RightInfName: rInfName}
					link.LeftNodeType = node.Type
					link.LeftIsSet = true
					links = append(links, link)
					n2nLinkCmds := inf.N2nLink(node.Name)
					utils.PrintCmds(os.Stdout, n2nLinkCmds)
				}
			} else if inf.Type == "bridge" {
				s2nLinkCmds := inf.S2nLink(node.Name)
				utils.PrintCmds(os.Stdout, s2nLinkCmds)
			} else if inf.Type == "veth" {
				v2cLinkCmds := inf.V2cLink(node.Name)
				utils.PrintCmds(os.Stdout, v2cLinkCmds)
			} else if inf.Type == "phys" {
				p2cLinkCmds := inf.P2cLink(node.Name)
				utils.PrintCmds(os.Stdout, p2cLinkCmds)
			} else {
				err := fmt.Errorf("not supported interface type: %s", inf.Type)
				log.Fatal(err)
			}
		}
	}

	// check
	err = CmdCheck(c)
	if err != nil {
		return err
	}

	for _, node := range fpnconfig.Nodes {
		delNsCmd := node.DelNsCmd()
		utils.PrintCmd(os.Stdout, delNsCmd)
	}

	for _, p := range fpnconfig.Pods {
		cmd := p.CreatePod()
		fmt.Println(cmd)
	}

	// Conf
	for _, podConfig := range fpnconfig.PodsConfigs {
		podCmds := podConfig.ParseCmd()
		utils.PrintCmds(os.Stdout, podCmds)
	}

	return nil
}

func CmdDown(c *cli.Context) error {
	fpnconfig, err := LoadCfg(c)
	if err != nil {
		return err
	}

	for _, node := range fpnconfig.Nodes {
		delNodeCmds := node.DeleteNode()
		utils.PrintCmds(os.Stdout, delNodeCmds)
	}

	for _, pod := range fpnconfig.Pods {
		delPodCmds := pod.DeletePod()
		utils.PrintCmds(os.Stdout, delPodCmds)
	}

	return nil
}

func CmdReUp(c *cli.Context) error {

	fpnconfig, err := LoadCfg(c)
	if err != nil {
		return err
	}

	// Down
	for _, node := range fpnconfig.Nodes {
		delNodeCmds := node.DeleteNode()
		utils.PrintCmds(os.Stdout, delNodeCmds)
	}

	for _, pod := range fpnconfig.Pods {
		delPodCmds := pod.DeletePod()
		utils.PrintCmds(os.Stdout, delPodCmds)
	}

	// Up
	for _, n := range fpnconfig.Nodes {
		cmd := n.CreateNode()
		fmt.Println(cmd)
		mountDockerNsCmds := n.Mount_docker_netns()
		for _, mountDockerNsCmd := range mountDockerNsCmds {
			fmt.Println(mountDockerNsCmd)
		}
	}

	if len(fpnconfig.Switches) != 0 {
		for _, sw := range fpnconfig.Switches {
			createSwCmds := sw.CreateSwitch()
			utils.PrintCmds(os.Stdout, createSwCmds)
		}
	}

	var links []config.LinkStatus

	for _, node := range fpnconfig.Nodes {
		for _, inf := range node.Interfaces {
			if inf.Type == "direct" {
				rNodeArgs := strings.Split(inf.Args, "#")
				rNodeName := rNodeArgs[0]
				rInfName := rNodeArgs[1]
				peerFound := false
				for _, link := range links {
					if !link.RightIsSet {
						nodecheck := link.LeftNodeName == rNodeName
						infcheck := link.LeftInfName == rInfName
						if nodecheck && infcheck {
							link.RightNodeName = node.Name
							link.RightInfName = inf.Name
							link.RightNodeType = node.Type
							link.RightIsSet = true
							peerFound = true
						}
					}
				}
				if !peerFound {
					link := config.LinkStatus{LeftNodeName: node.Name, LeftInfName: inf.Name, RightNodeName: rNodeName, RightInfName: rInfName}
					link.LeftNodeType = node.Type
					link.LeftIsSet = true
					links = append(links, link)
					n2nLinkCmds := inf.N2nLink(node.Name)
					utils.PrintCmds(os.Stdout, n2nLinkCmds)
				}
			} else if inf.Type == "bridge" {
				s2nLinkCmds := inf.S2nLink(node.Name)
				utils.PrintCmds(os.Stdout, s2nLinkCmds)
			} else if inf.Type == "veth" {
				v2cLinkCmds := inf.V2cLink(node.Name)
				utils.PrintCmds(os.Stdout, v2cLinkCmds)
			} else if inf.Type == "phys" {
				p2cLinkCmds := inf.P2cLink(node.Name)
				utils.PrintCmds(os.Stdout, p2cLinkCmds)
			} else {
				err := fmt.Errorf("not supported interface type: %s", inf.Type)
				log.Fatal(err)
			}
		}
	}

	// check
	err = CmdCheck(c)
	if err != nil {
		return err
	}

	for _, node := range fpnconfig.Nodes {
		delNsCmd := node.DelNsCmd()
		utils.PrintCmd(os.Stdout, delNsCmd)
	}

	for _, p := range fpnconfig.Pods {
		cmd := p.CreatePod()
		fmt.Println(cmd)
	}

	return nil
}

func CmdReConf(c *cli.Context) error {
	fpnconfig, err := LoadCfg(c)
	if err != nil {
		return err
	}

	// Down
	for _, node := range fpnconfig.Nodes {
		delNodeCmds := node.DeleteNode()
		utils.PrintCmds(os.Stdout, delNodeCmds)
	}

	for _, pod := range fpnconfig.Pods {
		delPodCmds := pod.DeletePod()
		utils.PrintCmds(os.Stdout, delPodCmds)
	}

	// Up
	for _, n := range fpnconfig.Nodes {
		cmd := n.CreateNode()
		fmt.Println(cmd)
		mountDockerNsCmds := n.Mount_docker_netns()
		for _, mountDockerNsCmd := range mountDockerNsCmds {
			fmt.Println(mountDockerNsCmd)
		}
	}

	if len(fpnconfig.Switches) != 0 {
		for _, sw := range fpnconfig.Switches {
			createSwCmds := sw.CreateSwitch()
			utils.PrintCmds(os.Stdout, createSwCmds)
		}
	}

	var links []config.LinkStatus

	for _, node := range fpnconfig.Nodes {
		for _, inf := range node.Interfaces {
			if inf.Type == "direct" {
				rNodeArgs := strings.Split(inf.Args, "#")
				rNodeName := rNodeArgs[0]
				rInfName := rNodeArgs[1]
				peerFound := false
				for _, link := range links {
					if !link.RightIsSet {
						nodecheck := link.LeftNodeName == rNodeName
						infcheck := link.LeftInfName == rInfName
						if nodecheck && infcheck {
							link.RightNodeName = node.Name
							link.RightInfName = inf.Name
							link.RightNodeType = node.Type
							link.RightIsSet = true
							peerFound = true
						}
					}
				}
				if !peerFound {
					link := config.LinkStatus{LeftNodeName: node.Name, LeftInfName: inf.Name, RightNodeName: rNodeName, RightInfName: rInfName}
					link.LeftNodeType = node.Type
					link.LeftIsSet = true
					links = append(links, link)
					n2nLinkCmds := inf.N2nLink(node.Name)
					utils.PrintCmds(os.Stdout, n2nLinkCmds)
				}
			} else if inf.Type == "bridge" {
				s2nLinkCmds := inf.S2nLink(node.Name)
				utils.PrintCmds(os.Stdout, s2nLinkCmds)
			} else if inf.Type == "veth" {
				v2cLinkCmds := inf.V2cLink(node.Name)
				utils.PrintCmds(os.Stdout, v2cLinkCmds)
			} else if inf.Type == "phys" {
				p2cLinkCmds := inf.P2cLink(node.Name)
				utils.PrintCmds(os.Stdout, p2cLinkCmds)
			} else {
				err := fmt.Errorf("not supported interface type: %s", inf.Type)
				log.Fatal(err)
			}
		}
	}

	// check
	err = CmdCheck(c)
	if err != nil {
		return err
	}

	for _, node := range fpnconfig.Nodes {
		delNsCmd := node.DelNsCmd()
		utils.PrintCmd(os.Stdout, delNsCmd)
	}

	for _, p := range fpnconfig.Pods {
		cmd := p.CreatePod()
		fmt.Println(cmd)
	}

	// Conf
	for _, podConfig := range fpnconfig.PodsConfigs {
		podCmds := podConfig.ParseCmd()
		utils.PrintCmds(os.Stdout, podCmds)
	}

	return nil
}

func CmdTest(c *cli.Context) error {
	fpnconfig, err := LoadCfg(c)
	if err != nil {
		return err
	}

	testName := c.Args().Get(0)

	var tnTestCmds []string

	if testName == "all" || testName == "" {
		for _, test := range fpnconfig.Tests {
			tnTestCmds = test.TestCmdExec()
		}
	} else {
		for _, test := range fpnconfig.Tests {
			if testName == test.Name {
				tnTestCmds = test.TestCmdExec()
			}
		}
	}

	if len(tnTestCmds) == 0 {
		return fmt.Errorf("not found test name\n")
	}

	fmt.Fprintln(os.Stdout, strings.Join(tnTestCmds, "\n"))

	return nil
}
