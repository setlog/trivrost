// +build mage

package main

func Generate() {
	run(s().SetEnv("GO111MODULE", "off").Command("go", "get", "-u", "github.com/josephspurrier/goversioninfo/cmd/goversioninfo"))
	run(s().Command("go", "generate", "-installsuffix", "_separate", "github.com/setlog/trivrost/cmd/launcher"))
}
