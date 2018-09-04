package main

import (
    "net/http"
    "strings"
    "fmt"
)

func routeRequest(resp http.ResponseWriter, req *http.Request) {
    ctx := establishRequestContext(req, resp)

    db, err := getDatabaseConnection()
    if err != nil {
        fmt.Fprintf(resp, fmt.Sprintf("Database error: %v\n", err))
        return
    }
    ctx.database = db

    mainRoute := firstOrDefault(ctx.routes)
    method   := ctx.method
    printList(ctx.routes)
    printStrLStrMap(getQueryParams(*req))

    if method == "GET" {
        if mainRoute == "" || mainRoute == "home" {
            getHome(ctx)
        } else if mainRoute == "settings" {
            getSettings(ctx)
        } else if mainRoute == "transaction" {
            getTransaction(ctx)
        } else if mainRoute == "feed" {
            getFeed(ctx)
        } else {
            ctx.notFoundPage("Invalid route")
        }
    } else if method == "POST" {
        if mainRoute == "transaction" {
            postTransaction(ctx)
        } else if mainRoute == "settings" {
            postSettings(ctx)
        } else if mainRoute == "login" {
            postLogin(ctx)
        } else {
            ctx.notFoundPage("Invalid route")
        }
    } else {
        ctx.badRequestRaw("HTTP Method not supported")
    }
}

func establishRequestContext(req *http.Request, resp http.ResponseWriter) RequestContext {
    session, err := store.Get(req, config.SessionConfig.SessionName)
    if err != nil {
        panic(err)
    }

    ctx := RequestContext {
        request  : req,
        response : resp,
        method   : req.Method,
        userId   : 0,
        routes   : getPathRoutes(req.URL.String()),
        session  : session,
    }
    return ctx
}

func splitPathRoutes(path string) []string {
    return denullStrList(strings.Split(path, "/"))
}

func getPathRoutes(path string) []string {
    routes := splitPathRoutes(path)
    if len(routes) == 0 {
        return routes
    }
    // Split the query from the last route
    lastRoute := getLastPathRoute(path)
    // "replace" the last route in `routes` with the new lastRoute
    return append(routes[:len(routes)-1], lastRoute)
}

func getLastPathRoute(path string) string {
    // Get all routes
    routes := splitPathRoutes(path)
    if len(routes) == 0 {
        // If there were no routes returned, return empty string
        return ""
    }
    // Get the last route and the location of the beginning of the url query
    lastRoute := routes[len(routes)-1]
    queryIndex   := strings.Index(lastRoute, "?")
    if queryIndex == -1 {
        // If there is no query, we don't have to go further. Just return the last route.
        return lastRoute
    }
    // Return the last route without the query
    return lastRoute[:queryIndex]
}