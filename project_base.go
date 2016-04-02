package main

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
    "path/filepath"
)

////// README
// Nut must remain backward compatible with respect to nut configuration files.
// It must also be easy to add new features in nut files, without issues of
// performance to parse the file, and without applying modifications to all
// older versions to insure backward compatibility.
//
// The chosen solution is that new features should be defined in the interface,
// and default behavior/feature should be defined in base class (to be
// accessible) to all syntaxes version. Then override this default behavior in
// the new syntax version.
//
//  Interfaces Names| Base Class Names    | Version 2         | Version 3         | ...
//  ----------------+---------------------+-------------------+-------------------+----
//  Macro           | MacroBase           | MacroV2           | MacroV3           | ...
//  MountingPoint   | MountingPointBase   | MountingPointV2   | MountingPointV3   | ...
//  BaseEnvironment | BaseEnvironmentBase | BaseEnvironmentV2 | BaseEnvironmentV3 | ...
//  Project         | ProjectBase         | ProjectV2         | ProjectV3         | ...
//
// For each new syntax, copy/paste a project_vXXX.go file (ie the latest one)
// and change implementation to suit new requirements.
// For any new feature, create a virtual method in the interface Project
// (or its components Macro, MountingPoint, BaseEnvironment, etc) and
// implement a method in ProjectBase class (or its components MacroBase, etc)
// to define a default behavior for older versions.
//
// Old syntax versions should never be modified. Only the interface, the
// base class, and the new syntax should be updated.
//
// Note: When working with structs and embedding, everything is STATICALLY LINKED. All references are resolved at compile time.
// See https://github.com/luciotato/golang-notes/blob/master/OOP.md for more details.


////// Interfaces
/// Define what a project (and its components) should be able to do
type Macro interface {
    getUsage() string
    getActions() []string
    getAliases() []string
    getUsageText() string
    getDescription() string
}

type MountingPoint interface {
    fullHostPath() (string, error)
    fullContainerPath() (string, error)
}

type BaseEnvironment interface {
    getDockerImageName() (string, error)
}

type Project interface {
    getSyntaxVersion() string
    getName() string
    getBaseEnv() BaseEnvironment
    getWorkingDir() string
    getMountingPoints() map[string]MountingPoint
    getMacros() map[string]Macro
    getEnableGui() bool
    toYAML() string
    fromYAML(bytes []byte) error
}

/////// Base classes
/// Placeholder for basic behaviors, and default behaviors for retro-compatibility.
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

type MountingPointBase struct {
}
        func (self *MountingPointBase) fullHostPath() (string, error) {
            return "", errors.New("MountingPointBase.fullHostPath() must be overloaded.")
        }
        func (self *MountingPointBase) fullContainerPath() (string, error) {
            return "", errors.New("MountingPointBase.fullContainerPath() must be overloaded.")
        }

type BaseEnvironmentBase struct {
}
        func (self *BaseEnvironmentBase) getDockerImageName() (string, error) {
            return "", errors.New("BaseEnvironmentBase.getDockerImageName() must be overloaded.")
        }

type ProjectBase struct {
}
        func (self *ProjectBase) getSyntaxVersion() string {
            return "1"
        }
        func (self *ProjectBase) getName() string {
        	return "nut"
        }
        func (self *ProjectBase) getBaseEnv() BaseEnvironment {
            return &BaseEnvironmentBase{}
        }
        func (self *ProjectBase) getWorkingDir() string {
        	str, _ := filepath.Abs(".")
            return str
        }
        func (self *ProjectBase) getMountingPoints() map[string]MountingPoint {
            return make(map[string]MountingPoint)
        }
        func (self *ProjectBase) getMacros() map[string]Macro {
        	return make(map[string]Macro)
        }
        func (self *ProjectBase) toYAML() string {
            d, err := yaml.Marshal(&self)
            if err != nil {
                log.Fatalf("error: %v", err)
            }
            return string(d)
        }
	    func (self *ProjectBase) fromYAML(bytes []byte) error {
	    	return errors.New("ProjectBase.fromYAML() must be overloaded.")
	    }
        func (self *ProjectBase) getEnableGui() bool {
        	return false
        }

func loadProject() (Project, error) {
	return parseNutFileAtPath("nut.yml")
}

func ProjectFromBytes(bytes []byte) (Project, error) {
	var err error
	for _, syntax := range getSyntaxes() {
		version := syntax.getSyntaxVersion()
		err = syntax.fromYAML(bytes)
		if err == nil {
			return syntax, nil
		} else {
			log.Debug("Not syntax ", version, ": ", err)
		}
	}
	return nil, err
}

func getSyntaxes() []Project {
	return []Project{
		NewProjectV3(),
		NewProjectV2(),
	}
}
