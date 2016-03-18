package main

import (
	"fmt"
	"io/ioutil"

	"github.com/codegangsta/cli"
)

func initSubcommand(c *cli.Context) {
	// check is nut.yml exists at the current path
	if nutFileExistsAtPath(".") {
		fmt.Println("Cannot init new Nut project because a nut.yml file already exists")
		return
	}
	// TODO: define and retrieve CLI parameters from the "c" argument

	// create a nut.yml at the current path
	var p *project = NewProject()
	p.SyntaxVersion = "1"
	p.ProjectName = "nut" // TODO: use the name of the current folder
	p.Base.DockerImage = "golang:1.6"
	// p.Mount.append(MountArgument{
	// 	host: ".",
	// 	container: "/go/src/project"
	// 	})
	// append(p.Mount, MountArgument{".", "/go/src/project"})
	// p.Mount = []MountArgument{MountArgument{".", "/go/src/project"}}
	// p.Mount["main"] = MountArgument{".", "/go/src/project"}
	// p.Mount["main"] = MountArgument{
		// host:".",
		// container:"/go/src/project"}
	p.Mount["main"] = []string{".", "/go/src/project"}
	p.Macros["build"] = []string{"go build -o output"}
	p.Macros["run"] = []string{"./output"}
	p.WorkingDir = "/go/src/project"

	data := p.Marshal()
	nutfileName := "./nut.yml"  // TODO: pick this name from a well documented and centralized list of legit nut file names
	err := ioutil.WriteFile(nutfileName, []byte(data), 0644)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Project created: configuration saved in", nutfileName)
		fmt.Println("")
		fmt.Println(data)
	}

}
