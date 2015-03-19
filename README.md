# dexec
Executes source code via Docker images at https://github.com/docker-exec.

# Target Features

## Version 1

### Execute single source file

```sh
dexec foo.cpp
```

### Execute multiple source files

```sh
dexec foo.cpp bar.cpp
```

### Pass arguments for build

```sh
dexec foo.cpp --arg -std=c++11 --arg --oO
```

or

```sh

dexec foo.cpp -a -std=c++11
```

### Build only

```sh
dexec --build foo.cpp
```

or

```sh
dexec -b foo.cpp
```

### Pass arguments for execution

```sh
dexec foo.cpp --exec-arg=hello --exec-arg=world
```

or

```sh
dexec foo.cpp -e hello -e world
```

## Version 2

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

### Execute source code in a directory (requires image override?)

```sh
dexec --image=cpp .
```