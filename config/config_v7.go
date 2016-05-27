package config

import(
    "path/filepath"
    Utils "github.com/matthieudelaro/nut/utils"
    containerFilepath "github.com/matthieudelaro/nut/container/filepath"
)

type VolumeV7 struct {
    VolumeBase `yaml:"inheritedValues,inline"`

    Host string `yaml:"host_path,omitempty"`
    Container string `yaml:"container_path,omitempty"`
    Options string `yaml:"options,omitempty"`
}
        func (self *VolumeV7) getHostPath() string {
            return self.Host
        }
        func (self *VolumeV7) getContainerPath() string {
            return self.Container
        }
        func (self *VolumeV7) getOptions() string {
            return self.Options
        }
        func (self *VolumeV7) fullHostPath(context Utils.Context) (string, error) {
            clean := filepath.Clean(self.Host)
            if filepath.IsAbs(clean) {
                return clean, nil
            } else {
                return filepath.Join(context.GetRootDirectory(), clean), nil
            }
        }
        func (self *VolumeV7) fullContainerPath(context Utils.Context) (string, error) {
            clean := containerFilepath.Clean(self.Container)
            if containerFilepath.IsAbs(clean) {
                return clean, nil
            } else {
                return containerFilepath.Join(context.GetRootDirectory(), clean), nil
            }
        }


type BaseEnvironmentV7 struct {
    BaseEnvironmentBase `yaml:"inheritedValues,inline"`

    FilePath string `yaml:"nut_file_path,omitempty"`
    GitHub string `yaml:"github,omitempty"`
}
        func (self *BaseEnvironmentV7) getFilePath() string{
            return self.FilePath
        }
        func (self *BaseEnvironmentV7) getGitHub() string{
            return self.GitHub
        }

type ConfigV7 struct {
    ConfigBase `yaml:"inheritedValues,inline"`

    DockerImage string `yaml:"docker_image,omitempty"`
    Mount map[string][]string `yaml:"mount,omitempty"`
    WorkingDir string `yaml:"container_working_directory,omitempty"`
    EnvironmentVariables map[string]string `yaml:"environment,omitempty"`
    Ports []string `yaml:"ports,omitempty"`
    EnableGUI string `yaml:"enable_gui,omitempty"`
    EnableNvidiaDevices string `yaml:"enable_nvidia_devices,omitempty"`
    Privileged string `yaml:"privileged,omitempty"`
    SecurityOpts []string `yaml:"security_opts,omitempty"`
    Detached string `yaml:"detached,omitempty"`
    UTSMode string `yaml:"uts,omitempty"`
    NetworkMode string `yaml:"net,omitempty"`
    Devices map[string]VolumeV7 `yaml:"devices,omitempty"`
    parent Config
}
        func (self *ConfigV7) getDockerImage() string {
            return self.DockerImage
        }
        func (self *ConfigV7) getNetworkMode() string {
            return self.NetworkMode
        }
        func (self *ConfigV7) getUTSMode() string {
            return self.UTSMode
        }
        func (self *ConfigV7) getParent() Config {
            return self.parent
        }
        func (self *ConfigV7) getWorkingDir() string {
            return self.WorkingDir
        }
        func (self *ConfigV7) getVolumes() map[string]Volume {
            cacheVolumes := make(map[string]Volume)
            for name, data := range(self.Mount) {
                cacheVolumes[name] = &VolumeV7{
                    Host: data[0],
                    Container: data[1],
                }
            }
            return cacheVolumes
        }
        func (self *ConfigV7) getEnvironmentVariables() map[string]string {
            return self.EnvironmentVariables
        }
        func (self *ConfigV7) getDevices() map[string]Volume {
            cacheVolumes := make(map[string]Volume)
            for name, data := range(self.Devices) {
                cacheVolumes[name] = &data
            }
            return cacheVolumes
        }
        func (self *ConfigV7) getPorts() []string {
            return self.Ports
        }
        func (self *ConfigV7) getEnableGui() (bool, bool) {
            return TruthyString(self.EnableGUI)
        }
        func (self *ConfigV7) getEnableNvidiaDevices() (bool, bool) {
            return TruthyString(self.EnableNvidiaDevices)
        }
        func (self *ConfigV7) getPrivileged() (bool, bool) {
            return TruthyString(self.Privileged)
        }
        func (self *ConfigV7) getDetached() (bool, bool) {
            return TruthyString(self.Detached)
        }
        func (self *ConfigV7) getSecurityOpts() []string {
            return self.SecurityOpts
        }

type ProjectV7 struct {
    SyntaxVersion string `yaml:"syntax_version"`
    ProjectName string `yaml:"project_name"`
    Base BaseEnvironmentV7 `yaml:"based_on,omitempty"`
    Macros map[string]*MacroV7 `yaml:"macros,omitempty"`
    parent Project

    ProjectBase `yaml:"inheritedValues,inline"`
    ConfigV7 `yaml:"inheritedValues,inline"`
}
        func (self *ProjectV7) getSyntaxVersion() string {
            return self.SyntaxVersion
        }
        func (self *ProjectV7) getProjectName() string {
            return self.ProjectName
        }
        func (self *ProjectV7) getBaseEnv() BaseEnvironment {
            return &self.Base
        }
        func (self *ProjectV7) getMacros() map[string]Macro {
            // make the list of macros
            cacheMacros := make(map[string]Macro)
            for name, data := range self.Macros {
                data.parent = self
                cacheMacros[name] = data
            }
            return cacheMacros
        }
        // func (self *ProjectV7) createMacro(usage string, commands []string) Macro {
        //     return &MacroV7 {
        //         ConfigV7: *NewConfigV7(self,),
        //         Usage: usage,
        //         Actions: commands,
        //     }
        // }
        func (self *ProjectV7) getParent() Config {
            return self.parent
        }
        func (self *ProjectV7) getParentProject() Project {
            return self.parent
        }
        func (self *ProjectV7) setParentProject(project Project) {
            self.parent = project
        }

type MacroV7 struct {
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
    ConfigV7 `yaml:"inheritedValues,inline"`
}
        func (self *MacroV7) setParentProject(project Project) {
            self.ConfigV7.parent = project
        }
        func (self *MacroV7) getUsage() string {
            return self.Usage
        }
        func (self *MacroV7) getActions() []string {
            return self.Actions
        }
        func (self *MacroV7) getAliases() []string {
            return self.Aliases
        }
        func (self *MacroV7) getUsageText() string {
            return self.UsageText
        }
        func (self *MacroV7) getDescription() string {
            return self.Description
        }


func NewConfigV7(parent Config) *ConfigV7 {
    return &ConfigV7{
        Mount: make(map[string][]string),
        parent: parent,
    }
}

func NewProjectV7(parent Project) *ProjectV7 {
    project := &ProjectV7 {
        SyntaxVersion: "7",
        Macros: make(map[string]*MacroV7),
        ConfigV7: *NewConfigV7(nil),
        parent: parent,
    }
    return project
}
