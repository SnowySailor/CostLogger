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
    "github.com/gorilla/context"
    "github.com/gorilla/sessions"
    "strconv"
)

func main() {
    config.populateAppConfig()
    store = sessions.NewCookieStore([]byte(config.SessionConfig.SecretKey))
    store.Options = &sessions.Options{
        MaxAge:   config.SessionConfig.MaxAge,
        HttpOnly: config.SessionConfig.internHttpOnly,
    }
    portStr := strconv.Itoa(config.WebConfig.Port)
    http.HandleFunc("/", routeRequest)
    fmt.Println("Listening on " + "localhost:" + portStr)
    log.Fatal(http.ListenAndServe(":" + portStr, context.ClearHandler(http.DefaultServeMux)))
}
