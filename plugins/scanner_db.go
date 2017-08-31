package plugins

import (
	"github.com/src-d/enry"
)

func Enry() {

	langs := enry.GetLanguages("foo.h",  "<cpp-code>")
	// result: []string{"C++", "C"}

	keywords := enry.GetLanguages("foo.h",  "<cpp-code>")
	// result: []string{"C++", "C"}

	authors := enry.GetLanguages("foo.h",  "<cpp-code>")
	// result: []string{"C++", "C"}

}