package main

import (
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"github.com/codegangsta/cli"
	"fmt"
	"path/filepath"
)

// create a nut.yml at the current path
func initSubcommand(c *cli.Context, gitHubFlag string) {
	name, err := filepath.Abs(".") // TODO: handle path properly, instead of just "."
	if err != nil {
		log.Debug("Could not get name of the project: use default.")
		name = "nut"
	} else {
		name = filepath.Base(name)
	}

	var project *ProjectV4 = NewProjectV4()
	project.ProjectName = name
	if gitHubFlag == "" {// create an example of configuration
		// project.ProjectName = "nut" // TODO: use the name of the current folder
		project.Base.DockerImage = "golang:1.6"
		project.Mount["main"] = []string{".", "/go/src/project"}
		project.Macros["run"] = &MacroV4{
			Usage: "run the project in the container",
			Actions: []string{"./nut"},
		}
		project.Macros["build"] = &MacroV4{
			Usage: "build the project",
			Actions: []string{"go build -o nut"},
		}
		project.Macros["build-osx"] = &MacroV4{
			Usage: "build the project for OSX",
			Actions: []string{"env GOOS=darwin GOARCH=amd64 go build -o nutOSX"},
			Aliases: []string{"bo"},
			Description: "cross-compile the project to run on OSX, with architecture amd64.",
		}

		project.Macros["build-linux"] = &MacroV4{
			Usage: "build the project for Linux",
			Actions: []string{"env GOOS=linux GOARCH=amd64 go build -o nutLinux"},
			Aliases: []string{"bl"},
		}
		project.Macros["build-windows"] = &MacroV4{
			Usage: "build the project for Windows",
			Actions: []string{"env GOOS=windows GOARCH=amd64 go build -o nutWindows"},
			Aliases: []string{"bw"},
		}
		project.Macros["build-all"] = &MacroV4{
			Usage: "build the project for linux, OSX, and for Windows",
			Actions: []string{
				"echo Building for linux...",
				"env GOOS=linux GOARCH=amd64 go build -o nutLinux",
				"echo Building for OSX...",
				"env GOOS=darwin GOARCH=amd64 go build -o nutOSX",
				"echo Building for Windows...",
				"env GOOS=windows GOARCH=amd64 go build -o nutWindows",
			},
			Aliases: []string{"ba"},
		}
		project.WorkingDir = "/go/src/project"
	} else {
		log.Debug("init from github: ", gitHubFlag)
		// Download dependencies, and analyse base configuration
		baseConfig := NewProjectV4()
		baseConfig.getBaseEnv().setGitHub(gitHubFlag)
		err := LoadProjectHierarchy(baseConfig)
		if err != nil {
			log.Error("Could not init from GitHub repository ", gitHubFlag,
				": ", err)
			return
		}

		// create project depending on parent configuration
			// 1 - mount folder "." if not already mounted by parent configuration
			mountingPointName := "name"
			hostDir := "."
			containerDir := "/nut/" + name
			mountingPoint := &MountingPointV4{
				host: hostDir,
				container: containerDir,
			}
			if CheckConflict(mountingPointName, mountingPoint, baseConfig.getMountingPoints()) == nil {
				project.Mount[mountingPointName] = []string{hostDir, containerDir}
			}
			// 2 - set working directory to "." if not specified otherwise
			if baseConfig.getWorkingDir() == "" {
				project.WorkingDir = containerDir
			}
		project.getBaseEnv().setGitHub(gitHubFlag)
		project.setParentProject(baseConfig.getParentProject())
	}

	data := project.toYAML()
	fmt.Println("Project configuration:")
	fmt.Println("")
	fmt.Println(data)
	// check is nut.yml exists at the current path
	if nutFileExistsAtPath(".") {
		log.Error("Could not save new Nut project because a nut.yml file already exists")
		return
	} else {
		nutfileName := "./nut.yml"  // TODO: pick this name from a well documented and centralized list of legit nut file names
		err := ioutil.WriteFile(nutfileName, []byte(data), 0644) // TODO: discuss this permission level
		if err != nil {
			log.Error(err)
		} else {
			fmt.Printf("Project configuration saved in %s.", nutfileName) // TODO: make sure that bug Project configuration saved in %!(EXTRA string=./nut.yml) is fixed.
		}
	}
}
