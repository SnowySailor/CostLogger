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
        feed, err := ctx.getFeedHtml()
        if err != nil {
            ctx.badRequestRaw("Error 1 - Internal error rendering page")
            println(err.Error())
            return
        }
        pageData  := makePageData("Feed", feed, []Link{}, []Link{})
        home, err := ctx.makeHtmlWithHeader("../templates/home.template", pageData)
        if err != nil {
            ctx.badRequestRaw("Error 2 - Internal error rendering page")
            return
        }
        pageData = makePageData("Home", home, []Link{{Url:"/static/styles/global.css"}}, []Link{{Url:"/static/scripts/global.js"}})
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

func (ctx *RequestContext) getFeedHtml() (string, error) {
    // Get all transactions
    transactions, err := ctx.getUserTransactions(ctx.userId)
    if err != nil {
        println(err.Error())
        return "", err
    }

    // Get all users
    users, err := ctx.getAllUsers()
    if err != nil {
        println(err.Error())
        return "", err
    }

    // Convert users to json
    minimalJSON, err := marshalJSON(toMinimalUsers(users))
    if err != nil {
        println(err.Error())
        return "", err
    }

    // Construct feed data and feed it to templates
    feedData := FeedData{Transactions:transactions, UsersJSON:minimalJSON, UserCount:len(users)}
    feed, err := makeHtml("../templates/feed.template", feedData)
    if err != nil {
        return "", err
    }
    return feed, nil
}

func getRegisterUser(ctx RequestContext) {
    inputForm, err  := ctx.makeHtmlWithHeader("../templates/register.template", PageData{})
    if err != nil {
        ctx.badRequestRaw("Internal error rendering page")
        return
    }
    pageData        := makePageData("Register", inputForm, []Link{{Url:"/static/styles/global.css"}}, []Link{{Url:"/static/scripts/global.js"}})
    ctx.successPage(pageData)
}

func getLogin(ctx RequestContext) {
    if !ctx.isUserLoggedIn() {
        inputForm, err  := ctx.makeHtmlWithHeader("../templates/login.template", PageData{})
        if err != nil {
            ctx.badRequestRaw("Internal error rendering page")
            return
        }
        pageData        := makePageData("Login", inputForm, []Link{{Url:"/static/styles/global.css"}}, []Link{{Url:"/static/scripts/global.js"}})
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

func getUsers(ctx RequestContext) {
    response := makeJSONResponse("")
    if !ctx.isUserLoggedIn() {
        response.Message = "Not authorized"
        ctx.notAuthorizedJSON(response)
        return
    }

    // Get all users
    users, err     := ctx.getAllUsers()
    if err != nil {
        ctx.internalErrorJSON()
        return
    }
    minimalUsers := toMinimalUsers(users)

    // Convert users to JSON
    usersJSON, err := marshalJSON(minimalUsers)
    if err != nil {
        ctx.internalErrorJSON()
        return
    }

    // Return users
    response.Message = usersJSON
    ctx.successJSON(response)
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