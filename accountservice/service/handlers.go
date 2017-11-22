package service


import (
    "fmt"
    "strconv"
    "encoding/json"
    "net/http"
    "github.com/alculquicondor/gotter/accountservice/dbclient"
    "github.com/gorilla/mux"
)


var DBClient dbclient.IBoltClient
var IsHealthy = true


func GetAccount(w http.ResponseWriter, r *http.Request) {
    var accountId = mux.Vars(r)["accountId"]
    account, err := DBClient.QueryAccount(accountId)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    data, _ := json.Marshal(account)
    WriteJsonResponse(w, http.StatusOK, data)
}


func HealthCheck(w http.ResponseWriter, r *http.Request) {
    dbUp := DBClient.Check()
    if dbUp && IsHealthy {
        data, _ := json.Marshal(HealthCheckResponse{Status: "UP"})
        WriteJsonResponse(w, http.StatusOK, data)
    } else {
        data, _ := json.Marshal(HealthCheckResponse{Status: "Database inaccessible"})
        WriteJsonResponse(w, http.StatusServiceUnavailable, data)
    }
}


func SetHealthyState(w http.ResponseWriter, r *http.Request) {
    state, err := strconv.ParseBool(mux.Vars(r)["state"])
    if err != nil {
        fmt.Println("Invalid request to SetHealthyState, allowed values are true or false")
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    IsHealthy = state
    w.WriteHeader(http.StatusOK)
}


func WriteJsonResponse(w http.ResponseWriter, status int, data []byte) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Content-Length", strconv.Itoa(len(data)))
    w.WriteHeader(status)
    w.Write(data)
}


type HealthCheckResponse struct {
    Status string `json:"status"`
}
