package main

import (
	"fmt"
)

func status() {
	// check is nut.yml exists at the current path
	if nutFileExistsAtPath(".") {
		fmt.Println("nut.yml is present")
	} else {
		fmt.Println("nut.yml not found")
	}
}
