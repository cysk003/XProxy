package main

import (
    log "github.com/sirupsen/logrus"
    "gopkg.in/yaml.v3"
    "net"
    "strings"
)

var v4Bypass []string
var v6Bypass []string
var dnsServer []string

var v4Gateway string
var v4Address string
var v6Gateway string
var v6Address string

type NetConfig struct {
    Gateway string `yaml:"gateway"` // network gateway
    Address string `yaml:"address"` // network address
}

type Config struct {
    Update struct {
        Cron string            `yaml:"cron"`
        Url  map[string]string `yaml:"url"`
    } `yaml:"update"`
    Proxy struct {
        Sniff    bool           `yaml:"sniff"`
        Redirect bool           `yaml:"redirect"`
        Http     map[string]int `yaml:"http"`
        Socks    map[string]int `yaml:"socks"`
        AddOn    []interface{}  `yaml:"addon"`
    } `yaml:"proxy"`
    Network struct {
        DNS    []string  `yaml:"dns"`    // system dns server
        ByPass []string  `yaml:"bypass"` // cidr bypass list
        IPv4   NetConfig `yaml:"ipv4"`   // ipv4 network configure
        IPv6   NetConfig `yaml:"ipv6"`   // ipv6 network configure
    } `yaml:"network"`
}

func isIP(ipAddr string, isCidr bool) bool {
    if !isCidr {
        return net.ParseIP(ipAddr) != nil
    }
    _, _, err := net.ParseCIDR(ipAddr)
    return err == nil
}

func isIPv4(ipAddr string, isCidr bool) bool {
    return isIP(ipAddr, isCidr) && strings.Contains(ipAddr, ".")
}

func isIPv6(ipAddr string, isCidr bool) bool {
    return isIP(ipAddr, isCidr) && strings.Contains(ipAddr, ":")
}

func loadConfig(rawConfig []byte) {
    config := Config{}
    log.Debugf("Decode yaml content -> \n%s", string(rawConfig))
    err := yaml.Unmarshal(rawConfig, &config) // yaml (or json) decode
    if err != nil {
        log.Panicf("Decode config file error -> %v", err)
    }
    log.Debugf("Decoded config -> %v", config)

    for _, address := range config.Network.DNS { // dns options
        if isIPv4(address, false) || isIPv6(address, false) {
            dnsServer = append(dnsServer, address)
        } else {
            log.Panicf("Invalid DNS server -> %s", address)
        }
    }
    log.Infof("DNS server -> %v", dnsServer)

    for _, address := range config.Network.ByPass { // bypass options
        if isIPv4(address, true) {
            v4Bypass = append(v4Bypass, address)
        } else if isIPv6(address, true) {
            v6Bypass = append(v6Bypass, address)
        } else {
            log.Panicf("Invalid bypass CIDR -> %s", address)
        }
    }
    log.Infof("IPv4 bypass CIDR -> %s", v4Bypass)
    log.Infof("IPv6 bypass CIDR -> %s", v6Bypass)

    v4Address = config.Network.IPv4.Address
    v4Gateway = config.Network.IPv4.Gateway
    if v4Address != "" && !isIPv4(v4Address, true) {
        log.Panicf("Invalid IPv4 address -> %s", v4Address)
    }
    if v4Gateway != "" && !isIPv4(v4Gateway, false) {
        log.Panicf("Invalid IPv4 gateway -> %s", v4Gateway)
    }
    log.Infof("IPv4 -> address = %s | gateway = %s", v4Address, v4Gateway)

    v6Address = config.Network.IPv6.Address
    v6Gateway = config.Network.IPv6.Gateway
    if v6Address != "" && !isIPv6(v6Address, true) {
        log.Panicf("Invalid IPv6 address -> %s", v6Address)
    }
    if v6Gateway != "" && !isIPv6(v6Gateway, false) {
        log.Panicf("Invalid IPv6 gateway -> %s", v6Gateway)
    }
    log.Infof("IPv6 -> address = %s | gateway = %s", v6Address, v6Gateway)

    enableSniff = config.Proxy.Sniff
    log.Infof("Connection sniff -> %v", enableSniff)
    enableRedirect = config.Proxy.Redirect
    log.Infof("Connection redirect -> %v", enableRedirect)
    httpInbounds = config.Proxy.Http
    log.Infof("Http inbounds -> %v", httpInbounds)
    socksInbounds = config.Proxy.Socks
    log.Infof("Socks5 inbounds -> %v", socksInbounds)
    addOnInbounds = config.Proxy.AddOn
    log.Infof("Add-on inbounds -> %v", addOnInbounds)

    updateCron = config.Update.Cron
    log.Infof("Update cron -> %s", updateCron)
    updateUrls = config.Update.Url
    log.Infof("Update url -> %v", updateUrls)
}
