package main

import (
  "os"
	"os/exec"
  "github.com/codegangsta/cli"
)

func main() {
  cli.NewApp().Run(os.Args)
}