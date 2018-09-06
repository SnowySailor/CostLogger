package main

func (ctx *RequestContext) attemptUserLogin() (string, bool) {
    userId := ctx.getUserId()
    if userId > 0 {
        return "", true
    }

    // Get posted form fields
    //login    := ""
    hash     := ""
    password := ""

    // Get user by login name

    // Hash password with user's salt

    // Compare hashed password with hash in database

    if hash == password {
        // Set session user id to user.id
        return "", true
    } else {
        return "Invalid password", false
    }
    return "", false
}

func (ctx *RequestContext) getFormField(name string) (string, bool) {
    return "", true
}
