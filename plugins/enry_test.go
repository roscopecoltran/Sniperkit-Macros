package plugins

/*
import (
    "testing"
    "os"
    "io/ioutil"
)

func TestEnry(t *testing.T) {

	// Examples
	lang, safe := enry.GetLanguageByExtension("foo.go")
	fmt.Println(lang)
	// result: Go

	lang, safe := enry.GetLanguageByContent("foo.m", "<matlab-code>")
	fmt.Println(lang)
	// result: Matlab

	lang, safe := enry.GetLanguageByContent("bar.m", "<objective-c-code>")
	fmt.Println(lang)
	// result: Objective-C

	// all strategies together
	lang := enry.GetLanguage("foo.cpp", "<cpp-code>")

	// Languages
	langs := enry.GetLanguages("foo.h",  "<cpp-code>")
	// result: []string{"C++", "C"}

	langs := enry.GetLanguagesByExtension("foo.asc", "<content>", nil)
	// result: []string{"AGS Script", "AsciiDoc", "Public Key"}

	langs := enry.GetLanguagesByFilename("Gemfile", "<content>", []string{})
	// result: []string{"Ruby"}

	// Keywords
	keywords := enry.GetKeywords("foo.h",  "<cpp-code>")
	// result: []string{"C++", "C"}

	keywords := enry.GetKeywordsByExtension("foo.asc", "<content>", nil)
	// result: []string{"AGS Script", "AsciiDoc", "Public Key"}

	// Brands
	brands := enry.GetBrands("foo.h",  "<cpp-code>")
	// result: []string{"C++", "C"}

	brands := enry.GetBrandsByExtension("foo.asc", "<content>", nil)
	// result: []string{"AGS Script", "AsciiDoc", "Public Key"}

	// Authors
	authors := enry.GetAuthors("Luc Michalski",  "<cpp-code>")
	// result: []string{"C++", "C"}

	authors := enry.GetAuthorsByExtension("foo.asc", "Luc Michalski", nil)
	// result: []string{"AGS Script", "AsciiDoc", "Public Key"}

}

*/