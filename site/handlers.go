package main

import (
    "fmt"
    "net/http"
)

func serveFile(ctx RequestContext) {
    location := removeLeadingSlash(ctx.request.URL.Path)
    http.ServeFile(ctx.response, ctx.request, location)
}

func getHome(ctx RequestContext) {
    ctx.successRaw("Get home")
}

func getSettings(ctx RequestContext) {
    ctx.successRaw("Get settings")
}

func getTransaction(ctx RequestContext) {
    ctx.successRaw("Get transaction")
}

func getFeed(ctx RequestContext) {
    ctx.successRaw("Get feed")
}

func getRegisterUser(ctx RequestContext) {
    var pageData PageData
    inputForm := makeHtmlWithTemplate("../templates/user_create.template", pageData)
    pageData   = makePageData("Register", inputForm, []Link{{Url:"/static/styles/global.css"}}, []Link{{Url:"/static/scripts/global.js"}})
    ctx.successPage(pageData)
}

func postTransaction(ctx RequestContext) {
    ctx.successRaw("Post transaction")
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
        ctx.successJSON(response)
    }
}

func postRegisterUser(ctx RequestContext) {
    user, err := ctx.extractNewUser()

    if err != nil {
        ctx.badRequestRaw(err.Error())
        return
    }

    userId, err := ctx.insertUser(user)
    if err != nil {
        ctx.badRequestRaw(err.Error())
        return
    }

    ctx.successRaw(fmt.Sprintf("Created user %v", userId))
}