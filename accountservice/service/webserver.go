package service

import (
    "net/http"
    "github.com/sirupsen/logrus"
)


func StartWebServer(port string) {
    r := NewRouter()
    http.Handle("/", r)

    logrus.Info("Starting HTTP service at " + port)
    err := http.ListenAndServe(":" + port, nil)  // Go routine will block here

    if err != nil {
        logrus.Error("An error occurred starting HTTP listener at port " + port)
        logrus.Error("Error: " + err.Error())
    }
}
