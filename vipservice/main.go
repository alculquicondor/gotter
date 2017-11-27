package main

import (
    "flag"
    "fmt"
    "github.com/alculquicondor/gotter/vipservice/service"
    "github.com/alculquicondor/gotter/common/messaging"
    "github.com/alculquicondor/gotter/common/config"
    "github.com/spf13/viper"
    "github.com/streadway/amqp"
    "os"
    "os/signal"
    "syscall"
    "github.com/sirupsen/logrus"
)


var appName = "vipservice"
var messagingClient messaging.IMessagingClient


func init() {
    configServerUrl := flag.String("configServerUrl", "http://configserver:8888", "Address to config server")
    profile := flag.String("profile", "test", "Environment profile, something similar to spring profiles")
    configBranch := flag.String("configBranch", "master", "git branch to fetch configuration from")
    flag.Parse()

    viper.Set("profile", *profile)
    viper.Set("configServerUrl", *configServerUrl)
    viper.Set("configBranch", *configBranch)
}


func main() {
    logrus.Info("Starting " + appName + "...")

    config.LoadConfigurationFromBranch(viper.GetString("configServerUrl"), appName,
        viper.GetString("profile"), viper.GetString("configBranch"))

    initializeMessaging()
    handleSigterm(func () {
        if messagingClient != nil {
            messagingClient.Close()
        }
    })
    service.StartWebServer(viper.GetString("server_port"))
}


func onMessage(delivery amqp.Delivery) {
    logrus.Infof("Got a message: %v", string(delivery.Body))
}


func initializeMessaging() {
    if !viper.IsSet("amqp_server_url") {
        panic("No 'broker_url' set in configuration, cannot start")
    }
    messagingClient = &messaging.MessagingClient{}
    messagingClient.ConnectToBroker(viper.GetString("amqp_server_url"))

    // Call the subscribe method with queue name and callback function
    err := messagingClient.SubscribeToQueue("vip_queue", appName, onMessage)
    failOnError(err, "Could not start subscribe to vip_queue")

    err = messagingClient.Subscribe(viper.GetString("config_event_bus"), "topic", appName, config.HandleRefreshEvent)
    failOnError(err, "Could not start subscribe to " + viper.GetString("config_event_bus") + " topic")
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


func failOnError(err error, msg string) {
    if err != nil {
        logrus.Fatalf("%s: %s", msg, err)
        panic(fmt.Sprintf("%s: %s", msg, err))
    }
}