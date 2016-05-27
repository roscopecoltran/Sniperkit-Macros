package config

import(
    "path/filepath"
    Utils "github.com/matthieudelaro/nut/utils"
    containerFilepath "github.com/matthieudelaro/nut/container/filepath"
)

type VolumeV5 struct {
    VolumeBase `yaml:"inheritedValues,inline"`

    Host string `yaml:host_path`
    Container string `yaml:container_path`
}
        func (self *VolumeV5) getFullHostPath(context Utils.Context) (string, error) {
            clean := filepath.Clean(self.Host)
            if filepath.IsAbs(clean) {
                return clean, nil
            } else {
                return filepath.Join(context.GetRootDirectory(), clean), nil
            }
        }
        func (self *VolumeV5) getFullContainerPath(context Utils.Context) (string, error) {
            clean := containerFilepath.Clean(self.Container)
            if containerFilepath.IsAbs(clean) {
                return clean, nil
            } else {
                return containerFilepath.Join(context.GetRootDirectory(), clean), nil
            }
        }

type BaseEnvironmentV5 struct {
    BaseEnvironmentBase `yaml:"inheritedValues,inline"`

    FilePath string `yaml:"nut_file_path,omitempty"`
    GitHub string `yaml:"github,omitempty"`
    DockerImage string `yaml:"docker_image,omitempty"`
}
        func (self *BaseEnvironmentV5) getFilePath() string{
            return self.FilePath
        }
        func (self *BaseEnvironmentV5) getGitHub() string{
            return self.GitHub
        }

type ConfigV5 struct {
    ConfigBase `yaml:"inheritedValues,inline"`

    Mount map[string][]string `yaml:"mount,omitempty"`
    WorkingDir string `yaml:"container_working_directory,omitempty"`
    EnvironmentVariables map[string]string `yaml:"environment,omitempty"`
    Ports []string `yaml:"ports,omitempty"`
    EnableGUI string `yaml:"enable_gui,omitempty"`
    EnableNvidiaDevices string `yaml:"enable_nvidia_devices,omitempty"`
    Privileged string `yaml:"privileged,omitempty"`
    SecurityOpts []string `yaml:"security_opts,omitempty"`
    parent Config
}
        func (self *ConfigV5) getParent() Config {
            return self.parent
        }
        func (self *ConfigV5) getWorkingDir() string {
            return self.WorkingDir
        }
        func (self *ConfigV5) getVolumes() map[string]Volume {
            cacheVolumes := make(map[string]Volume)
            for name, data := range(self.Mount) {
                cacheVolumes[name] = &VolumeV5{
                    Host: data[0],
                    Container: data[1],
                }
            }
            return cacheVolumes
        }
        func (self *ConfigV5) getEnvironmentVariables() map[string]string {
            return self.EnvironmentVariables
        }
        func (self *ConfigV5) getPorts() []string {
            return self.Ports
        }
        func (self *ConfigV5) getEnableGui() (bool, bool) {
            return TruthyString(self.EnableGUI)
        }
        func (self *ConfigV5) getEnableNvidiaDevices() (bool, bool) {
            return TruthyString(self.EnableNvidiaDevices)
        }
        func (self *ConfigV5) getPrivileged() (bool, bool) {
            return TruthyString(self.Privileged)
        }
        func (self *ConfigV5) getSecurityOpts() []string {
            return self.SecurityOpts
        }

type ProjectV5 struct {
    SyntaxVersion string `yaml:"syntax_version"`
    ProjectName string `yaml:"project_name"`
    Base BaseEnvironmentV5 `yaml:"based_on"`
    Macros map[string]*MacroV5 `yaml:"macros,omitempty"`
    parent Project

    ProjectBase `yaml:"inheritedValues,inline"`
    ConfigV5 `yaml:"inheritedValues,inline"`
}
        func (self *ProjectV5) getDockerImage() string {
            return self.Base.DockerImage
        }
        func (self *ProjectV5) getSyntaxVersion() string {
            return self.SyntaxVersion
        }
        func (self *ProjectV5) getProjectName() string {
            return self.ProjectName
        }
        func (self *ProjectV5) getBaseEnv() BaseEnvironment {
            return &self.Base
        }
        func (self *ProjectV5) getMacros() map[string]Macro {
            // make the list of macros
            cacheMacros := make(map[string]Macro)
            for name, data := range self.Macros {
                data.parent = self
                cacheMacros[name] = data
            }
            return cacheMacros
        }
        // func (self *ProjectV5) createMacro(usage string, commands []string) Macro {
        //     return &MacroV5 {
        //         ConfigV5: *NewConfigV5(self,),
        //         Usage: usage,
        //         Actions: commands,
        //     }
        // }
        func (self *ProjectV5) getParent() Config {
            return self.parent
        }
        func (self *ProjectV5) getParentProject() Project {
            return self.parent
        }
        func (self *ProjectV5) setParentProject(project Project) {
            self.parent = project
        }

type MacroV5 struct {
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

    MacroBase `yaml:"inheritedValues,inline"`
    ConfigV5 `yaml:"inheritedValues,inline"`
}
        func (self *MacroV5) setParentProject(project Project) {
            self.ConfigV5.parent = project
        }
        func (self *MacroV5) getUsage() string {
            return self.Usage
        }
        func (self *MacroV5) getActions() []string {
            return self.Actions
        }
        func (self *MacroV5) getAliases() []string {
            return self.Aliases
        }
        func (self *MacroV5) getUsageText() string {
            return self.UsageText
        }
        func (self *MacroV5) getDescription() string {
            return self.Description
        }


func NewConfigV5(parent Config) *ConfigV5 {
    return &ConfigV5{
        Mount: make(map[string][]string),
        parent: parent,
    }
}

func NewProjectV5(parent Project) *ProjectV5 {
    project := &ProjectV5 {
        SyntaxVersion: "5",
        Macros: make(map[string]*MacroV5),
        ConfigV5: *NewConfigV5(nil),
        parent: parent,
    }
    return project
}
