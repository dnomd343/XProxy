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

    fmt.Println("v4Gateway ->", v4Gateway)
    fmt.Println("v4Address ->", v4Address)
    fmt.Println("v6Gateway ->", v6Gateway)
    fmt.Println("v6Address ->", v6Address)

    fmt.Println("v4Forward ->", v4Forward)
    fmt.Println("v6Forward ->", v6Forward)
}
