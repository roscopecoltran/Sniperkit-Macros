package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// parseNutFileAtPath reads and parses the nut file at the given path
func parseNutFileAtPath(nutFilePath string) (*Project, error) {
	// check file exists
	exists, err := fileExists(nutFilePath)
	if err != nil {
		return nil, err
	}
	if exists == false {
		return nil, errors.New("nut file not found")
	}
	// file exists
	bytes, err := ioutil.ReadFile(nutFilePath)
	if err != nil {
		return nil, err
	}

	var project *Project = &Project{}
	err = yaml.Unmarshal(bytes, &project)
	return project, err
}

// nutFileExistsAtPath return whether a nut.yml file exists at the given path
func nutFileExistsAtPath(parentPath string) bool {
	nutFilePath := path.Join(parentPath, "nut.yml")
	exists, err := fileExists(nutFilePath)
	if err != nil {
		logrus.Fatal(err)
	}
	return exists
}

// fileExists returns whether the given file or directory exists or not
func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
