package plugins

import (
	"github.com/michaelsauter/crane/crane"
)

func Crane() {
	crane.checkDockerClient()
}