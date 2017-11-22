package service

import "net/http"


type Route struct {
    Name        string
    Method 	    string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

var routes = []Route {
    {
        "GetAccount",
        "GET",
        "/accounts/{accountId}",
        GetAccount,
    },
    {
        "HealthCheck",
        "GET",
        "/health",
        HealthCheck,
    },
    {
        "Testability",
        "GET",
        "/testability/healthy/{state}",
        SetHealthyState,
    },
}
