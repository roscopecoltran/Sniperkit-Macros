package main

import (
	"os"
    "sort"

	"github.com/codegangsta/cli"
	log "github.com/Sirupsen/logrus"
    "fmt"
)

func main() {
	log.SetLevel(log.ErrorLevel)

	// try to parse the folder's configuration to add macros
	var macros []cli.Command
	project, err := loadProject()
	if err == nil {
        // Macros are stored in a random order.
        // But we want to display them in the same order everytime.
        // So sort the names of the macros.
        macroNamesOrdered := make([]string, 0, len(macros))
        for key, _ := range project.Macros {
            macroNamesOrdered = append(macroNamesOrdered, key)
        }
        sort.Strings(macroNamesOrdered)


		macros = make([]cli.Command, len(project.Macros))
		for index, name := range macroNamesOrdered {
            macro := project.Macros[name]
			log.Debug(name, ":", macro.Actions[0])

			nameClosure := name
            usage := "macro: "
            if macro.Usage == "" {
                usage += "undefined usage. define one with 'usage' property for this macro."
            } else {
                usage += macro.Usage
            }

			macros[index] = cli.Command{
				Name:  name,
				Usage: usage,
                Aliases: macro.Aliases,
                UsageText: macro.UsageText,
                Description: macro.Description,
				Action: func(c *cli.Context) {
					log.Debug(nameClosure, ":", macro.Actions[0])
					execMacro(&macro, project)
				},
				// Action: func(name string, macro []string) func(*cli.Context)  {
				// 	return func(c *cli.Context) {
				// 		log.Debug(name, ":", macro[0])
				// 		execMacro(macro, project)
				// 	}
				// }(name, macro),
			}
		}
	} else {
		macros = []cli.Command{}
	}

    initFlag := false
    statusFlag := false
    execFlag := ""

    app := cli.NewApp()
	app.Name = "nut"
	app.Version = "0.0.2 dev"
	app.Usage = "the development environment, containerized"
	app.EnableBashCompletion = true
    app.Flags = []cli.Flag {
        cli.BoolFlag{
            Name:        "init",
            Usage:       "initialize a nut project",
            Destination: &initFlag,
        },
        cli.BoolFlag{
            Name:  "status",
            Usage: "gives status of the dev env",
            Destination: &statusFlag,
        },
        cli.StringFlag{
            Name:  "exec",
            Usage: "execute a command in a container.",
            Destination: &execFlag,
        },
    }
    defaultAction := app.Action
    app.Action = func(c *cli.Context) {
        if statusFlag {
            status()
            return
        }
        if initFlag {
            initSubcommand(c)
            return
        }
        if execFlag != "" {
            if project != nil {
                exec(project, c, execFlag)
            } else {
                fmt.Println("Could not find nut configuration.")
            }
            return
        }
        defaultAction(c)
    }

	app.Commands = macros

	app.Run(os.Args)
}
