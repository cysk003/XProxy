package main

import (
    "fmt"
)

func main() {
    fmt.Println("XProxy start")

    yamlContent := []byte(`
network:
  dns:  # system's dns server
    - 223.5.5.5
    - 119.29.29.29
#    - fesdc.fardf.afa
  ipv4:  # ipv4 network configure
    gateway: 192.168.2.1
    address: 192.168.2.2/24
    forward: false
  ipv6: null
  bypass:
    - 169.254.0.0/16
    - fc00::/7
    - 224.0.0.0/3
`)

    loadConfig(yamlContent)
}
