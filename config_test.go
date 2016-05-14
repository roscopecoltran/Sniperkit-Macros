package main

// import (
//     log "github.com/Sirupsen/logrus"
//     Config "github.com/matthieudelaro/nut/config"
//     "gopkg.in/yaml.v2"
//     "testing"
// )


// func TestFromNutPackage(t *testing.T) {
//     log.SetLevel(log.DebugLevel)


//     log.Debug("------Tests of main.go")

//     var conf Config.Config
//     // conf = &Config.ConfigBase{}
//     // conf = &Config.ProjectBase{}
//     // conf = &Config.ProjectBase{}
//     conf = Config.NewConfigV5(
//         "V5/",
//         []string{"V5port"},
//         nil,
//     )
//     log.Debug("OK ", conf)

//     // var proj Project
//     // proj = &ProjectBase{}
//     proj := Config.NewProjectV5(map[string]*Config.MacroV5{
//         "run": Config.NewSimpleMacroV5("run, from project"),
//         "build": Config.NewSimpleMacroV5("build, from project"),
//         }, *Config.NewConfigV5(
//             "ProjectV5.WorkingDir",
//             []string{"ProjectV5.Ports"},
//             nil,
//         ),)
//     log.Debug("OK ", proj)

//     var macro Config.Macro
//     // macro = &Config.MacroBase{}
//     macro = &Config.MacroV5{
//         Usage: "UsageV5",
//         ConfigV5: *Config.NewConfigV5(
//             "MacroV5.WorkingDir",
//             []string{"MacroV5.Ports"},
//             proj,
//         ),
//     }
//     log.Debug("OK ", macro)

//     macros := Config.GetMacros(macro)
//     log.Debug("OK merge macros ", macros)

//     ports := Config.GetPorts(macro)
//     log.Debug("OK merge ports ", ports)

//     macroExec := proj.CreateMacro([]string{"echo hello"})
//     log.Debug("OK macroExec ", macroExec)
//     log.Debug("proj before YAML ", proj)

//     bytes, err := yaml.Marshal(proj)
//     log.Debug("proj ", err, string(bytes))

//     input := `ports:
// - ProjectV5.Ports
// macros:
//   build:
//     usage: build, from project
//   run:
//     usage: run, from project`
//     err = yaml.Unmarshal([]byte(input), proj)
//     log.Debug("proj ", err, proj)
//     bytes, err = yaml.Marshal(proj)
//     log.Debug("proj ", err, string(bytes))
//     err = yaml.Unmarshal([]byte(input), proj)
//     log.Debug("proj ", err, proj)
// }

// func TestFromConfigPackage(t *testing.T) {
//     log.SetLevel(log.DebugLevel)
//     Config.Tests()
// }
