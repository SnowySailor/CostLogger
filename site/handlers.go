package main

import "fmt"

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

func postTransaction(ctx RequestContext) {
    ctx.successRaw("Post transaction")
}

func postSettings(ctx RequestContext) {
    ctx.successRaw("Post settings")
}

func postLogin(ctx RequestContext) {
    ctx.successRaw("Post login")
}

func getCreateUser(ctx RequestContext) {
    var pageData PageData
    inputForm := makeHtmlWithTemplate("../templates/user_create.template", pageData)
    pageData = makePageData("Create user", inputForm, make([]Link, 0), make([]Link, 0))
    ctx.successPage(pageData)
}

func postCreateUser(ctx RequestContext) {
    username, ok       := ctx.getFormValue("username")
    email, ok          := ctx.getFormValue("email")
    displayName, ok    := ctx.getFormValue("displayName")
    password, ok       := ctx.getFormValue("password")
    passwordVerify, ok := ctx.getFormValue("passwordVerify")

    // Just to get it to compile for now
    if ok {

    }

    if password != passwordVerify {
        ctx.badRequestRaw(fmt.Sprintf("Passwords do not match: %v, %v", password, passwordVerify))
        return
    }

    user := User {
        Username: username,
        DisplayName: displayName,
        Email: email,
        PasswordHash: "testinghash",
    }

    userId, err := ctx.insertUser(user)
    if err != nil {
        ctx.badRequestRaw(fmt.Sprintf("%v", err))
        return
    }

    ctx.successRaw(fmt.Sprintf("Created user %v", userId))
}