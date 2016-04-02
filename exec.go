package main

import (
    "github.com/codegangsta/cli"
)

func exec(project Project, c *cli.Context, command string) {
    execInContainer([]string{command}, project)
}
