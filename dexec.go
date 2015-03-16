package main

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "github.com/codegangsta/cli"
)

func is_docker_present() bool {
    return true
}

func is_docker_running() bool {
    return true
}

func run_anonymous_container(image string) {
    out, err := exec.Command("docker", "run", "-t", "--rm", image, "echo", "test").Output()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s: %s", image, out)
}

func main() {
    app := cli.NewApp()
    app.Name = "dexec"
    app.Usage = "dexec"
    app.Flags = []cli.Flag {
      cli.StringFlag{
        Name: "srcfile",
        Usage: "Source file",
      },
    }
    app.Action = func(c *cli.Context) {
        if len(c.Args()) == 0 {
            cli.ShowAppHelp(c)
        } else {
            run_anonymous_container("ubuntu")
            run_anonymous_container("debian")
        }
    }

    app.Run(os.Args)
}