package service


import (
    "fmt"
    "strconv"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "github.com/alculquicondor/gotter/accountservice/dbclient"
    "github.com/alculquicondor/gotter/common/netutils"
    "github.com/gorilla/mux"
    "github.com/alculquicondor/gotter/accountservice/model"
)


var DbClient dbclient.IBoltClient
var isHealthy = true
var client = &http.Client{}


func init() {
    var transport http.RoundTripper = &http.Transport{
        DisableKeepAlives: true,
    }
    client.Transport = transport
}


func GetAccount(w http.ResponseWriter, r *http.Request) {
    var accountId = mux.Vars(r)["accountId"]
    account, err := DbClient.QueryAccount(accountId)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }
    account.ServedBy = netutils.GetIP()

    quote, err := getQuote()
    if err == nil {
        account.Quote = quote
    }

    data, _ := json.Marshal(account)
    writeJsonResponse(w, http.StatusOK, data)
}


func getQuote() (model.Quote, error) {
    req, _ := http.NewRequest("GET", "http://quotesservice:8080/api/quote?strength=4", nil)
    resp, err := client.Do(req)
    if err == nil && resp.StatusCode == 200 {
        quote := model.Quote{}
        bytes, _ := ioutil.ReadAll(resp.Body)
        json.Unmarshal(bytes, &quote)
        return quote, nil
    } else {
        return model.Quote{}, fmt.Errorf("some error ocurred with quotesservice")
    }
}


func HealthCheck(w http.ResponseWriter, r *http.Request) {
    dbUp := DbClient.Check()
    if dbUp && isHealthy {
        data, _ := json.Marshal(healthCheckResponse{Status: "UP"})
        writeJsonResponse(w, http.StatusOK, data)
    } else {
        data, _ := json.Marshal(healthCheckResponse{Status: "Database inaccessible"})
        writeJsonResponse(w, http.StatusServiceUnavailable, data)
    }
}


func SetHealthyState(w http.ResponseWriter, r *http.Request) {
    state, err := strconv.ParseBool(mux.Vars(r)["state"])
    if err != nil {
        fmt.Println("Invalid request to SetHealthyState, allowed values are true or false")
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    isHealthy = state
    w.WriteHeader(http.StatusOK)
}


func writeJsonResponse(w http.ResponseWriter, status int, data []byte) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Content-Length", strconv.Itoa(len(data)))
    w.WriteHeader(status)
    w.Write(data)
}


type healthCheckResponse struct {
    Status string `json:"status"`
}
