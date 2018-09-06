package main

import ("fmt")

func getHome(ctx RequestContext) {
    val, ok := ctx.getSessionInt("times")
    if ok {
        ctx.setSession("times", val + 1)
    } else {
        ctx.setSession("times", 1)
    }
    val = val + 1
    ctx.successPage("<h3>Get home: " + fmt.Sprintf("%v", val) + "</h3>")
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
