package main

import (
    "fmt"
    "os"
)

func main() {
    if host, ok := os.LookupEnv("SHIPA_HOST"); ok {
        fmt.Println("host:", host)
    }

    if token, ok := os.LookupEnv("SHIPA_TOKEN"); ok {
        fmt.Println("token:", token)
    }

    args := os.Args[1:]
    fmt.Println("Input args:", args)
    if len(args) > 0 {
        path := args[0]
        fmt.Println("Input file path:", path)

        if _, err := os.Stat(path); err != nil {
            fmt.Println("ERR:", err)
        } else {
            fmt.Println("File exists: OK")
        }
    }
}
