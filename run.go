package main

import (
    "github.com/Sirupsen/logrus"
)

func run() {
    p, err := loadProject()
    if err == nil {
        commands := p.Macros["run"]
        execInContainer(commands, p)
        // execInContainer([][]string{[]string{"run"}}, p)
    } else {
        logrus.Println("Could not load nut file")
    }
}
