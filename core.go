package main

import (
    // "github.com/Sirupsen/logrus"
    "fmt"
    "github.com/fsouza/go-dockerclient"
    // "bytes"
    "os"
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
    fmt.Println("CreateExec")
    if dExec, err = client.CreateExec(de); err != nil {
        fmt.Println("CreateExec Error: %s", err)
        return
    }
    fmt.Println("Created Exec")
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
    fmt.Println("StartExec")
    if err = client.StartExec(execId, opts); err != nil {
        fmt.Println("StartExec Error: %s", err)
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
func execInContainer(commands []string, proj *project) {
    fmt.Println(commands)
    for _, command := range commands {
        fmt.Printf(command)
        // for _, word := range command {
        //     logrus.Printf(word)
        //     logrus.Printf(", ")
        // }
        fmt.Printf("\n")
    }
    // import "strings" // words := strings.Fields(someString) // http://play.golang.org/

    endpoint := "unix:///var/run/docker.sock"
    client, err := docker.NewClient(endpoint)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    fmt.Println("Created client")

    //Pull image from Registry, if not present
    imageName := proj.Base.DockerImage
    _, err = client.InspectImage(imageName)
    if err != nil {
        fmt.Println(err.Error())
        fmt.Println("Pulling image...")

        opts := docker.PullImageOptions{Repository: imageName}
        err = client.PullImage(opts, docker.AuthConfiguration{})
        if err != nil {
            fmt.Println(err.Error())
            return
        }
        fmt.Println("Pulled image")
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
        WorkingDir:   proj.WorkingDir,
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
        fmt.Println(err.Error())
        return
    } else {
        defer func() {
            err = client.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID})
            if( err != nil) {
                fmt.Println(err.Error())
                return
            }
            fmt.Println("Removed container with ID", container.ID)
        }()
    }
    fmt.Println("Created container with ID", container.ID)


    //Try to start the container

    // prepare names of directories to mount
    mountingPoints := proj.getMountingPoints()
    binds := make([]string, 0, len(mountingPoints))
    for _, directory := range(mountingPoints) {
        hostPath, hostPathErr := directory.fullHostPath()
        containerPath, containerPathErr := directory.fullContainerPath()
        if hostPathErr != nil {
            fmt.Println(hostPathErr.Error())
            return
        }
        if containerPathErr != nil {
            fmt.Println(containerPathErr.Error())
            return
        }
        binds = append(binds, hostPath + ":" + containerPath)
    }
    fmt.Println("binds", binds)
    err = client.StartContainer(container.ID, &docker.HostConfig{
        Binds: binds,
    })
    if( err != nil) {
        fmt.Println(err.Error())
        return
    } else {
        // And once it is done with all the commands, remove the container.
        defer func () {
            err = client.StopContainer(container.ID, 0)
            if( err != nil) {
                fmt.Println(err.Error())
                return
            }
            fmt.Println("Stopped container with ID", container.ID)
        }()
    }
    fmt.Println("Started container with ID", container.ID)

    for _, command := range commands {
        execCommandInContainer(client, container, command)
    }


}

func execMacro(commands []string, proj *project) {
    execInContainer(commands, proj)
}
