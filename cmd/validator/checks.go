package main

import "fmt"

type checkDetails struct {
	reason      checkReason
	os          string
	arch        string
	othersCount int
}

func (cd checkDetails) String() string {
	if cd.othersCount > 0 {
		if cd.othersCount > 1 {
			return fmt.Sprintf("%s on platform %s-%s and %d others", cd.reason, cd.os, cd.arch, cd.othersCount)
		}
		return fmt.Sprintf("%s on platform %s-%s and one other", cd.reason, cd.os, cd.arch)
	}
	return fmt.Sprintf("%s on platform %s-%s", cd.reason, cd.os, cd.arch)
}

type checkReason int

const (
	reasonUpdate checkReason = iota
	reasonBundle
	reasonCommand
	reasonJar
)

func (cr checkReason) String() string {
	switch cr {
	case reasonUpdate:
		return "URL required for self-update"
	case reasonBundle:
		return "URL required for bundle-update"
	case reasonCommand:
		return "URL required for command binary"
	case reasonJar:
		return "URL required for Java application .jar"
	}
	panic(fmt.Sprintf("Unknown checkReason %d", cr))
}

type Check struct {
	URL     string
	Details checkDetails
	Error   error
}
