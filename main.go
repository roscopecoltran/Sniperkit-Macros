package main

import (
	"os"
    "sort"
	"github.com/codegangsta/cli"
	log "github.com/Sirupsen/logrus"
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
        projectMacros := project.getMacros()
        macroNamesOrdered := make([]string, 0, len(projectMacros))
        for key, _ := range projectMacros {
            macroNamesOrdered = append(macroNamesOrdered, key)
        }
        sort.Strings(macroNamesOrdered)


		macros = make([]cli.Command, len(projectMacros))
		for index, name := range macroNamesOrdered {
            macro := projectMacros[name]

			nameClosure := name
            usage := "macro: "
            if macro.getUsage() == "" {
                usage += "undefined usage. define one with 'usage' property for this macro."
            } else {
                usage += macro.getUsage()
            }

			macros[index] = cli.Command{
				Name:  nameClosure,
				Usage: usage,
                Aliases: macro.getAliases(),
                UsageText: macro.getUsageText(),
                Description: macro.getDescription(),
				Action: func(c *cli.Context) {
					execMacro(macro, project)
				},
			}
		}
	} else {
		macros = []cli.Command{}
	}

    initFlag := false
    statusFlag := false
    logsFlag := false
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
            Name:        "logs",
            Usage:       "display log messages. Useful for contributors and to report an issue",
            Destination: &logsFlag,
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
        if logsFlag {
            log.SetLevel(log.DebugLevel)
        }
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
                log.Error("Could not find nut configuration.")
            }
            return
        }
        defaultAction(c)
    }

	app.Commands = macros

    app.Run(os.Args)
}
