package main

import (
    "github.com/Sirupsen/logrus"
    "github.com/codegangsta/cli"
)

func run(c *cli.Context) {
    p, err := loadProject()
    if err == nil {
        logrus.Println(p)
        commands := p.Macros["run"]
        execInContainer(commands, p)
        // execInContainer([][]string{[]string{"run"}}, p)
    } else {
        logrus.Println("Could not load nut file")
    }
}
