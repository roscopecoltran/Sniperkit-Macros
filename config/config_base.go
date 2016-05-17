package config

import (
    "errors"
    containerFilepath "github.com/matthieudelaro/nut/container/filepath"
    Utils "github.com/matthieudelaro/nut/utils"
)

type VolumeBase struct {
}
        func (self *VolumeBase) fullHostPath(context Utils.Context) (string, error) {
            return "", errors.New("VolumeBase.fullHostPath() must be overloaded.")
        }
        func (self *VolumeBase) fullContainerPath(context Utils.Context) (string, error) {
            return "", errors.New("VolumeBase.fullContainerPath() must be overloaded.")
        }

type BaseEnvironmentBase struct {
}
        func (self *BaseEnvironmentBase) getFilePath() string{
            return ""
        }
        func (self *BaseEnvironmentBase) getGitHub() string{
            return ""
        }

type ConfigBase struct {
}
        func (self *ConfigBase) getDockerImage() string {
            return ""
        }
        func (self *ConfigBase) getProjectName() string {
            return ""
        }
        func (self *ConfigBase) getParent() Config {
            return nil
        }
        func (self *ConfigBase) getSyntaxVersion() string {
            return ""
        }
        func (self *ConfigBase) getBaseEnv() BaseEnvironment {
            return nil
        }
        func (self *ConfigBase) getWorkingDir() string {
            str, _ := containerFilepath.Abs(".")
            return str
        }
        func (self *ConfigBase) getVolumes() map[string]Volume {
            return make(map[string]Volume)
        }
        func (self *ConfigBase) getMacros() map[string]Macro {
            return make(map[string]Macro)
        }
        func (self *ConfigBase) getEnvironmentVariables() map[string]string {
            return make(map[string]string)
        }
        func (self *ConfigBase) getPorts() []string {
            return []string{}
        }
        func (self *ConfigBase) getEnableGui() (bool, bool) {
            return false, false
        }
        func (self *ConfigBase) getEnableNvidiaDevices() (bool, bool)  {
            return false, false
        }
        func (self *ConfigBase) getPrivileged() (bool, bool)  {
            return false, false
        }
        func (self *ConfigBase) getSecurityOpts() []string {
            return []string{}
        }

type ProjectBase struct {
}
        func (self *ProjectBase) getParentProject() Project {
            return nil
        }

type MacroBase struct {
}
        func (self *MacroBase) getUsage() string {
            return ""
        }
        func (self *MacroBase) getActions() []string {
            return []string{}
        }
        func (self *MacroBase) getAliases() []string {
            return []string{}
        }
        func (self *MacroBase) getUsageText() string {
            return ""
        }
        func (self *MacroBase) getDescription() string {
            return ""
        }
