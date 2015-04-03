# dexec [![Build Status](https://travis-ci.org/docker-exec/dexec.svg?branch=master)](https://travis-ci.org/docker-exec/dexec)

Executes source code via Docker images at https://github.com/docker-exec.

Currently in development...

### Execute source files

```sh
dexec foo.cpp
dexec foo.cpp bar.cpp
```

### Pass individual arguments for build

```sh
dexec foo.cpp --build-arg=-std=c++11
dexec foo.cpp --build-arg -std=c++11
dexec foo.cpp -a -std=c++11
```

### Pass argument string for build

```sh
dexec foo.cpp --build-argline='-std=c++11 -o bar'
dexec foo.cpp --build-argline '-std=c++11 -o bar'
dexec foo.cpp -B '-std=c++11 -o bar'
```

### Pass arguments for execution

```sh
dexec foo.cpp --exec-arg=hello --exec-arg=world
dexec foo.cpp --exec-arg hello --exec-arg world
dexec foo.cpp -e hello -e world
```

### Pass argument string for execution

```sh
dexec foo.cpp --argline='hello world'
dexec foo.cpp --argline 'hello world'
dexec foo.cpp --A 'hello world'
```

### Specify location of source files

```sh
dexec -C /path/to/sources foo.cpp
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
chmod +x foo.cpp
./foo.cpp
```

### Override the image used to perform build

```sh
dexec --image=cpp foo.cpp
```

or

```sh
dexec -i cpp foo.cpp
```
