# dexec [![Build Status](https://travis-ci.org/docker-exec/dexec.svg?branch=master)](https://travis-ci.org/docker-exec/dexec)  [ ![Download](https://api.bintray.com/packages/dexec/release/dexec/images/download.svg) ](https://bintray.com/dexec/release/dexec/_latestVersion)

A command line utility for executing code in many different languages using the Docker Exec images, written in Go.

![dexec demo animation](https://docker-exec.github.io/images/dexec-short-1.0.1.gif)

## Installation

### Using Bintray

Download the appropriate binary for your OS and architecture, then unzip or untar and move the ```dexec``` executable to where it can be found on your PATH.

| OS      | 64-bit | 32-bit |
| ------- | ------ | ------ |
| Linux   |  [64-bit](https://bintray.com/artifact/download/dexec/release/dexec_1.0.3_linux_amd64.tar.gz) | [32-bit](https://bintray.com/artifact/download/dexec/release/dexec_1.0.3_linux_386.tar.gz) |
| Mac     |  [64-bit](https://bintray.com/artifact/download/dexec/release/dexec_1.0.3_darwin_amd64.zip) | [32-bit](https://bintray.com/artifact/download/dexec/release/dexec_1.0.3_darwin_386.zip) |
| Windows |  [64-bit](https://bintray.com/artifact/download/dexec/release/dexec_1.0.3_windows_amd64.zip) | [32-bit](https://bintray.com/artifact/download/dexec/release/dexec_1.0.3_windows_386.zip) |

Binaries for other distributions are available on [Bintray](https://bintray.com/dexec/release/dexec/1.0.3).

### Using Go

If you have Go installed and configured correctly, you can install ```dexec``` with the ```go get``` command.

```sh
$ go get github.com/docker-exec/dexec
```

### Using Homebrew

If you're on OSX you can install the latest release of ```dexec``` with brew.

```sh
$ brew install docker-exec/formula/dexec
```

## Reference

These examples use a .cpp source file, but any of the supported languages can be used instead. Arguments can be passed in any order, using any style of the acceptable switch styles described.

The application provides help and version information as follows:

```sh
$ dexec --version
$ dexec --help
```

### Pass source files to execute

Multiple source files can be passed to the compiler or interpreter as follows. The first source file's extension is used to pick the appropriate Docker Exec image, e.g. .cpp retrieves dexec/cpp from the Docker registry.

```sh
$ dexec foo.cpp
$ dexec foo.cpp bar.cpp
```

The sources are mounted individually using the default Docker mount permissions (rw) and can be specified by appending :ro or :rw to the source file.

### Pass individual arguments for build

For compiled languages, arguments can be passed to the compiler.

```sh
$ dexec foo.cpp --build-arg=-std=c++11
$ dexec foo.cpp --build-arg -std=c++11
$ dexec foo.cpp -b -std=c++11
```

### Pass arguments for execution

Arguments can be passed to the executing code. Enclose arguments with single quotes to preserve whitespace.

```sh
$ dexec foo.cpp --arg=hello --arg=world --arg='hello world'
$ dexec foo.cpp --arg hello --arg world --arg 'hello world'
$ dexec foo.cpp -a hello -a world -a 'hello world'
```

### Specify location of source files

By default, ```dexec``` assumes the sources are in the directory from which it is being invoked from. It is possible to override the working directory by passing the ```-C``` flag.

```sh
$ dexec -C /path/to/sources foo.cpp bar.cpp
```

### Include files and folders mounted in Docker container

Individual files can be mounted without being passed to the compiler, for example header files in C & C++, or input files for program execution. These can be included in the following way.

```sh
$ dexec foo.cpp --include=bar.hpp
$ dexec foo.cpp --include bar.hpp
$ dexec foo.cpp -i bar.hpp
```

In addition, a program may require read and/or write access to several files on the host system. The most efficient way to achieve this is to include a directory.

```sh
$ dexec foo.cpp --include=.
$ dexec foo.cpp --include .
$ dexec foo.cpp -i .
```

Files and directories are relative to either the current working directory, or the directory specified with the ```-C``` flag.

As with sources, included files and directories are mounted using the default Docker mount permissions (rw) and can be specified by appending :ro or :rw to the source file.

### Override the image used by dexec

```dexec``` stores a map of file extensions to Docker images and uses this to look up the right image to run for a given source file. This can be overridden in the following ways:

#### Override image by name/tag

```sh
$ dexec foo.c --image=dexec/lang-cpp
$ dexec foo.c --image dexec/lang-cpp
$ dexec foo.c -m dexec/lang-cpp
```

This will cause ```dexec``` to attempt to use the supplied image. If no image version is specified, "latest" is used.

#### Override image by file extension

```sh
$ dexec foo.c --extension=cpp
$ dexec foo.c --extension cpp
$ dexec foo.c -e cpp
```

This will cause ```dexec``` to attempt to lookup the image for the supplied extension in its map.

### Reading from STDIN

If no source files are specified, ```dexec``` will read from STDIN. Both manual entry and piping are supported. In both cases ```dexec``` writes the contents of STDIN to a temporary file, executes it and then removes it.

```dexec``` requires either the [extension](#override-image-by-file-extension) or [image](#override-image-by-nametag) to be supplied.

#### Manual entry

Text may be entered until EOF is received (Ctrl-D).

```sh
$ dexec --extension js
Enter your code. Ctrl-D to exit
console.log('hello world')
<Ctrl-D>

hello world
```

#### Pipe

```sh
$ echo "console.log('hello world')" | dexec --extension js
```

```sh
$ curl "http://some-url.com/foo.cpp" | dexec --image dexec/lang-cpp
```

### Force dexec to pull latest version of image

Primarily for debugging purposes, the --update command triggers a ```docker pull``` of the target image before executing the code.

```sh
$ dexec foo.cpp -u
$ dexec foo.cpp --update
```

### Executable source with shebang

```dexec``` can be used to make source code executable by adding a shebang that invokes it at the top of a source file.

The shebang is stripped out at execution time but the original source containing the shebang is preserved.

```c++
#!/usr/bin/env dexec
#include <iostream>
int main() {
    std::cout << "hello world" << std::endl;
}
```

then

```sh
$ chmod +x foo.cpp
$ ./foo.cpp
```

## Contributors

#### [docker-exec/dexec](https://github.com/docker-exec/dexec/graphs/contributors)

 * [Alix Axel](https://github.com/alixaxel)
 * [kroton](https://github.com/kroton)
 * [John Albietz](https://github.com/inthecloud247)

#### [docker-exec/perl](https://github.com/docker-exec/perl/graphs/contributors)

 * [Øyvind Skaar](https://github.com/oyvindsk)

## See also

* [Docker Exec GitHub Page](https://docker-exec.github.io/)
* [Docker Exec GitHub Repositories](https://github.com/docker-exec)
* [Docker Exec Images on Docker Hub](https://hub.docker.com/repos/dexec/)
* [dexec on Bintray](https://bintray.com/dexec/release/dexec/view)
