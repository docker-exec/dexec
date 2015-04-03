# dexec [![Build Status](https://travis-ci.org/docker-exec/dexec.svg?branch=master)](https://travis-ci.org/docker-exec/dexec)

Executes source code via Docker images at https://github.com/docker-exec.

Currently in development...

### Execute source files

```sh
$ dexec foo.cpp
$ dexec foo.cpp bar.cpp
```

### Pass individual arguments for build

```sh
$ dexec foo.cpp --build-arg=-std=c++11
$ dexec foo.cpp --build-arg -std=c++11
$ dexec foo.cpp -b -std=c++11
```

### Pass arguments for execution

```sh
$ dexec foo.cpp --arg=hello --arg=world --arg='hello world'
$ dexec foo.cpp --arg hello --arg world --arg 'hello world'
$ dexec foo.cpp -a hello -a world -a 'hello world'
```

### Specify location of source files

```sh
$ dexec -C /path/to/sources foo.cpp bar.cpp
```

### Include files and folders mounted in Docker container without passing to compiler

```sh
$ dexec foo.cpp --include=bar.hpp
$ dexec foo.cpp --include bar.hpp
$ dexec foo.cpp -i bar.hpp
```

### Support shebang in source files

```c++
#!/usr/bin/dexec
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

