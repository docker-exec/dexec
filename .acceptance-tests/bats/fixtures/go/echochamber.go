package main

import "os"
import "fmt"

func main() {
    for _, arg := range os.Args[1:] {
        fmt.Printf("%s\n", arg)
    }
}
