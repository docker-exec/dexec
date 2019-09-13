# Change Log
All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased][unreleased]
### Added
- Destroy option in acceptance test script.
- Github Actions support.
- Dropped Travis support.
- Renamed AcceptanceTests dir to .acceptance_tests

### Fixed
- Fixed Stdin example in Readme.md.

### Changed
- Migrate to Go Modules for dependency management.
- Move sources to root.
- Renamed Image struct to ContainerImage to avoid clash with Image enum value.
- Removed .vscode dir.
- Renamed .test dir to AcceptanceTests
- Updated years in LICENSE.

### Removed
- Removed redundant stdin reading code.

## [1.0.7] - 2016-09-18
### Fixed
- Compatibility with Docker remote API 1.24.

## [1.0.6] - 2016-04-24
### Fixed
- Shebang support for Java files now works.
- Ruby and Objective C no longer output 'stdin: not a tty' before program output.
- Regex for extracting source filenames no longer ignores single character filenames e.g. 'a.cpp'.

### Changed
- Re-enabled Java, Ruby and Objective C in acceptance tests.

## [1.0.5] - 2016-04-23
### Fixed
- Receive input from STDIN correctly and display STDOUT/STDERR messages only once.

## [1.0.4] - 2016-04-23
### Added
- Added ability to read from STDIN either manually or via pipe.
- Added ability to specify image to use by file extension.
- Added clean functionality to remove all locally downloaded dexec images.
- Add Vagrant-based acceptance tests.
- Added installation instructions for Homebrew on OSX.
- Add contributing instructions to readme.
- Add vscode tasks configuration.

### Changed
- Migrated from custom code that called the Docker CLI to the library 'fsouza/go-dockerclient' which uses the Docker Remote API.
- Deleted custom code that called the Docker CLI.
- Switch to go vendoring for dependency management using Godeps.

## [1.0.3] - 2015-11-13
### Fixed
- Forward application output to stdout/stderr correctly.

## [1.0.2] - 2015-05-12
### Changed
- Moved Docker, CLI & miscellaneous functionality to separate packages.
- Bumped patch versions of each image by 1 to enable unicode support introduced by those versions.
- Changed the format of the image names from ```dexec/{{language abbreviation}}``` to ```dexec/lang-{{language abbreviation}}```.
- Add contributors section in the readme.

## [1.0.1] - 2015-04-20
### Added
- New languages: R, Nim and Lua.
- Support for Perl 6 via .p6 extension.
- Ability to specify the image used with --specify-image or -s.

### Changed
- Perl 5 is now the default for .pl extension.
- Version command is now suffixed with newline.
- Fixed typo in RunDexecContainer comments.

### Fixed
- Bug in IsDockerPresent and IsDockerPresent where defer method was not correctly called on panic.
- Corrected how paths are handled in Windows allowing volumes to be mounted.

## 1.0.0 - 2015-04-06
### Added
- Command line interface 'dexec'.
- Ability to pass source files to container.
- Container image selected based on source file extension.
- Support for Bash, C, Clojure, CoffeeScript, C++, C#, D, Erlang, F#, Go, Groovy, Haskell, Java, Lisp, Node JS, Objective C, OCaml, Perl, PHP, Python, Racket, Ruby, Rust & Scala.
- Ability to pass arguments to the language's compiler if it has one with --build-arg or -b.
- Ability to pass arguments to the executing code with --arg or -a.
- Ability to pass other files or folders to be mounted in the container with --include or -i.
- Ability to augment source files with a shebang resulting in dexec being called.
- Help dialog.
- Version dialog.

[unreleased]: https://github.com/docker-exec/dexec/compare/v1.0.3...HEAD
[1.0.3]: https://github.com/docker-exec/dexec/compare/v1.0.2...v1.0.3
[1.0.2]: https://github.com/docker-exec/dexec/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/docker-exec/dexec/compare/v1.0.0...v1.0.1
