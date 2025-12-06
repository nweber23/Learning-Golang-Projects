package main

import (
    "fmt"
    "os"

    "github.com/nweber/github-activity/cmd"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: github-activity <username>")
        os.Exit(1)
    }
    username := os.Args[1]
    if err := cmd.Run(username); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}