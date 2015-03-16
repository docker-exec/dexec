package main

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "github.com/codegangsta/cli"
)

func main() {
    app := cli.NewApp()
    app.Name = "dexec"
    app.Usage = "dexec"
    app.Action = func(c *cli.Context) {
        out, err := exec.Command("docker", "ps").Output()
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("Running Docker processes:\n%s\n", out)
    }

    app.Run(os.Args)
}