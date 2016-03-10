package main

import (
	"github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// project represents a nut project
type project struct {
	ProjectName string `yaml:"project_name"`
	DockerImage string `yaml:"docker_image"`
}

func NewProject() *project {
	var p *project = &project{
		ProjectName: "",
		DockerImage: "",
	}
	return p
}

// object to YAML
func (p *project) Marshal() string {
	d, err := yaml.Marshal(&p)
	if err != nil {
		logrus.Fatalf("error: %v", err)
	}
	logrus.Printf("--- t dump:\n%s\n\n", string(d))
	return string(d)
}

// func (p *project) InitFromYaml(yaml string) {
// 	logrus.Fatal("NOT IMPLEMENTED YET")
// }
