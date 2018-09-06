package main

func (ctx *RequestContext) getUserId() int {
    userId, exists := ctx.getSessionInt("UserId")
    if exists {
        return userId
    }
    return -1
}