package main

import (
    "fmt"
    "flag"
    "github.com/alculquicondor/gotter/accountservice/dbclient"
    "github.com/alculquicondor/gotter/accountservice/service"
    "github.com/spf13/viper"
    "github.com/alculquicondor/gotter/common/config"
    "github.com/alculquicondor/gotter/common/messaging"
    "os"
    "os/signal"
    "syscall"
)


var appName = "accountservice"


func init() {
    profile := flag.String("profile", "test", "Environment profile, something similar to spring profiles")
    configServerUrl := flag.String("configServerUrl", "http://configserver:8888", "Address to config server")
    configBranch := flag.String("configBranch", "master", "git branch to fetch configuration from")
    flag.Parse()

    // Pass the flag values into viper.
    viper.Set("profile", *profile)
    viper.Set("configServerUrl", *configServerUrl)
    viper.Set("configBranch", *configBranch)
}


func main() {
    fmt.Printf("Starting %v\n", appName)

    config.LoadConfigurationFromBranch(
        viper.GetString("configServerUrl"),
        appName,
        viper.GetString("profile"),
        viper.GetString("configBranch"))

    initializeBoltClient()
    initializeMessaging()

    handleSigterm(func() {
        service.MessagingClient.Close()
    })

    service.StartWebServer(viper.GetString("server_port"))
}


func initializeBoltClient() {
    service.DbClient = &dbclient.BoltClient{}
    service.DbClient.OpenBoltDb()
    service.DbClient.Seed()
}


func initializeMessaging() {
    if !viper.IsSet("amqp_server_url") {
        panic("No 'amqp_server_url' set in configuration, cannot start")
    }

    service.MessagingClient = &messaging.MessagingClient{}
    service.MessagingClient.ConnectToBroker(viper.GetString("amqp_server_url"))
    service.MessagingClient.Subscribe(viper.GetString("config_event_bus"), "topic", appName,
        config.HandleRefreshEvent)
}


func handleSigterm(handleExit func()) {
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    signal.Notify(c, syscall.SIGTERM)
    go func() {
        <-c
        handleExit()
        os.Exit(1)
    }()
}