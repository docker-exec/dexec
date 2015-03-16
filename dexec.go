package main

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "regexp"
    "github.com/codegangsta/cli"
    "code.google.com/p/go-uuid/uuid"
)

func is_docker_present() (bool, string) {
    dockerVersionPattern := regexp.MustCompile(`^Docker version (\d+\.\d+\.\d+), build [a-z0-9]{7}`)

    out, err := exec.Command("docker", "-v").Output()
    if err != nil {
        return false, err.Error()
    } else if !dockerVersionPattern.Match(out) {
        return false, "Did not match Docker version string"
    } else {
        return true, dockerVersionPattern.FindStringSubmatch(string(out))[1]
    }
}

func is_docker_running() bool {
    return true
}

func run_anonymous_container(args ...string) {
    newArgs := append([]string{"run", "-t", "--rm"}, args...)
    out := exec.Command("docker", newArgs...)
    out.Stdin = os.Stdin
    out.Stdout = os.Stdout
    out.Stderr = os.Stderr
    out.Run()
}

func run_dexec_container(language string, sourcefile string, entrypointargs ...string) {
    workdir := fmt.Sprintf("/tmp/%s", uuid.New())
    abssourcefile, _ := filepath.Abs(sourcefile)

    run_anonymous_container(
        append(
            []string {
                "-w", workdir,
                "-v", fmt.Sprintf("%s:%s/%s", abssourcefile, workdir, sourcefile),
                fmt.Sprintf("dexec/%s", language),
                fmt.Sprintf("%s/%s", workdir, sourcefile)},
            entrypointargs...
        )...
    )
}

func getExtension(filename string) string {
    filenamePattern := regexp.MustCompile(`.*\.(.*)`)
    return filenamePattern.FindStringSubmatch(string(filename))[1]
}

func main() {
    app := cli.NewApp()
    app.Name = "dexec"
    app.Usage = "dexec"

    app.Action = func(c *cli.Context) {
        found, msg := is_docker_present()
        if (!found) {
            log.Fatal(msg)
        } else if len(c.Args()) == 0 {
            cli.ShowAppHelp(c)
        } else {
            run_dexec_container(getExtension(c.Args()[0]), c.Args()[0], c.Args()[1:]...)
        }
    }

    app.Run(os.Args)
}