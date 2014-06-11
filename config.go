package controller

import (
    "code.google.com/p/gcfg"
    "os"
    "log"
)

type Config struct {
    Connection struct {
        Ip string
        Port string
        User string
        Pass string
    }
    Bindings struct {
        Exchange string
        Routes []string
        QueueName string
    }
}

var Cfg Config

func LoadConfig() {
    cfgFile := os.Getenv("CONSUMER_CONFIG_FILE")
    log.Printf("Loading config using file [%s]", cfgFile)
    err := gcfg.ReadFileInto(&Cfg, cfgFile)
    haltOnError(err, "config load error")
    log.Printf("config loaded")
}
