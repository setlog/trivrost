// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"

	sh "github.com/codeskyblue/go-sh"
)

func Generate() error {
	s := sh.NewSession()
	s.SetEnv("GO111MODULE", "off")
	if err := s.Command("go", "get", "-u", "github.com/josephspurrier/goversioninfo/cmd/goversioninfo").Run(); err != nil {
		return fmt.Errorf("could not go get goversioninfo: %w", err)
	}
	s = sh.NewSession()
	s.SetDir(filepath.Dir(wd()))
	return s.Command("go", "generate", "-installsuffix", "_separate", "github.com/setlog/trivrost/cmd/launcher").Run()
}

func wd() string {
	workingDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return workingDirectory
}
