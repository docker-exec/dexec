package main

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "github.com/codegangsta/cli"
)

// ${docker_cmd} run -t --rm \
//     -w ${work_dir} \
//     -v $(pwd -P)/${source_file}:${work_dir}/${source_file} \
//     ${docker_image} ${source_file} ${entrypoint_args}

func run_anonymous_container(image string) {
    out, err := exec.Command("docker", "run", "-t", "--rm", image, "echo", "test").Output()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("output:\n%s\n", out)
}

func main() {
    app := cli.NewApp()
    app.Name = "dexec"
    app.Usage = "dexec"
    app.Action = func(c *cli.Context) {
        run_anonymous_container("ubuntu")
        out, err := exec.Command("docker", "ps").Output()
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("Running Docker processes:\n%s\n", out)
    }

    app.Run(os.Args)
}