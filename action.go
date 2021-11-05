package main

import (
    "fmt"
    "github.com/brunoa19/shipa-github-actions/shipa"
    "os"
)

func main() {
    if host, ok := os.LookupEnv("SHIPA_HOST"); ok {
        fmt.Println("host:", host)
    }

    if token, ok := os.LookupEnv("SHIPA_TOKEN"); ok {
        fmt.Println("token:", token)
    }

    _, err := shipa.New()
    if err != nil {
        fmt.Println("ERR failed create shipa client:", err)
    } else {
        fmt.Println("shipa client created")
    }

    args := os.Args[1:]
    fmt.Println("Input args:", args)
    if len(args) > 0 {
        path := args[0]
        fmt.Println("Input file path:", path)

        if _, err := os.Stat(path); err != nil {
            fmt.Println("ERR file stats:", err)
        } else {
            fmt.Println("File exists: OK")

            rawContent, err := os.ReadFile(path)
            if err != nil {
                fmt.Println("ERR read file failed:", err)
            } else {
                fmt.Println(string(rawContent))
            }

        }
    }
}
