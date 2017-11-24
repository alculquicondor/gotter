package main

import (
    "fmt"
    "github.com/alculquicondor/gotter/accountservice/dbclient"
    "github.com/alculquicondor/gotter/accountservice/service"
)


var appName = "accountService"


func main() {
    fmt.Printf("Starting %v\n", appName)
    initializeBoltClient()
    service.StartWebServer("6767")
}


func initializeBoltClient() {
    service.DbClient = &dbclient.BoltClient{}
    service.DbClient.OpenBoltDb()
    service.DbClient.Seed()
}
