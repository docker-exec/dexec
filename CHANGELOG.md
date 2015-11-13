# Change Log
All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

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
