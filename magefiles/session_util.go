// +build mage

package main

import "github.com/codeskyblue/go-sh"

func s() *sh.Session {
	session := sh.NewSession()
	session.ShowCMD = true
	return session
}

func run(command *sh.Session) {
	if err := command.Run(); err != nil {
		panic(err)
	}
}
