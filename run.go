package main

import (
    "github.com/Sirupsen/logrus"
    "github.com/codegangsta/cli"
    "fmt"
    "github.com/fsouza/go-dockerclient"
    // "bytes"
    "os"
    "strings" // words := strings.Fields(someString)
)

func execCommandInContainer(client *docker.Client, container *docker.Container, command string) {
    var (
        dExec  *docker.Exec
        err    error
    )
    de := docker.CreateExecOptions{
        AttachStderr: true,
        AttachStdin:  true,
        AttachStdout: true,
        Tty:          true,
        // Cmd:          []string{"echo", "hello world2", "&&", "echo", "blabla"}, //, ";", "echo", "blabla2"},
        // Cmd:          []string{"/bin/sh"},
        // Cmd:          []string{"top"},
        Cmd:          strings.Fields(command),
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
    }
    fmt.Println("Started container with ID", container.ID)

    for _, command := range commands {
        execCommandInContainer(client, container, command)
    }

    // And once it is done with all the commands, remove the container.
    err = client.StopContainer(container.ID, 0)
    if( err != nil) {
        fmt.Println(err.Error())
        return
    }
    fmt.Println("Stopped container with ID", container.ID)

    err = client.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID})
    if( err != nil) {
        fmt.Println(err.Error())
        return
    }
    fmt.Println("Removed container with ID", container.ID)

}

func run(c *cli.Context) {
    p, err := loadProject()
    if err == nil {
        logrus.Println(p)
        commands := p.Macros["run"]
        execInContainer(commands, p)
        // execInContainer([][]string{[]string{"run"}}, p)
    } else {
        logrus.Println("Could not load nut file")
    }
}
