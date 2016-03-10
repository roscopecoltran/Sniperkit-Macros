package main

import (
	"fmt"
)

func initProject() {
	// check is nut.yml exists at the current path
	if nutFileExistsAtPath(".") {
		fmt.Println("a nut.yml file already exists")
	} else {
		// create a nut.yml at the current path
		// err := ioutil.WriteFile("./nut.yml", data, perm)
		fmt.Println("project created (NOT IMPLEMENTED YET)")
	}
}
