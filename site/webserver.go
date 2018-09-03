// User can insert new transaction
    // Contains amount and other users from a dropdown
// User can edit their own transaction (both people and amounts)
// At the end of each week, the previous week's transactions (7 days old to 14 days old) are sent to all users that owe someone else money
// 

package main

import (
    "log"
    "net/http"
    "fmt"
)

var config AppConfig

func main() {
    config.getAppConfig()
    http.HandleFunc("/", routeRequest)
    fmt.Println("Listening on " + "localhost" + ":8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
