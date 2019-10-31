// +build mage

package main

import "os"

func init() {
	if err := os.Chdir(".."); err != nil {
		panic(err)
	}
}
