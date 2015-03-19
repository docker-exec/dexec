package dexec

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

func getRawDockerVersion() string {
  out, err := exec.Command("docker", "-v").Output()
  if err != nil {
    panic(err.Error())
  } else {
    return string(out)
  }
}

func extractDockerVersion(rawVersion string) string {
  dockerVersionPattern := regexp.MustCompile(`^Docker version (\d+\.\d+\.\d+), build [a-z0-9]{7}`)

  if dockerVersionPattern.MatchString(rawVersion) {
      return dockerVersionPattern.FindStringSubmatch(rawVersion)[1]
  } else {
    panic("Did not match Docker version string")
  }
}

func isDockerPresent() bool {
  defer func() {
      if r := recover(); r != nil {
          fmt.Println("Recovered in f", r)
      }
  }()
  present := false

  present = extractDockerVersion(getRawDockerVersion()) != ""

  return present
}

func is_docker_running() (bool, string) {
    return true, "Running"
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
        found := isDockerPresent()
        running, msg := is_docker_running()
        if (!found) {
            log.Fatal(msg)
        } else if(!running) {
            log.Fatal(msg)
        } else if len(c.Args()) == 0 {
            cli.ShowAppHelp(c)
        } else {
            run_dexec_container(getExtension(c.Args()[0]), c.Args()[0], c.Args()[1:]...)
        }
    }

    app.Run(os.Args)
}