package main

import ("fmt")

func getHome(ctx RequestContext) {
    user, err := ctx.getUserBy("username", "usera")
    output := ""
    if err == nil {
        output = fmt.Sprintf("Got user %v", user)
    } else {
        output = "No user"
    }
    user = User {
        Username: "usera",
        DisplayName: "usera123",
        Email: "Hello@user.com",
    }
    userId, err := ctx.insertUser(user)
    if err != nil {
        panic(err)
    } else {
        output = output + fmt.Sprintf(", new user %v", userId)
    }
    ctx.successPage("<h3>Get home: " + output + "</h3>")
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
