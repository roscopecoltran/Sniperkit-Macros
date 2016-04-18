package main

import (
    "github.com/fsouza/go-dockerclient"
    "github.com/fgrehm/go-dockerpty"
    log "github.com/Sirupsen/logrus"
    "strings"
    // Persist "github.com/matthieudelaro/nut/persist"
)

// Start the container, execute the list of commands, and then stops the container.
// Inspired from https://github.com/fsouza/go-dockerclient/issues/287
// and from https://github.com/fsouza/go-dockerclient/issues/220#issuecomment-77777365
// DONE: solve key issue (Ctrl+C, q, ...). See if that can help: https://github.com/fgrehm/go-dockerpty : it does :)
// DONE: to solve tty size, use func (c *Client) ResizeContainerTTY(id string, height, width int) error : solved with dockerpty
// and/or func (c *Client) ResizeExecTTY(id string, height, width int) error
// TODO: report bug "StartExec Error: %s read /dev/stdin: bad file descriptor" when executing several commands : post issue on dockerpty
func execInContainer(commands []string, project Project) {
    log.Debug("commands (len = ", len(commands), ") : ", commands)
    var cmdConfig []string
    if len(commands) == 0 {
        log.Debug("Given list of commands is empty.")
        cmdConfig = []string{}
    } else if len(commands) == 1 {
        cmdConfig = []string{"bash", "-c", commands[0]}
    } else {
        cmdConfig = []string{"bash", "-c", strings.Join(commands, "; ")}
    }
    log.Debug("cmdConfig: ", cmdConfig)

    imageName, err := project.getBaseEnv().getDockerImageName()
    if err != nil {
        log.Error(err.Error())
        return
    }

    // endpoint := getDockerEndpoint()
    // client, err := docker.NewClient(endpoint) // TODO : fix https for boot2docker (with https://github.com/fsouza/go-dockerclient/issues/166)
    client, err := getDockerClient() // TODO : fix https for boot2docker (with https://github.com/fsouza/go-dockerclient/issues/166)
    if err != nil {
        // log.Error("Could not reach Docker host (", endpoint, "): ", err.Error())
        log.Error(err)
        return
    }
    log.Debug("Created client")

    //Pull image from Registry, if not present
    _, err = client.InspectImage(imageName)
    if err != nil {
        log.Error("Could not inspect image ", imageName, ": ", err.Error())

        log.Debug("Pulling image...")
        opts := docker.PullImageOptions{Repository: imageName}
        err = client.PullImage(opts, docker.AuthConfiguration{})
        if err != nil {
            log.Error("Could not pull image ", imageName, ": ", err.Error())
            // log.Debug(err.Error())
            return
        }
        log.Debug("Pulled image")
    }

    // prepare names of directories to mount
    // inspired from https://github.com/fsouza/go-dockerclient/issues/220#issuecomment-77777365
    mountingPoints := project.getMountingPoints()
    binds := make([]string, 0, len(mountingPoints))
    for _, directory := range(mountingPoints) {
        hostPath, hostPathErr := directory.fullHostPath()
        containerPath, containerPathErr := directory.fullContainerPath()
        if hostPathErr != nil {
            log.Error("Couldn't mount host directory: ", hostPathErr.Error())
            return
        }
        if containerPathErr != nil {
            log.Error("Couldn't container host directory: ", containerPathErr.Error())
            return
        }
        binds = append(binds, hostPath + ":" + containerPath)
    }
    log.Debug("binds", binds)

    portBindings := map[docker.Port][]docker.PortBinding{}
    envVariables := []string{}
    if project.getEnableGui() {
        portBindingsGUI, envVariablesGUI, bindsGUI, err := enableGui(project)
        if err != nil {
            log.Error("Could not enable GUI: ", err)
        } else {
            envVariables = append(envVariables, envVariablesGUI...)
            binds = append(binds, bindsGUI...)
            for k, v := range portBindingsGUI {
                portBindings[k] = v
            }
        }
    }

    //Try to create a container from the imageID
    config := docker.Config{
        Image: imageName,
        Cmd:          cmdConfig,
        OpenStdin:    true,
        StdinOnce:    true,
        AttachStdin:  true,
        AttachStdout: true,
        AttachStderr: true,
        Tty:          true,
        WorkingDir:   project.getWorkingDir(),
        Env:          envVariables,
    }
    // TODO : set following config options: https://godoc.org/github.com/fsouza/go-dockerclient#Config
    // User: set it to the user of the host, instead of root, to manage file permissions properly

    // TODO: ? Give the container a name? Can be done with docker.CreateContainerOptions{Name: "nut_myproject"}
    // opts2 := docker.CreateContainerOptions{Name: "nut_" + , Config: &config}
    opts2 := docker.CreateContainerOptions{Config: &config}
    container, err := client.CreateContainer(opts2)
    if err != nil {
        log.Error("Couldn't create container: ", err.Error())
        return
    } else {
        defer func() {
            err = client.RemoveContainer(docker.RemoveContainerOptions{
                ID: container.ID,
                Force: true,
            })
            if( err != nil) {
                log.Error("Coundn't remove container: ", container.ID, ": ", err.Error())
                return
            }
            log.Debug("Removed container with ID", container.ID)
        }()
    }
    log.Debug("Created container with ID", container.ID)

    //Try to start the container
    if err = dockerpty.Start(client, container, &docker.HostConfig{
        Binds: binds,
        PortBindings: portBindings,
    }); err != nil {
        log.Error(err.Error())
        return
    } else {
        // And once it is done with all the commands, remove the container.
        defer func () {
            err = client.StopContainer(container.ID, 0)
            if( err != nil) {
                log.Debug("Could not stop container ", container.ID, ": ", err.Error())
                return
            }
            log.Debug("Stopped container with ID", container.ID)
        }()
    }
    log.Debug("Started container with ID", container.ID)
}

func execMacro(macro Macro, project Project) {
    execInContainer(macro.getActions(), project)
}
