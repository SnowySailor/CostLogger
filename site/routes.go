package main

import (
    "net/http"
)

func handleRequest(resp http.ResponseWriter, req *http.Request) {
    ctx := establishRequestContext(req ,resp)
    routeRequest(resp, req, ctx)
    err := ctx.destroy()
    if err != nil {
        print(err.Error())
    }
}

func routeRequest(resp http.ResponseWriter, req *http.Request, ctx RequestContext) {
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
        } else if mainRoute == "register" {
            getRegisterUser(ctx)
        } else if mainRoute == "login" {
            getLogin(ctx)
        } else if mainRoute == "logout" {
            postLogout(ctx)
        } else if mainRoute == "users" {
            getUsers(ctx)
        } else if mainRoute == "static" {
            // Try to serve a static file
            serveFile(ctx)
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
        } else if mainRoute == "logout" {
            postLogout(ctx)
        } else if mainRoute == "register" {
            postRegisterUser(ctx)
        } else {
            ctx.notFoundPage("Invalid route")
        }
    } else if ctx.method == "DELETE" {
        if mainRoute == "transaction" {
            deleteTransaction(ctx)
        }
    } else {
        ctx.badRequestRaw("HTTP Method not supported")
    }
}
