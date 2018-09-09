package main

import (
    "fmt"
)

func (ctx RequestContext) notFoundPage(msg string) {
    // TOOD: 
    ctx.writeResponse(msg, 404, "text/html")
}

// HTML writers

func (ctx RequestContext) successPage(data PageData) {
    result := makeHtmlWithTemplate("../templates/page_wrapper.template", data)
    ctx.writeResponse(result, 200, "text/html")
}

func (ctx RequestContext) badRequestPage(data PageData) {
    result := makeHtmlWithTemplate("../templates/page_wrapper.template", data)
    ctx.writeResponse(result, 400, "text/html")
}

// Raw writers

func (ctx RequestContext) successRaw(msg string) {
    ctx.writeResponse(msg, 200, "text/plain")
}

func (ctx RequestContext) badRequestRaw(msg string) {
    ctx.writeResponse(msg, 400, "text/plain")
}

// JSON writers

func (ctx RequestContext) successJSON(resp JSONResponse) {
    msg, err := marshalJSON(resp)
    if err != nil {
        panic(err)
    }
    ctx.writeResponse(msg, 200, "application/json")
}

func (ctx RequestContext) badRequestJSON(resp JSONResponse) {
    msg, err := marshalJSON(resp)
    if err != nil {
        panic(err)
    }
    ctx.writeResponse(msg, 400, "application/json")
}

// Main writer

func (ctx RequestContext) writeResponse(msg string, status int, contentType string) {
    ctx.response.Header().Set("Content-Type", contentType)
    ctx.response.WriteHeader(status)
    fmt.Fprintf(ctx.response, msg)
}
