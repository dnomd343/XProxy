package main

import (
    "fmt"
    "os"
)

func main() {
    fmt.Println("XProxy start")
    content, err := os.ReadFile("test.yml")
    if err != nil {
        panic(err)
    }
    loadConfig(content)
    fmt.Println("DNS ->", dnsServer)
    fmt.Println("v4Bypass ->", v4Bypass)
    fmt.Println("v6Bypass ->", v6Bypass)
}
