package main

import (
    "fmt"
)

func (ctx RequestContext) notFoundPage(msg string) {
    ctx.writeResponse(msg, 404, "text/html")
}

func (ctx RequestContext) successPage(msg string) {
    ctx.writeResponse(msg, 200, "text/html")
}

func (ctx RequestContext) successRaw(msg string) {
    ctx.writeResponse(msg, 200, "text/plain")
}

func (ctx RequestContext) successJson(msg string) {
    ctx.writeResponse(msg, 200, "application/json")
}

func (ctx RequestContext) badRequestPage(msg string) {
    ctx.writeResponse(makeBadRequestPage(msg), 400, "text/html")
}

func (ctx RequestContext) badRequestRaw(msg string) {
    ctx.writeResponse(msg, 400, "text/plain")
}

func (ctx RequestContext) writeResponse(msg string, status int, contentType string) {
    ctx.response.Header().Set("Content-Type", contentType)
    ctx.response.WriteHeader(status)
    fmt.Fprintf(ctx.response, msg)
}
