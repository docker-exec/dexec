package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

const timeoutStatusCode = 124

// RunDexecContainer runs an anonymous Docker container with a Docker Exec
// image, mounting the specified sources and includes and passing the
// list of sources and arguments to the entrypoint.
func RunDexecContainer(cliParser CLI) int {
	options := cliParser.Options

	shouldClean := len(options[CleanFlag]) > 0
	updateImage := len(options[UpdateFlag]) > 0

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

	dexecImage, err := ImageFromOptions(options)
	if err != nil {
		log.Fatal(err)
	}

	dockerImage := fmt.Sprintf("%s:%s", dexecImage.Image, dexecImage.Version)

	if err = FetchImage(
		dexecImage.Image,
		dexecImage.Version,
		updateImage,
		client); err != nil {
		log.Fatal(err)
	}

	var sourceBasenames []string
	for _, source := range options[Source] {
		basename, _ := ExtractBasenameAndPermission(source)
		sourceBasenames = append(sourceBasenames, []string{basename}...)
	}

	entrypointArgs := JoinStringSlices(
		sourceBasenames,
		AddPrefix(options[BuildArg], "-b"),
		AddPrefix(options[Arg], "-a"),
	)

	readFromStdin := false

	if stat, _ := os.Stdin.Stat(); (stat.Mode() & os.ModeCharDevice) == 0 {
		readFromStdin = true
	}

	container, err := client.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:        dockerImage,
			Cmd:          entrypointArgs,
			StdinOnce:    true,
			OpenStdin:    true,
			AttachStdin:  true,
			AttachStderr: true,
			AttachStdout: true,
			Tty:          !readFromStdin,
		},
		HostConfig: &docker.HostConfig{
			Binds: BuildVolumeArgs(
				RetrievePath(options[TargetDir]),
				append(options[Source], options[Include]...)),
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

	waiter, err := client.AttachToContainerNonBlocking(docker.AttachToContainerOptions{
		Container:    container.ID,
		InputStream:  os.Stdin,
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
		Stream:       true,
		Stdin:        true,
		Stdout:       true,
		Stderr:       true,
		Logs:         true,
		RawTerminal:  !readFromStdin,
	})

	type ClientRunResult struct {
		Code  int
		Error error
	}

	if timeout, ok := options[Timeout]; ok && len(timeout) == 1 {
		done := make(chan ClientRunResult)
		go func() {
			if err := waiter.Wait(); err != nil {
				log.Fatal(err)
			}

			code, err := client.WaitContainer(container.ID)
			result := ClientRunResult{
				code,
				err,
			}
			done <- result
		}()

		// Start a timer
		num, _ := strconv.Atoi(timeout[0])
		t := time.Duration(num) * time.Second
		timeout := time.After(t)

		select {
		case <-timeout:
			// kill container using container ID
			err := client.KillContainer(docker.KillContainerOptions{ID: container.ID})
			if err != nil {
				log.Fatal(err)
			}
			return timeoutStatusCode
		case result := <-done:
			code := result.Code
			err := result.Error
			if err != nil {
				log.Fatal(err)
			}
			return code
		}
	} else {
		if err := waiter.Wait(); err != nil {
			log.Fatal(err)
			return 1
		}

		code, err := client.WaitContainer(container.ID)
		if err != nil {
			log.Fatal(err)
		}
		return code
	}
}

func validate(cliParser CLI) bool {
	options := cliParser.Options

	hasVersionFlag := len(options[VersionFlag]) == 1
	hasSources := len(options[Source]) > 0
	shouldClean := len(options[CleanFlag]) > 0

	if hasSources || shouldClean {
		return true
	}

	if hasVersionFlag {
		DisplayVersion(cliParser.Filename)
		return false
	}

	DisplayHelp(cliParser.Filename)
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
		return fmt.Errorf("request to Docker host timed out")
	}
}

func main() {
	cliParser := ParseOsArgs(os.Args)

	if validate(cliParser) {
		if err := validateDocker(); err != nil {
			log.Fatal(err)
		} else {
			os.Exit(RunDexecContainer(cliParser))
		}
	}
}
