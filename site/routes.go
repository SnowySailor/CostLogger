package main

import (
    "net/http"
)

func routeRequest(resp http.ResponseWriter, req *http.Request) {
    ctx := establishRequestContext(req, resp)

    mainRoute := firstOrDefault(ctx.routes)
    printList(ctx.routes)
    printStrLStrMap(getQueryParams(*req))

    if ctx.method == "GET" {
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
    } else if ctx.method == "POST" {
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
