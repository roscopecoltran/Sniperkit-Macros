package main

import (
    // "errors"
    log "github.com/Sirupsen/logrus"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "testing"
    "reflect"
    "strconv"
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

func TestParsingV5(t *testing.T) {
    log.SetLevel(log.DebugLevel)

    type Tuple struct {
        file string
        env map[string]string
        ports []string
    }

    nutFiles := []Tuple{}

    nutFiles = append(nutFiles,
Tuple{ file:`
syntax_version: "5"
based_on:
  docker_image: golang:1.6
`,
ports: []string{},
},
Tuple{ file:`
syntax_version: "5"
based_on:
  docker_image: golang:1.6
environment:
  A: 1
  B: 2
`,
env: map[string]string{
    "A": "1",
    "B": "2",
}, ports: []string{},
},
Tuple{ file:`
syntax_version: "5"
based_on:
  docker_image: golang:1.6
environment:
  A: 1
  B:
ports:
  - "3000:3000"
  - 100:100
`,
env: map[string]string{
    "A": "1",
    "B": "",
}, ports: []string{
    "3000:3000",
    "100:100",
},
},
Tuple{ file:`
syntax_version: "5"
based_on:
  docker_image: golang:1.6
environment:
  A: 1
  B:
ports:
  - "3000:3000"
`,
env: map[string]string{
    "A": "1",
    "B": "",
}, ports: []string{
    "3000:3000",
},
},
)

    for index, tuple := range nutFiles {
        byteArray := []byte(tuple.file)
        project := NewProjectV5()
        assert.Nil(t, project.fromYAML(byteArray))

        assert.Equal(t, len(reflect.ValueOf(tuple.env).MapKeys()),
            len(reflect.ValueOf(project.getEnvironmentVariables()).MapKeys()),
            "Error with tuple " + strconv.Itoa(index) + ": not same keys")

        for name, value := range project.getEnvironmentVariables() {
            assert.Equal(t, value, tuple.env[name],
                "Error with tuple " + strconv.Itoa(index) + ": not same " + name)
            // execInContainer([]string{"echo $" + name}, project)
            // TODO: automate this test
        }

        require.Equal(t, len(tuple.ports), len(project.getPorts()),
            "Error with tuple " + strconv.Itoa(index) + ": not same port quantity")
        for key, value := range project.getPorts() {
            assert.Equal(t, value, tuple.ports[key],
                "Error with tuple " + strconv.Itoa(index) +
                ": not same port " + strconv.Itoa(key))
        }
    }

    // test nginx
    // TODO: automate this test
//     nutFile := `
// syntax_version: "5"
// based_on:
//   docker_image: nginx
// ports:
// #  - "80:80"  # works
//   - "80"  # works
// `
//     project := NewProjectV5()
//     byteArray := []byte(nutFile)
//     assert.Nil(t, project.fromYAML(byteArray))
//     execInContainer([]string{}, project)

//     log.Error("end")
}
