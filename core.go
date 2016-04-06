package main

import (
    // "github.com/Sirupsen/logrus"
    "github.com/fsouza/go-dockerclient"
    // "bytes"
    "os"
    // "net"
    // "io"
    log "github.com/Sirupsen/logrus"
    // "time"
    // "strings" // words := strings.Fields(someString)
)

func execCommandInContainer(client *docker.Client, container *docker.Container, command string) {
    var (
        dExec  *docker.Exec
        err    error
    )

    // Issues with command = "ls -la | grep '323'" (the pipe is an issue in most cases)
    // cmd := []string{command}
        // "ls -la | grep '323'": executable file not found in $PATH
    // cmd := strings.Fields(command)
        // ls: cannot access |: No such file or directory
        // ls: cannot access grep: No such file or directory
        // ls: cannot access '323': No such file or directory
    // cmd := strings.Fields("bash -c \"" + command + "\"")
        // -la: -c: line 0: unexpected EOF while looking for matching `"'
        // -la: -c: line 1: syntax error: unexpected end of file
    cmd := []string{"bash", "-c", command}
        // runs properly. Using bash does not seem like an elegant solution,
        // but this is the best so far.
    // for i, v := range cmd {
    //     logrus.Println("command field ", i, ": ", v)
    // }
    de := docker.CreateExecOptions{
        AttachStderr: true,
        AttachStdin:  true,
        AttachStdout: true,
        Tty:          true,
        // Cmd:          []string{"echo", "hello world2", "&&", "echo", "blabla"}, //, ";", "echo", "blabla2"},
        // Cmd:          []string{"/bin/sh"},
        // Cmd:          []string{"top"},
        Cmd:          cmd,
        Container:    container.ID,
    }
    log.Debug("CreateExec")
    if dExec, err = client.CreateExec(de); err != nil {
        log.Debug("CreateExec Error: %s", err)
        return
    }
    log.Debug("Created Exec")
    // var stdout, stderr bytes.Buffer
    // var reader = strings.NewReader("echo hello world")
    execId := dExec.ID

    opts := docker.StartExecOptions{
        OutputStream: os.Stdout,
        ErrorStream:  os.Stderr,
        // InputStream:  reader,
        InputStream:  os.Stdin,
        RawTerminal:  true,
    }
    log.Debug("StartExec")
    if err = client.StartExec(execId, opts); err != nil {
        log.Debug("StartExec Error: %s", err)
        return
    }
}

// Start the container, execute the list of commands, and then stops the container.
// This function should be used by run(), build(), etc.
// Inspired from https://github.com/fsouza/go-dockerclient/issues/287
// and from https://github.com/fsouza/go-dockerclient/issues/220#issuecomment-77777365
// TODO: solve key issue (Ctrl+C, q, ...). See if that can help: https://github.com/fgrehm/go-dockerpty
// TODO: to solve tty size, use func (c *Client) ResizeContainerTTY(id string, height, width int) error
// and/or func (c *Client) ResizeExecTTY(id string, height, width int) error
// TODO: fix bug "StartExec Error: %s read /dev/stdin: bad file descriptor" when executing several commands
func execInContainer(commands []string, project Project) {
    imageName, err := project.getBaseEnv().getDockerImageName()
    if err != nil {
        log.Error(err.Error())
        return
    }

    endpoint := getDockerEndpoint()
    client, err := docker.NewClient(endpoint)
    if err != nil {
        log.Debug(err.Error())
        return
    }
    log.Debug("Created client")

    //Pull image from Registry, if not present
    _, err = client.InspectImage(imageName)
    if err != nil {
        log.Error("Could not pull image:", err.Error())
        log.Debug("Pulling image...")

        opts := docker.PullImageOptions{Repository: imageName}
        err = client.PullImage(opts, docker.AuthConfiguration{})
        if err != nil {
            log.Debug(err.Error())
            return
        }
        log.Debug("Pulled image")
    }


    log.SetLevel(log.DebugLevel)
    portBindings := map[docker.Port][]docker.PortBinding{}
    envVariables := []string{}
    if project.getEnableGui() {
        portBindings, envVariables, err = enableGui(project)
        if err != nil {
            log.Error("Could not enable GUI: ", err)
        }
    }

    //Try to create a container from the imageID
    // config := docker.Config{AttachStdout: true, AttachStdin: true, Image: imageName, Tty: true, OpenStdin: true}
    config := docker.Config{
        Image: imageName,
        // Cmd:          []string{"/bin/sh"},
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
    // User: set it to the user of the host, instead of root

    // I think https://github.com/fsouza/go-dockerclient/issues/220#issuecomment-77777365
    // is a good starting point to mount volumes and bind sockets.

    //
    // opts2 := docker.CreateContainerOptions{Name: "nut_" + , Config: &config}
    opts2 := docker.CreateContainerOptions{Config: &config}
    container, err := client.CreateContainer(opts2)
    if err != nil {
        log.Debug(err.Error())
        return
    } else {
        defer func() {
            err = client.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID})
            if( err != nil) {
                log.Debug(err.Error())
                return
            }
            log.Debug("Removed container with ID", container.ID)
        }()
    }
    log.Debug("Created container with ID", container.ID)


    //Try to start the container

    // prepare names of directories to mount
    mountingPoints := project.getMountingPoints()
    binds := make([]string, 0, len(mountingPoints))
    for _, directory := range(mountingPoints) {
        hostPath, hostPathErr := directory.fullHostPath()
        containerPath, containerPathErr := directory.fullContainerPath()
        if hostPathErr != nil {
            log.Debug(hostPathErr.Error())
            return
        }
        if containerPathErr != nil {
            log.Debug(containerPathErr.Error())
            return
        }
        binds = append(binds, hostPath + ":" + containerPath)
    }
    log.Debug("binds", binds)

    err = client.StartContainer(container.ID, &docker.HostConfig{
        Binds: binds,
        PortBindings: portBindings,
    })
    if( err != nil) {
        log.Debug(err.Error())
        return
    } else {
        // And once it is done with all the commands, remove the container.
        defer func () {
            err = client.StopContainer(container.ID, 0)
            if( err != nil) {
                log.Debug(err.Error())
                return
            }
            log.Debug("Stopped container with ID", container.ID)
        }()
    }
    log.Debug("Started container with ID", container.ID)

    for _, command := range commands {
        execCommandInContainer(client, container, command)
    }
}

func execMacro(macro Macro, project Project) {
    execInContainer(macro.getActions(), project)
}
