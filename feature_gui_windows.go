// Build for Windows
package main

import (
    "github.com/fsouza/go-dockerclient"
    "errors"
)

func enableGui(project Project) (map[docker.Port][]docker.PortBinding, []string, error) {
    portBindings := map[docker.Port][]docker.PortBinding{}
    envVariables := []string{}

    return portBindings, envVariables, errors.New("Could not enable GUI: it has not been implemented for linux yet.")
}
