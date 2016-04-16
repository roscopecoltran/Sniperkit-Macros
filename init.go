package main

import (
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"github.com/codegangsta/cli"
	"fmt"
)

// create a nut.yml at the current path
func initSubcommand(c *cli.Context) {
	var project *ProjectV4 = NewProjectV4()

	project.ProjectName = "nut" // TODO: use the name of the current folder
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

	data := project.toYAML()
	fmt.Println("Project configuration:")
	fmt.Println("")
	fmt.Println(data)
	// check is nut.yml exists at the current path
	if nutFileExistsAtPath(".") {
		log.Error("Could not save new Nut project because a nut.yml file already exists")
		return
		// TODO: define and retrieve CLI parameters from the "c" argument
	} else {
		nutfileName := "./nut.yml"  // TODO: pick this name from a well documented and centralized list of legit nut file names
		err := ioutil.WriteFile(nutfileName, []byte(data), 0644)
		if err != nil {
			log.Error(err)
		} else {
			fmt.Printf("Project configuration saved in ", nutfileName)
		}
	}
}
