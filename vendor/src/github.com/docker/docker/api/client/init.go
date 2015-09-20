package client

import (
	"fmt"
	"os"
	"io/ioutil"
)

func (cli *DockerCli) CmdInit(args ...string) error {
	info, err := os.Stat(".gattai")
	if err == nil && info.IsDir() {
		return fmt.Errorf(".gattai is already existed")
	}

	err = os.Mkdir(".gattai", 0644)
	if err != nil {
		return err
	}

	provisionYml := `---
machines:
#
# Put a group of your machines here.
# For example,
#
#  ocean:
#    driver: digitalocean
#    instances: 5
#

`
	_, err = os.Stat("provision.yml")
	// err != nil is file not found, OK to write
	if err != nil {
		err = ioutil.WriteFile("provision.yml",
			[]byte(provisionYml), 0644)
		if err != nil {
			return err
		}
	}

	compositionYml := `---
#
# Composition is simply a docker-compose.yml file.
#
`
	_, err = os.Stat("composition.yml")
	// err != nil is file not found, OK to write
	if err != nil {
		err = ioutil.WriteFile("composition.yml",
			[]byte(compositionYml), 0644)
		if err != nil {
			return err
		}
	}

	fmt.Println("Gattai mission repository is initialized.")

	return nil
}