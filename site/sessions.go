package main

func (ctx *RequestContext) getSession(key string) (interface{}, bool) {
    sess := ctx.session
    val, ok := sess.Values[key]
    return val, ok
}

func (ctx *RequestContext) getSessionInt(key string) (int, bool) {
    if val, ok := ctx.getSession(key); ok {
        valInt, okInt := val.(int)
        return valInt, okInt
    } else {
        return 0, false
    }
}

func (ctx *RequestContext) getSessionString(key string) (string, bool) {
    if val, ok := ctx.getSession(key); ok {
        valStr, okStr := val.(string)
        return valStr, okStr
    } else {
        return "", false
    }
}

func (ctx *RequestContext) getSessionBool(key string) (bool, bool) {
    if val, ok := ctx.getSession(key); ok {
        valBool, okBool := val.(bool)
        return valBool, okBool
    } else {
        return false, false
    }
}

func (ctx *RequestContext) setSession(key string, val interface{}) {
    sess := ctx.session
    sess.Values[key] = val
    sess.Save(ctx.request, ctx.response)
}

func (ctx *RequestContext) removeSession(key string) {
    sess := ctx.session
    delete(sess.Values, key)
    sess.Save(ctx.request, ctx.response)
}