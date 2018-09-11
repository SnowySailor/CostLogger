package main

func (ctx *RequestContext) extractNewUser() (User, error) {
    var user User
    username, ok := ctx.getFormValue("username")
    if !ok || username == "" {
        return user, makeError("Username not provided.")
    }

    email, ok := ctx.getFormValue("email")
    if !ok || email == "" {
        return user, makeError("Email not provided.")
    }

    displayName, ok := ctx.getFormValue("displayName")
    if !ok || displayName == "" {
        return user, makeError("Display name not provided.")
    }

    password, ok := ctx.getFormValue("password")
    if !ok || password == "" {
        return user, makeError("Password not provided.")
    }

    passwordVerify, ok := ctx.getFormValue("passwordVerify")
    if !ok || passwordVerify == "" {
        return user, makeError("Password verification not provided.")
    }

    if password != passwordVerify {
        return user, makeError("Passwords do not match.")
    }

    user = User {
        Username:     username,
        DisplayName:  displayName,
        Email:        email,
        PasswordHash: hashPassword(password),
    }
    return user, nil
}