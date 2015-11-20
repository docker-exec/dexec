# Contributing

Fork the repo. Make your changes. Raise a PR.

All PRs are welcome, and changes don't just have to be additive - if you have ideas about how to improve existing code, please feel free to submit these.

## Code

```dexec``` is written in go and changes must adhere to the default settings for all the following tools:

 * go vet
 * gofmt
 * golint

A good place to start is with the [go-plus](https://github.com/joefitzgerald/go-plus/) package for [Atom](https://atom.io/) which automatically verifies the code against those tools.

## Dependencies

Dependencies are managed using the experimental vendor feature introduced in Go 1.5. Versioned library code is committed to the dexec repo to guarantee reproducible builds, and the tool used to achieve this is [govendor](https://github.com/kardianos/govendor/). Please read the [govendor documentation](https://github.com/kardianos/govendor/blob/master/README.md#vendor-tool-for-go) if you need to add a new library.

## Unit Tests

Unit tests are required for new contributions in most cases. Don't break the existing ones. Run the following command from the path that you checked out dexec to run the unit tests:

```sh
go test ./...
```

## Acceptance Tests

A suite of acceptance tests can be found in the [dexec/_bats](https://github.com/docker-exec/dexec/tree/master/_test) folder. If you are adding a new image, you must also add to the acceptance tests.

Vagrant's VirtualBox provider is used to spin up an Ubuntu VM on which Docker is installed. The current ```dexec``` folder is mounted in the VM, built and then run against four [bats](https://github.com/sstephenson/bats) tests for each docker-exec image. The tests are:

 * Verify simple output.
 * Verify unicode output.
 * Verify execution arguments are available to script.
 * Verify code can be run as standalone script.

Run the following command from the path that you checked out dexec to run the acceptance tests:

```sh
./acceptance-tests.sh run
```

By default, the vagrant box will be restored to its initial state at the start of every run. This guarantees sandboxing of the application and images, but means that all the images will be redownloaded every time you run the tests. This may not be desirable as it takes several minutes to download the images. To avoid this you can add the ```--no-clean``` option to reuse images that have already been downloaded to the VM:

```sh
./acceptance-tests.sh run --no-clean
```
