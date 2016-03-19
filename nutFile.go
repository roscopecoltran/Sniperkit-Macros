package main

import (
	"fmt"
	// "os"
	"github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type MountArgument struct {
	host string `yaml:host_path`
	container string `yaml:container_path`
}

type BaseEnvironment struct {
	DockerImage string `yaml:"docker_image,omitempty"`
	FileURL string `yaml:"nut_file_url,omitempty"`
	FilePath string `yaml:"nut_file_path,omitempty"`
}

// project represents a nut project
type project struct {
	SyntaxVersion string `yaml:"syntax_version"`
	ProjectName string `yaml:"project_name"`
	Base BaseEnvironment `yaml:"based_on"`
	WorkingDir string `yaml:"container_working_directory,omitempty"`
	Mount map[string][]string `yaml:"mount,omitempty"`
	Macros map[string][]string `yaml:"macros,omitempty"`
}

func NewProject() *project {
	var p *project = &project{
		ProjectName: "",
		Base: BaseEnvironment{},
		Macros: make(map[string][]string),
		Mount: make(map[string][]string),
	}
	for i, v := range p.Mount {
		fmt.Println(i)
		fmt.Println(v)
	}
	return p
}

func loadProject() (*project, error) {
	return parseNutFileAtPath("nut.yml")
}

// object to YAML
func (p *project) Marshal() string {
	d, err := yaml.Marshal(&p)
	if err != nil {
		logrus.Fatalf("error: %v", err)
	}
	// logrus.Printf("--- t dump:\n%s\n\n", string(d))
	return string(d)
}

// func (p *project) InitFromYaml(yaml string) {
// 	logrus.Fatal("NOT IMPLEMENTED YET")
// }
