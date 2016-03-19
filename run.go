package main

import (
    "github.com/Sirupsen/logrus"
    "github.com/codegangsta/cli"
    "fmt"
    "github.com/fsouza/go-dockerclient"
    // import "strings" // words := strings.Fields(someString)
)

// Start the container, execute the list of commands, and then stops the container.
// This function should be used by run(), build(), etc.
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
    config := docker.Config{AttachStdout: true, AttachStdin: true, Image: imageName, Tty: true, OpenStdin: true}
    // opts2 := docker.CreateContainerOptions{Name: "nut_" + , Config: &config}
    opts2 := docker.CreateContainerOptions{Config: &config}
    container, err := client.CreateContainer(opts2)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    fmt.Println("Created container")

    //Try to start the container
    err = client.StartContainer(container.ID, &docker.HostConfig{})
    if( err != nil) {
        fmt.Println(err.Error())
        return
    }
    fmt.Println("Started container with ID", container.ID)

    // TODO : finish this function
    logrus.Println("RUN is not fully implemented yet...")

    // I think https://github.com/fsouza/go-dockerclient/issues/220#issuecomment-77777365
    // is a good starting point to mount volumes and bind sockets.

    // Then, instead of just starting the container, we should make it run the commands.
    // And once it is done with all the commands, delete the container.
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
