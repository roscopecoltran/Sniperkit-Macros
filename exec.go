package main

import (
	"github.com/Sirupsen/logrus"
    "github.com/codegangsta/cli"
    "strings"
)

func exec(c *cli.Context) {
    p, err := loadProject()
    if err == nil {
        // commands := []string{"bash -c " + strings.Join(c.Args(), " ")}
        commands := []string{strings.Join(c.Args(), " ")}
        logrus.Println(commands)
        execInContainer(commands, p)
    } else {
        logrus.Println("Could not load nut file")
    }
}
