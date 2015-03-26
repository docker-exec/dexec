package main

import (
	"github.com/codegangsta/cli"
	"log"
	"os"
	"regexp"
)

func GetExtension(filename string) string {
	filenamePattern := regexp.MustCompile(`.*\.(.*)`)
	return filenamePattern.FindStringSubmatch(string(filename))[1]
}

func main() {
	app := cli.NewApp()
	app.Name = "dexec"
	app.Usage = "dexec"

	app.Action = func(c *cli.Context) {
		found := IsDockerPresent()
		running := IsDockerRunning()

		if !found {
			log.Fatal("Docker not found")
		} else if !running {
			log.Fatal("Docker not running")
		} else if len(c.Args()) == 0 {
			cli.ShowAppHelp(c)
		} else {
			RunDexecContainer("cpp", c.Args()[0], c.Args()[1:]...)
		}
	}

	app.Run(os.Args)
}
