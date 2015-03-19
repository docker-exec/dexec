# dexec
Executes source code via Docker images at https://github.com/docker-exec.

# Target Features

## Version 1

### Execute single source file
```
dexec foo.cpp
```

### Execute multiple source files
```
dexec foo.cpp bar.cpp
```

### Pass arguments for build
```
dexec foo.cpp --arg -std=c++11
```
or
```
dexec foo.cpp -a -std=c++11
```

### Build only
```
dexec --build foo.cpp
```
or
```
dexec -b foo.cpp
```

### Pass arguments for execution
```
dexec foo.cpp --exec-arg hello --exec-arg world
```
or
```
dexec foo.cpp -e hello -e world
```

## Version 2

### Override the image used to perform build
```
dexec --image=cpp foo.cpp
```
or
```
dexec -i cpp foo.cpp
```

### Execute source code in a directory (requires image override?)
```
dexec --image=cpp .
```