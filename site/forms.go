package main

func (ctx *RequestContext) getFormValue(name string) (string, bool) {
    if val, ok := ctx.getFormValueMulti(name); ok {
        if len(val) > 0 {
            return val[0], true
        }
    }
    return "", false
}

func (ctx *RequestContext) getFormValueMulti(name string) ([]string, bool) {
    if val, ok := ctx.FormValueMulti(name); ok {
        return val, true
    }
    if val, ok := ctx.PostFormValueMulti(name); ok {
        return val, true
    }
    return make([]string, 0), false
}

func (ctx *RequestContext) FormValueMulti(key string) ([]string, bool) {
    r := ctx.request
    if r.Form == nil {
        r.ParseMultipartForm(config.WebConfig.MaxUploadSize)
    }
    if vs := r.Form[key]; len(vs) > 0 {
        return vs, true
    }
    return make([]string, 0), false
}

func (ctx *RequestContext) PostFormValueMulti(key string) ([]string, bool) {
    r := ctx.request
    if r.PostForm == nil {
        r.ParseMultipartForm(config.WebConfig.MaxUploadSize)
    }
    if vs := r.PostForm[key]; len(vs) > 0 {
        return vs, true
    }
    return make([]string, 0), false
}