package main

import (
	"github.com/Sirupsen/logrus"
)

func build() {
    p, err := loadProject()
    if err == nil {
        commands := p.Macros["build"]
        execInContainer(commands, p)
    } else {
        logrus.Println("Could not load nut file")
    }
}
