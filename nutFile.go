package main

import (
	"fmt"
	// "os"
	"github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
    "path/filepath"

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

type Macro struct {
	// A short description of the usage of this macro
	Usage string `yaml:"usage,omitempty"`
	// The commands to execute when this macro is invoked
	Actions []string `yaml:"actions,omitempty"`
	// A list of aliases for the macro
	Aliases []string `yaml:"aliases,omitempty"`
    // Custom text to show on USAGE section of help
    UsageText string `yaml:"usage_for_help_section,omitempty"`
    // A longer explanation of how the macro works
    Description string `yaml:"description,omitempty"`
}

// project represents a nut project
type Project struct {
	SyntaxVersion string `yaml:"syntax_version"`
	ProjectName string `yaml:"project_name"`
	Base BaseEnvironment `yaml:"based_on"`
	WorkingDir string `yaml:"container_working_directory,omitempty"`
	Mount map[string][]string `yaml:"mount,omitempty"`
	// Macros map[string][]string `yaml:"macros,omitempty"`
	Macros map[string]Macro `yaml:"macros,omitempty"`
}

func NewProject() *Project {
	var project *Project = &Project{
		ProjectName: "",
		Base: BaseEnvironment{},
		Mount: make(map[string][]string),
		// Macros: make(map[string][]string),
		Macros: make(map[string]Macro),
	}
	for i, v := range project.Mount {
		fmt.Println(i)
		fmt.Println(v)
	}
	return project
}

func loadProject() (*Project, error) {
	return parseNutFileAtPath("nut.yml")
}

// object to YAML
func (project *Project) Marshal() string {
	d, err := yaml.Marshal(&project)
	if err != nil {
		logrus.Fatalf("error: %v", err)
	}
	// logrus.Printf("--- t dump:\n%s\n\n", string(d))
	return string(d)
}

func (project *Project) getMountingPoints() map[string]MountArgument {
	mountingPoints := make(map[string]MountArgument)
	for name, data := range(project.Mount) {
		mountingPoints[name] = MountArgument{
			host: data[0],
			container: data[1],
		}
	}
	return mountingPoints
}

func (self *MountArgument) fullHostPath() (string, error) {
	absolutePath, err := filepath.Abs(self.host)
	return absolutePath, err
}

func (self *MountArgument) fullContainerPath() (string, error) {
	absolutePath, err := filepath.Abs(self.container)
	return absolutePath, err
}

// func (project *project) InitFromYaml(yaml string) {
// 	logrus.Fatal("NOT IMPLEMENTED YET")
// }
