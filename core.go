package main

import (
    "github.com/fsouza/go-dockerclient"
    // "github.com/fgrehm/go-dockerpty"
    "github.com/matthieudelaro/go-dockerpty"
    log "github.com/Sirupsen/logrus"
    "strings"
    "fmt"
    "github.com/matthieudelaro/nut/nvidia"
    Config "github.com/matthieudelaro/nut/config"
    Utils "github.com/matthieudelaro/nut/utils"
)

// Start the container, execute the list of commands, and then stops the container.
// Inspired from https://github.com/fsouza/go-dockerclient/issues/287
// and from https://github.com/fsouza/go-dockerclient/issues/220#issuecomment-77777365
// DONE: solve key issue (Ctrl+C, q, ...). See if that can help: https://github.com/fgrehm/go-dockerpty : it does :)
// DONE: to solve tty size, use func (c *Client) ResizeContainerTTY(id string, height, width int) error : solved with dockerpty
// and/or func (c *Client) ResizeExecTTY(id string, height, width int) error
// TODO: report bug "StartExec Error: %s read /dev/stdin: bad file descriptor" when executing several commands : post issue on dockerpty
// func execMacro(macro Config.Macro) {
func execInContainer(commands []string, config Config.Config, context Utils.Context) {
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

    imageName := Config.GetDockerImage(config)
    if imageName == "" {
        log.Error("Docker image has not been defined.")
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
        fmt.Println("Could not inspect image", imageName, ":", err.Error())

        fmt.Println("Pulling image...")
        opts := docker.PullImageOptions{Repository: imageName}
        err = client.PullImage(opts, docker.AuthConfiguration{})
        if err != nil {
            log.Error("Could not pull image ", imageName, ": ", err.Error())
            return
        }
        fmt.Println("Pulled image")
    }

    // prepare names of directories to mount
    // inspired from https://github.com/fsouza/go-dockerclient/issues/220#issuecomment-77777365
    mountingPoints := Config.GetVolumes(config, context)
    binds := make([]string, 0, len(mountingPoints))
    portBindings := map[docker.Port][]docker.PortBinding{}
    exposedPorts := map[docker.Port]struct{}{}
    envVariables := []string{}
    volumeDriver := ""
    devices := []docker.Device{}

    for _, directory := range(mountingPoints) {
        hostPath, hostPathErr := Config.GetFullHostPath(directory, context)
        containerPath, containerPathErr := Config.GetFullContainerPath(directory, context)
        if hostPathErr != nil {
            log.Error("Couldn't mount host directory: ", hostPathErr.Error())
            return
        }
        if containerPathErr != nil {
            log.Error("Couldn't mount container directory: ", containerPathErr.Error())
            return
        }
        binds = append(binds, hostPath + ":" + containerPath)
    }
    log.Debug("binds", binds)

    for _, value := range Config.GetPorts(config) {
        parts := strings.Split(value, ":") // TODO: support ranges of ports
        hostPort := ""
        containerPort := ""
        if len(parts) == 2 {
            hostPort = parts[0]
            containerPort = parts[1]
        } else if len(parts) == 1 {
            hostPort = parts[0]
            containerPort = parts[0]
        } else {
            log.Error("Could not parse port: " + value)
            return
        }
        // name := containerPort + "/tcp" // TODO: support UDP
        // dockerPort := docker.Port{containerPort + "/tcp"}
        var dockerPort docker.Port = docker.Port(containerPort + "/tcp")
        exposedPorts[dockerPort] = struct{}{}
        portBindings[dockerPort] = []docker.PortBinding{
            // {HostIP: "0.0.0.0", HostPort: "8080",}}
            {HostPort: hostPort,}} // TODO: support HostIP
    }
    for name, value := range Config.GetEnvironmentVariables(config) {
        envVariables = append(envVariables, name + "=" + value)
    }
    if Config.IsGUIEnabled(config) {
        portBindingsGUI, envVariablesGUI, bindsGUI, err := enableGui()
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
    if Config.IsNvidiaDevicesEnabled(config) {
        nvidiaDevices, driverName, driverVolume, err := nvidia.GetConfiguration()
        if err != nil {
            log.Error("Could not enable Nvidia devices: ", err,
                "\nYou have to be on Linux for this to work. Also, make sure " +
                "that nvidia-docker-plugin is running.\n")
        } else {
            binds = append(binds, driverVolume)
            volumeDriver = driverName
            for _, devicePath := range nvidiaDevices {
                devices = append(devices, docker.Device{
                    PathOnHost: devicePath,
                    PathInContainer: devicePath,
                    CgroupPermissions: "mrw", // TODO: discuss proper CgroupPermissions
                })
            }
        }
    }

    //Try to create a container from the imageID
    dockerConfig := docker.Config{
        Image: imageName,
        Cmd:          cmdConfig,
        OpenStdin:    true,
        StdinOnce:    true,
        AttachStdin:  true,
        AttachStdout: true,
        AttachStderr: true,
        Tty:          true,
        WorkingDir:   Config.GetWorkingDir(config),
        Env:          envVariables,
        ExposedPorts: exposedPorts,
        VolumeDriver: volumeDriver,
    }
    // TODO : set following config options: https://godoc.org/github.com/fsouza/go-dockerclient#Config
    // User: set it to the user of the host, instead of root, to manage file permissions properly

    // TODO: ? Give the container a name? Can be done with docker.CreateContainerOptions{Name: "nut_myproject"}
    // opts2 := docker.CreateContainerOptions{Name: "nut_" + , Config: &dockerConfig}
    opts2 := docker.CreateContainerOptions{Config: &dockerConfig}
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
            log.Debug("Removed container with ID ", container.ID)
        }()
    }
    log.Debug("Created container with ID ", container.ID)

    dockerHostConfig := docker.HostConfig{
        Binds: binds,
        PortBindings: portBindings,
        Privileged: Config.IsPrivileged(config),
        SecurityOpt: Config.GetSecurityOpts(config),
        Devices: devices,
        UTSMode: Config.GetUTSMode(config),
    }

    //Try to start the container
    if err = dockerpty.Start(client, container, &dockerHostConfig); err != nil {
        log.Error("Error while starting container, and attaching to it: ", err.Error(),
            "\nBinds: ", dockerHostConfig.Binds,
            "\nPortBindings: ", dockerHostConfig.PortBindings,
            "\nPrivileged: ", dockerHostConfig.Privileged,
            "\nSecurityOpt: ", dockerHostConfig.SecurityOpt,
            "\nDevices: ", dockerHostConfig.Devices,

            "\nImage: ", dockerConfig.Image,
            "\nCmd: ", dockerConfig.Cmd,
            "\nWorkingDir: ", dockerConfig.WorkingDir,
            "\nEnv: ", dockerConfig.Env,
            "\nExposedPorts: ", dockerConfig.ExposedPorts,
            "\nVolumeDriver: ", dockerConfig.VolumeDriver,
            )
        return
    } else {
        // And once it is done with all the commands, remove the container.
        defer func () {
            err = client.StopContainer(container.ID, 0)
            if( err != nil) {
                log.Debug("Could not stop container ", container.ID, ": ", err.Error())
                return
            }
            log.Debug("Stopped container with ID ", container.ID)
        }()
    }
    log.Debug("Started container with ID ", container.ID)
}

func execMacro(macro Config.Macro, context Utils.Context) {
    execInContainer(Config.GetActions(macro), macro, context)
}
