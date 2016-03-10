package main

import (
	"fmt"

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
	p.ProjectName = "nut"
	p.DockerImage = "golang:1.6"

	p.Marshal()
	// err := ioutil.WriteFile("./nut.yml", data, perm)
	fmt.Println("project created (NOT IMPLEMENTED YET)")

}
