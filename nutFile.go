package main

import (
	"github.com/Sirupsen/logrus"
)

// project represents a nut project
type project struct {
	ProjectName string
	DockerImage string
}

func NewProject() *project {
	var p *project = &project{
		ProjectName: "",
		DockerImage: "",
	}
	return p
}

func (p *project) InitFromYaml(yaml string) {
	logrus.Fatal("NOT IMPLEMENTED YET")
}
