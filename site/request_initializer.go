package main

import (
    "net/http"
    "fmt"
)

func establishRequestContext(req *http.Request, resp http.ResponseWriter) RequestContext {
    session, err := store.Get(req, config.SessionConfig.SessionName)
    if err != nil {
        panic(err)
    }

    database, err := getDatabaseConnection()
    if err != nil {
        panic(fmt.Sprintf("Database error: %v\n", err))
    }

    ctx := RequestContext {
        request    : req,
        response   : resp,
        method     : req.Method,
        routes     : getPathRoutes(req.URL.String()),
        session    : session,
        redis      : getRedisConnection(),
        database   : database,
    }
    ctx.userId = ctx.getUserId()
    
    return ctx
}
