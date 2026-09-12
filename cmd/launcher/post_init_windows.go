package main

import "golang.org/x/sys/windows"

func postInit() {
	err := windows.SetDllDirectory("")
	if err != nil {
		panic(err)
	}
}
