package main

import (
    // "errors"
    // log "github.com/Sirupsen/logrus"
    "testing"
)

type projectPair struct {
    p1 Project
    p2 Project
    res interface{}
}

func makeProject(name string, workingDir string, enableGui string) Project {
    p := NewProjectV4()
    p.Mount["main"] = []string{".", "/go/src/project"}
    p.Macros["build"] = &MacroV4{
        Usage: "build the project",
        Actions: []string{"go build -o nut"},
    }
    p.ProjectName = name
    p.WorkingDir = workingDir
    p.EnableGUI = enableGui
    return p
}

func makeChild(name string, workingDir string, enableGui string) Project {
    return makeProject(name, workingDir, enableGui)
}
func makeParent(name string, workingDir string, enableGui string) Project {
    return makeProject(name, workingDir, enableGui)
}

func strToBool(str string) bool {
    if str == "true" {
        return true
    } else {
        return false
    }
}

func TestInheritanceEnableGui(t *testing.T) {
    possible := []string{"true", "false", "", "junk"}
    for _, vp := range possible {
        for _, vc := range possible {

            p1 := makeChild("child", "child", vc)
            p2 := makeParent("parent", "parent", vp)

            p1.setParentProject(p2)

            res := p1.getEnableGui()
            var expectation bool
            if vc == "" {
                expectation = strToBool(vp)
            } else {
                expectation = strToBool(vc)
            }

            if res != expectation {
                t.Error(
                    "For", vc, vp,
                    "expected", expectation,
                    "got", res,
                )
            }
        }
    }
}

func TestInheritanceDockerImage(t *testing.T) {
    possible := []string{"", "golang:1.5", "golang:1.6"}
    for _, vp := range possible {
        for _, vc := range possible {
            pc := NewProjectV4()
            pc.Base.DockerImage = vc
            pp := NewProjectV4()
            pp.Base.DockerImage = vp

            pc.setParentProject(pp)

            res, err := pc.getBaseEnv().getDockerImageName()
            expectation := "undefined"
            expectError := false
            if vc == "" {
                expectation = vp
                if vp == "" {
                    expectError = true
                }
            } else {
                expectation = vc
            }

            if res != expectation || (err != nil) != expectError {
                t.Error(
                    "For child=", vc, "parent=", vp,
                    "expected", expectation, expectError, "error",
                    "got", res, err, "error",
                )
            }
        }
    }
}

func TestInheritanceWorkingDir(t *testing.T) {
    possible := []string{"here", "", "or here"}
    for _, vp := range possible {
        for _, vc := range possible {

            p1 := makeChild("child", vc, "child")
            p2 := makeParent("parent", vp, "parent")

            p1.setParentProject(p2)

            res := p1.getWorkingDir()
            var expectation string
            if vc == "" {
                expectation = vp
            } else {
                expectation = vc
            }

            if res != expectation {
                t.Error(
                    "For", vc, vp,
                    "expected", expectation,
                    "got", res,
                )
            }
        }
    }
}

func TestInheritanceName(t *testing.T) {
    possible := []string{"name1", "", "name2"}
    for _, vp := range possible {
        for _, vc := range possible {

            p1 := makeChild(vc, "child", "child")
            p2 := makeParent(vp, "parent", "parent")

            p1.setParentProject(p2)

            res := p1.getName()
            var expectation string
            if vc == "" {
                expectation = vp
            } else {
                expectation = vc
            }

            if res != expectation {
                t.Error(
                    "For", vc, vp,
                    "expected", expectation,
                    "got", res,
                )
            }
        }
    }
}

func makeProjectMacros(macros []string, usage string) Project {
    p := NewProjectV4()
    for _, name := range macros {
        p.Macros[name] = &MacroV4{
            Usage: usage,
            Actions: []string{name},
        }
    }
    return p
}

func macroInList(name string, list []string) bool {
    for _, b := range list {
        if b == name {
            return true
        }
    }
    return false
}

func TestInheritanceMacros(t *testing.T) {
    possible := [][]string{[]string{}, []string{"run"}, []string{"build"}, []string{"build", "run"}}
    for _, vp := range possible {
        for _, vc := range possible {

            pc := makeProjectMacros(vc, "child")
            pp := makeProjectMacros(vp, "parent")

            pc.setParentProject(pp)

            nameOccurences := make(map[string]int)
            for _, name := range append(vc, vp...) {
                nameOccurences[name] += 1
            }

            macrosc := pc.getMacros()
            macrosp := pp.getMacros()

            for name, count := range nameOccurences {
                if count == 2 {
                    // it must be in parent and in child
                    // and we expect the one of the child
                    if macrosc[name] == nil ||
                       macrosp[name] == nil ||
                       macrosc[name].getUsage() != "child" {
                        t.Error("Macro", name, "defined in parent and in child",
                            "expected to be overriden by the child")
                    }
                } else if count == 1 {
                    // it comes either from parent or from child
                    if macrosp[name] != nil {
                        if macrosp[name].getUsage() != "parent" {
                            t.Error("Macro", name, "defined in parent",
                                "expected to be overriden by the child")
                        }
                    } else if macrosc[name] != nil{
                        if macrosc[name].getUsage() != "child" {
                            t.Error("Macro", name, "defined in child",
                                "expected to be overriden by the child")
                        }
                    } else {
                        t.Error("Macro", name, "has been ignored.")
                    }
                }
            }
        }
    }
}
