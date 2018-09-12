package main

import (
    "net/http"
)

func serveFile(ctx RequestContext) {
    location := removeLeadingSlash(ctx.request.URL.Path)
    http.ServeFile(ctx.response, ctx.request, location)
}

func getHome(ctx RequestContext) {
    if ctx.isUserLoggedIn() {
        home     := ctx.makeHtmlWithHeader("../templates/home.template", PageData{})
        pageData := makePageData("Home", home, []Link{{Url:"/static/styles/global.css"}}, []Link{{Url:"/static/scripts/global.js"}})
        ctx.successPage(pageData)
    } else {
        ctx.redirect("login")
    }
}

func getSettings(ctx RequestContext) {
    if ctx.isUserLoggedIn() {
        ctx.successRaw("Get settings")
    } else {
        ctx.redirect("login")
    }
}

func getTransaction(ctx RequestContext) {
    ctx.successRaw("Get transaction")
}

func getFeed(ctx RequestContext) {
    if ctx.isUserLoggedIn() {
        ctx.successRaw("Get feed")
    } else {
        ctx.redirect("login")
    }
}

func getRegisterUser(ctx RequestContext) {
    inputForm := ctx.makeHtmlWithHeader("../templates/register.template", PageData{})
    pageData  := makePageData("Register", inputForm, []Link{{Url:"/static/styles/global.css"}}, []Link{{Url:"/static/scripts/global.js"}})
    ctx.successPage(pageData)
}

func getLogin(ctx RequestContext) {
    if !ctx.isUserLoggedIn() {
        inputForm := ctx.makeHtmlWithHeader("../templates/login.template", PageData{})
        pageData  := makePageData("Login", inputForm, []Link{{Url:"/static/styles/global.css"}}, []Link{{Url:"/static/scripts/global.js"}})
        ctx.successPage(pageData)
    } else {
        ctx.redirect("home")
    }
}

func getLogout(ctx RequestContext) {
    if ctx.isUserLoggedIn() {
        ctx.logoutUser()
    }
    ctx.redirect("login")
}

func postTransaction(ctx RequestContext) {
    response := makeJSONResponse("")

    // Get transaction
    transaction, err := ctx.extractTransaction()
    if err != nil {
        response.Message = err.Error()
        ctx.badRequestJSON(response)
        return
    }

    // Validate transaction
    err = ctx.validateTransaction(transaction)
    if err != nil {
        response.Message = err.Error()
        ctx.badRequestJSON(response)
        return
    }

    // Insert transaction
    _, err = ctx.insertTransaction(transaction)
    if err != nil {
        response.Message = err.Error()
        ctx.badRequestJSON(response)
        return
    }
    ctx.successJSON(response)
}

func postSettings(ctx RequestContext) {
    ctx.successRaw("Post settings")
}

func postLogin(ctx RequestContext) {
    msg, succ := ctx.attemptUserLogin()
    response  := makeJSONResponse(msg)
    if !succ {
        ctx.badRequestJSON(response)
    } else {
        response.RedirectUrl = "home"
        ctx.successJSON(response)
    }
}

func postLogout(ctx RequestContext) {
    if ctx.isUserLoggedIn() {
        ctx.logoutUser()
    }
    ctx.redirect("login")
}

func postRegisterUser(ctx RequestContext) {
    response  := makeJSONResponse("")

    // Get user
    user, err := ctx.extractNewUser()
    if err != nil {
        response.Message = err.Error()
        ctx.badRequestJSON(response)
        return
    }

    // Validate user
    err = ctx.validateNewUser(user)
    if err != nil {
        response.Message = err.Error()
        ctx.badRequestJSON(response)
        return
    }

    // Create user
    _, err = ctx.insertUser(user)
    if err != nil {
        response.Message = err.Error()
        ctx.badRequestJSON(response)
        return
    }

    // Redirect to login
    response.RedirectUrl = "login"
    ctx.successJSON(response)
}