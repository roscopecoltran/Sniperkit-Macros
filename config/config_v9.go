package config

import(
    "path/filepath"
    Utils "github.com/roscopecoltran/sniperkit-macros/utils"
    containerFilepath "github.com/roscopecoltran/sniperkit-macros/container/filepath"

    log "github.com/sirupsen/logrus"
    "errors"
    "strings"

    "os"
)

type VolumeV9 struct {
    VolumeBase `yaml:"inheritedValues,inline"`

    VolumeName string `yaml:"volume_name,omitempty"`
    Host string `yaml:"host_path,omitempty"`
    Container string `yaml:"container_path"`
    Options string `yaml:"options,omitempty"`
}
        func (self *VolumeV9) getVolumeName() string {
            return self.VolumeName
        }
        func (self *VolumeV9) getOptions() string {
            return self.Options
        }
        // TODO: move the content of the two next function to an appropriate place
        // Same for V6 and V5
        func (self *VolumeV9) getFullHostPath(context Utils.Context) (string, error) {
            if self.Host == "" {
                return "", errors.New("Undefined host path")
            } else {
                // log.Debug("getFullHostPath value ", self.Host)
                res := filepath.Clean(self.Host)
                // log.Debug("getFullHostPath clean ", res)
                if !filepath.IsAbs(res) {
                    // log.Debug("getFullHostPath is not Abs")
                    res = filepath.Join(context.GetRootDirectory(), res)
                }
                // log.Debug("getFullHostPath value ", res)
                if strings.Contains(res, `:\`) { // path on windows. Eg: C:\\Users\
                    // log.Debug("getFullHostPath windows ", `:\`)
                    parts := strings.Split(res, `:\`)
                    parts[0] = strings.ToLower(parts[0]) // drive letter should be lower case
                    res = "//" + parts[0] + "/" + filepath.ToSlash(parts[1])
                }

                log.Debug("getFullHostPath res ", res)
                return res, nil
            }
        }
        func (self *VolumeV9) getFullContainerPath(context Utils.Context) (string, error) {
            if self.Container == "" {
                return "", errors.New("Undefined container path")
            } else {
                // log.Debug("getFullContainerPath value ", self.Container)
                clean := containerFilepath.ToSlash(containerFilepath.Clean(self.Container))
                // log.Debug("getFullContainerPath clean ", clean)
                if containerFilepath.IsAbs(clean) {
                    // log.Debug("getFullContainerPath isAbs")
                    return clean, nil
                } else {
                    // log.Debug("getFullContainerPath is not Abs")
                    log.Debug("getFullContainerPath return ", containerFilepath.Join(context.GetRootDirectory(), clean))
                    return containerFilepath.Join(context.GetRootDirectory(), clean), nil
                }
            }
        }


type DeviceV9 struct {
    DeviceBase `yaml:"inheritedValues,inline"`

    Host string `yaml:"host_path"`
    Container string `yaml:"container_path"`
    Options string `yaml:"options,omitempty"`
}
        func (self *DeviceV9) getHostPath() string {
            return self.Host
        }
        func (self *DeviceV9) getContainerPath() string {
            return self.Container
        }
        func (self *DeviceV9) getOptions() string {
            return self.Options
        }


type BaseEnvironmentV9 struct {
    BaseEnvironmentBase `yaml:"inheritedValues,inline"`

    FilePath string `yaml:"nut_file_path,omitempty"`
    GitHub string `yaml:"github,omitempty"`
}
        func (self *BaseEnvironmentV9) getFilePath() string{
            return self.FilePath
        }
        func (self *BaseEnvironmentV9) getGitHub() string{
            return self.GitHub
        }

type ConfigV9 struct {
    ConfigBase `yaml:"inheritedValues,inline"`

    DockerImage string `yaml:"docker_image,omitempty"`
    Volumes map[string]*VolumeV9 `yaml:"volumes,omitempty"`
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
    Devices map[string]*DeviceV9 `yaml:"devices,omitempty"`
    EnableCurrentUser string `yaml:"enable_current_user,omitempty"`
    WorkInProjectFolderAs string `yaml:"work_in_project_folder_as,omitempty"`
    parent Config
}
        func (self *ConfigV9) getDockerImage() string {
            return self.DockerImage
        }
        func (self *ConfigV9) getNetworkMode() string {
            return self.NetworkMode
        }
        func (self *ConfigV9) getUTSMode() string {
            return self.UTSMode
        }
        func (self *ConfigV9) getParent() Config {
            return self.parent
        }
        func (self *ConfigV9) getWorkingDir() string {
            return self.WorkingDir
        }
        func (self *ConfigV9) getVolumes() map[string]Volume {
            cacheVolumes := make(map[string]Volume)
            for name, data := range(self.Volumes) {
                cacheVolumes[name] = data
            }
            return cacheVolumes
        }
        func (self *ConfigV9) getEnvironmentVariables() map[string]string {
            return self.EnvironmentVariables
        }
        func (self *ConfigV9) getDevices() map[string]Device {
            cacheVolumes := make(map[string]Device)
            for name, data := range(self.Devices) {
                cacheVolumes[name] = data
            }
            return cacheVolumes
        }
        func (self *ConfigV9) getPorts() []string {
            return self.Ports
        }
        func (self *ConfigV9) getEnableGui() (bool, bool) {
            return TruthyString(self.EnableGUI)
        }
        func (self *ConfigV9) getEnableNvidiaDevices() (bool, bool) {
            return TruthyString(self.EnableNvidiaDevices)
        }
        func (self *ConfigV9) getPrivileged() (bool, bool) {
            return TruthyString(self.Privileged)
        }
        func (self *ConfigV9) getDetached() (bool, bool) {
            return TruthyString(self.Detached)
        }
        func (self *ConfigV9) getEnableCurrentUser() (bool, bool) {
            return TruthyString(self.EnableCurrentUser)
        }
        func (self *ConfigV9) getSecurityOpts() []string {
            return self.SecurityOpts
        }
        func (self *ConfigV9) getWorkInProjectFolderAs() string {
            return self.WorkInProjectFolderAs
        }

type ProjectV9 struct {
    SyntaxVersion string `yaml:"syntax_version"`
    ProjectName string `yaml:"project_name"`
    Base BaseEnvironmentV9 `yaml:"based_on,omitempty"`
    Macros map[string]*MacroV9 `yaml:"macros,omitempty"`
    parent Project

    ProjectBase `yaml:"inheritedValues,inline"`
    ConfigV9 `yaml:"inheritedValues,inline"`
}
        func (self *ProjectV9) getSyntaxVersion() string {
            return self.SyntaxVersion
        }
        func (self *ProjectV9) getProjectName() string {
            return self.ProjectName
        }
        func (self *ProjectV9) getBaseEnv() BaseEnvironment {
            return &self.Base
        }
        func (self *ProjectV9) getMacros() map[string]Macro {
            // make the list of macros
            cacheMacros := make(map[string]Macro)
            for name, data := range self.Macros {
                data.parent = self
                cacheMacros[name] = data
            }
            return cacheMacros
        }
        // func (self *ProjectV9) createMacro(usage string, commands []string) Macro {
        //     return &MacroV9 {
        //         ConfigV9: *NewConfigV9(self,),
        //         Usage: usage,
        //         Actions: commands,
        //     }
        // }
        func (self *ProjectV9) getParent() Config {
            return self.parent
        }
        func (self *ProjectV9) getParentProject() Project {
            return self.parent
        }
        func (self *ProjectV9) setParentProject(project Project) {
            self.parent = project
        }

type MacroV9 struct {
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
    ConfigV9 `yaml:"inheritedValues,inline"`
}
        func (self *MacroV9) setParentProject(project Project) {
            self.ConfigV9.parent = project
        }
        func (self *MacroV9) getUsage() string {
            return self.Usage
        }
        func (self *MacroV9) getActions() []string {
            for i := range self.Actions {
                self.Actions[i] = os.ExpandEnv(self.Actions[i])
            }
            return self.Actions
        }
        func (self *MacroV9) getAliases() []string {
            return self.Aliases
        }
        func (self *MacroV9) getUsageText() string {
            return self.UsageText
        }
        func (self *MacroV9) getDescription() string {
            return self.Description
        }


func NewConfigV9(parent Config) *ConfigV9 {
    return &ConfigV9{
        Volumes: make(map[string]*VolumeV9),
        Devices: make(map[string]*DeviceV9),
        EnvironmentVariables: map[string]string{},
        parent: parent,
    }
}

func NewProjectV9(parent Project) *ProjectV9 {
    project := &ProjectV9 {
        SyntaxVersion: "8",
        Macros: make(map[string]*MacroV9),
        ConfigV9: *NewConfigV9(nil),
        parent: parent,
    }
    return project
}
