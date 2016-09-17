package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/docker-exec/dexec/cli"
	"github.com/docker-exec/dexec/dexec"
	"github.com/docker-exec/dexec/util"
	"github.com/fsouza/go-dockerclient"
)

// RunDexecContainer runs an anonymous Docker container with a Docker Exec
// image, mounting the specified sources and includes and passing the
// list of sources and arguments to the entrypoint.
func RunDexecContainer(cliParser cli.CLI) int {
	options := cliParser.Options

	shouldClean := len(options[cli.CleanFlag]) > 0
	updateImage := len(options[cli.UpdateFlag]) > 0

	client, err := docker.NewClientFromEnv()

	if err != nil {
		log.Fatal(err)
	}

	if shouldClean {
		images, err := client.ListImages(docker.ListImagesOptions{
			All: true,
		})
		if err != nil {
			log.Fatal(err)
		}
		for _, image := range images {
			for _, tag := range image.RepoTags {
				repoRegex := regexp.MustCompile("^dexec/lang-[^:\\s]+(:.+)?$")
				if match := repoRegex.MatchString(tag); match {
					client.RemoveImage(image.ID)
				}
			}
		}
	}

	dexecImage, err := dexec.ImageFromOptions(options)
	if err != nil {
		log.Fatal(err)
	}

	dockerImage := fmt.Sprintf("%s:%s", dexecImage.Image, dexecImage.Version)

	if err = dexec.FetchImage(
		dexecImage.Image,
		dexecImage.Version,
		updateImage,
		client); err != nil {
		log.Fatal(err)
	}

	var sourceBasenames []string
	for _, source := range options[cli.Source] {
		basename, _ := dexec.ExtractBasenameAndPermission(source)
		sourceBasenames = append(sourceBasenames, []string{basename}...)
	}

	entrypointArgs := util.JoinStringSlices(
		sourceBasenames,
		util.AddPrefix(options[cli.BuildArg], "-b"),
		util.AddPrefix(options[cli.Arg], "-a"),
	)

	container, err := client.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:     dockerImage,
			Cmd:       entrypointArgs,
			StdinOnce: true,
			OpenStdin: true,
		},
		HostConfig: &docker.HostConfig{
			Binds: dexec.BuildVolumeArgs(
				util.RetrievePath(options[cli.TargetDir]),
				append(options[cli.Source], options[cli.Include]...)),
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = client.RemoveContainer(docker.RemoveContainerOptions{
			ID: container.ID,
		}); err != nil {
			log.Fatal(err)
		}
	}()

	if err = client.StartContainer(container.ID, &docker.HostConfig{}); err != nil {
		log.Fatal(err)
	}

	go func() {
		if err = client.AttachToContainer(docker.AttachToContainerOptions{
			Container:   container.ID,
			InputStream: os.Stdin,
			Stream:      true,
			Stdin:       true,
		}); err != nil {
			log.Fatal(err)
		}
	}()

	code, err := client.WaitContainer(container.ID)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Logs(docker.LogsOptions{
		Container:    container.ID,
		Stdout:       true,
		Stderr:       true,
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
	})

	if err != nil {
		log.Fatal(err)
	}

	return code
}

func validate(cliParser cli.CLI) bool {
	options := cliParser.Options

	hasVersionFlag := len(options[cli.VersionFlag]) == 1
	hasSources := len(options[cli.Source]) > 0
	shouldClean := len(options[cli.CleanFlag]) > 0

	if hasSources || shouldClean {
		return true
	}

	if hasVersionFlag {
		cli.DisplayVersion(cliParser.Filename)
		return false
	}

	cli.DisplayHelp(cliParser.Filename)
	return false
}

func validateDocker() error {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	ping := make(chan error, 1)
	go func() {
		ping <- client.Ping()
	}()

	select {
	case err := <-ping:
		return err
	case <-time.After(5 * time.Second):
		return fmt.Errorf("Request to Docker host timed out")
	}
}

func main() {
	cliParser := cli.ParseOsArgs(os.Args)

	if validate(cliParser) {
		if err := validateDocker(); err != nil {
			log.Fatal(err)
		} else {
			os.Exit(RunDexecContainer(cliParser))
		}
	}
}
