package plugins

import (
	"fmt"
	"github.com/hoop33/limo/config"
)

func Limo() {
	fmt.Sprintf(" name: %s", config.ProgramName)
	fmt.Sprintf(" version: %s", config.Version)
}


