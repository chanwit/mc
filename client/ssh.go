package client

import (
	"errors"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/chanwit/gattai/machine"
	"github.com/chanwit/gattai/utils"
	Cli "github.com/docker/docker/cli"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
)

func DoSsh(cli interface{}, args ...string) error {

	ssh.SetDefaultClient(ssh.External)

	cmd := Cli.Subcmd("ssh",
		[]string{"MACHINES COMMAND"},
		"Run SSH commands on machines specified. Use - to run SSH on the active host.", false)

	cmd.ParseFlags(args, true)

	store := machine.GetDefaultStore(utils.GetBaseDir())

	p, err := machine.ReadProvision("provision.yml")
	if err != nil {
		log.Debugf("err: %s", err)
		return err
	}

	pattern := cmd.Args()[0]

	if pattern == "" {
		log.Fatal("Error: Please specify a machine name or pattern.")
	}

	// TODO if ssh -all
	machineList := p.GetMachineList(pattern)

	if len(machineList) == 1 && len(cmd.Args()) == 1 {
		host, err := store.Load(machineList[0])
		if err != nil {
			log.Fatal(err)
		}

		currentState, err := host.Driver.GetState()
		if err != nil {
			log.Fatal(err)
		}

		if currentState != state.Running {
			log.Fatalf("Error: Cannot run SSH command: Host %q is not running", host.Name)
		}

		err = host.CreateSSHShell()
		if err != nil {
			return err
		}
	} else {

		sshCmd := strings.Join(cmd.Args()[1:], " ")
		if strings.TrimSpace(sshCmd) == "" {
			return errors.New("Interative shell is not allowed for multiple hosts.")
		}

		// TODO should limit string channel
		outputs := make(chan string, len(machineList))
		for _, name := range machineList {
			go func(name string) {
				host, err := store.Load(name)
				if err != nil {
					log.Fatal(err)
				}

				currentState, err := host.Driver.GetState()
				if err != nil {
					log.Fatal(err)
				}

				if currentState != state.Running {
					log.Fatalf("Error: Cannot run SSH command: Host %q is not running", host.Name)
				}

				output, err := host.RunSSHCommand(sshCmd)
				if err != nil {
					if len(machineList) == 1 {
						outputs <- err.Error()
					} else {
						outputs <- fmt.Sprintf("\n%s:\n%s", name, err.Error())
					}
				} else {
					if len(machineList) == 1 {
						outputs <- string(output)
					} else {
						outputs <- fmt.Sprintf("\n%s:\n%s", name, string(output))
					}
				}
			}(name)
		}

		for i := 0; i < len(machineList); i++ {
			fmt.Print(<-outputs)
		}
	}

	return nil
}