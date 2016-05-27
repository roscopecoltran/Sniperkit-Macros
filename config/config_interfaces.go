package config

import (
    Utils "github.com/matthieudelaro/nut/utils"
)


type Volume interface {
    getHostPath() string
    getContainerPath() string
    getOptions() (string)
    fullHostPath(context Utils.Context) (string, error)
    fullContainerPath(context Utils.Context) (string, error)
}
type BaseEnvironment interface {
    getFilePath() string
    getGitHub() string
}


type Config interface {
    getDockerImage() string
    getProjectName() string
    getNetworkMode() string
    getUTSMode() string
    getParent() Config
    getSyntaxVersion() string
    getBaseEnv() BaseEnvironment
    getWorkingDir() string
    getVolumes() map[string]Volume
    getMacros() map[string]Macro
    getEnvironmentVariables() map[string]string
    getDevices() map[string]Volume
    getPorts() []string
    getEnableGui() (bool, bool)
    getEnableNvidiaDevices() (bool, bool)
    getPrivileged() (bool, bool)
    getDetached() (bool, bool)
    getSecurityOpts() []string
}
type Project interface { // extends Config interface
    // pure Project methods
    // createMacro(name string, commands []string) Macro
    setParentProject(project Project)

    // Config methods
    getDockerImage() string
    getProjectName() string
    getNetworkMode() string
    getUTSMode() string
    getParent() Config
    getParentProject() Project
    getSyntaxVersion() string
    getBaseEnv() BaseEnvironment
    getWorkingDir() string
    getVolumes() map[string]Volume
    getMacros() map[string]Macro
    getEnvironmentVariables() map[string]string
    getDevices() map[string]Volume
    getPorts() []string
    getEnableGui() (bool, bool)
    getEnableNvidiaDevices() (bool, bool)
    getPrivileged() (bool, bool)
    getDetached() (bool, bool)
    getSecurityOpts() []string

}
type Macro interface { // extends Config interface
    // pure Macro methods
    setParentProject(project Project)
    getUsage() string
    getActions() []string
    getAliases() []string
    getUsageText() string
    getDescription() string

    // Config methods
    getDockerImage() string
    getProjectName() string
    getNetworkMode() string
    getUTSMode() string
    getParent() Config
    getSyntaxVersion() string
    getBaseEnv() BaseEnvironment
    getWorkingDir() string
    getVolumes() map[string]Volume
    getMacros() map[string]Macro
    getEnvironmentVariables() map[string]string
    getDevices() map[string]Volume
    getPorts() []string
    getEnableGui() (bool, bool)
    getEnableNvidiaDevices() (bool, bool)
    getPrivileged() (bool, bool)
    getDetached() (bool, bool)
    getSecurityOpts() []string
}

