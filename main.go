package main

import (
	"os"

	"github.com/codegangsta/cli"
	log "github.com/Sirupsen/logrus"
)

func main() {
	log.SetLevel(log.ErrorLevel)

	// try to parse the folder's configuration to add macros
	var macros []cli.Command
	proj, err := loadProject()
	if err == nil {
		macros = make([]cli.Command, len(proj.Macros))
		counter := 0
		for name, macro := range proj.Macros {
			log.Debug(name, ":", macro[0])

			nameClosure := name
			macroClosure := macro

			macros[counter] = cli.Command{
				Name:  name,
				Usage: name + " is a macro defined in configuration.",
				Action: func(c *cli.Context) {
					log.Debug(nameClosure, ":", macroClosure[0])
					execMacro(macroClosure, proj)
				},
				// Action: func(name string, macro []string) func(*cli.Context)  {
				// 	return func(c *cli.Context) {
				// 		log.Debug(name, ":", macro[0])
				// 		execMacro(macro, proj)
				// 	}
				// }(name, macro),
			}
			counter += 1
		}
	} else {
		macros = []cli.Command{}
	}


	app := cli.NewApp()
	app.Name = "nut"
	app.Version = "0.0.1 dev"
	app.Usage = "the development environment, containerized"
	app.EnableBashCompletion = true
	// define nut subcommands
	nutCommands := []cli.Command{
		{
			Name:  "init",
			Usage: "init a nut project",
			Action: func(c *cli.Context) {
				initSubcommand(c)
			},
		},
		{
			Name:  "status",
			Usage: "gives status of the dev env",
			Action: func(c *cli.Context) {
				status()
			},
		},
		{
			Name:  "exec",
			Usage: "execute a command in a container. Surround the whole command with simple quotes to obtain expected result.",
			Action: func(c *cli.Context) {
				exec(c)
			},
		},
	}

	log.Debug("debug")
	for index, command := range macros {
		log.Debug(index, command.Name)
	}
	// merge built-in nut commands with project-defined macros
	for _, command := range nutCommands {
		// find out whether there is or not a macro with the same name
		found := false
		for counter := len(macros) - 1; counter >= 0 && found == false; counter-- {
			log.Debug(counter)
			if macros[counter].Name == command.Name {
				found = true
			}
		}
		// if there isn't any macro with the same name, add the command
		if found == false {
			macros = append(macros, command)
		}
	}

	app.Commands = macros

	app.Run(os.Args)
}
