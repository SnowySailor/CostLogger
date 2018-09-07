package main

import (
    "golang.org/x/crypto/bcrypt"
)


func (ctx *RequestContext) attemptUserLogin() (string, bool) {
    _invalid := "Invalid username or password."
    userId := ctx.getUserId()
    if userId > 0 {
        // User is already authenticated
        return "", true
    }

    // Get posted form fields
    providedUsername, _   := ctx.getFormValue("username")
    providedPassword, _   := ctx.getFormValue("password")
    providedPasswordBytes := []byte(providedPassword)

    // Try to get the user by their username or email
    user, ok := ctx.getUserBy("username", providedUsername)
    if !ok {
        user, ok = ctx.getUserBy("email", providedUsername)
        if !ok {
            return _invalid, false
        }
    }

    // Compare existing password hash and provided password
    existingPassword := []byte(user.PasswordHash)
    comparison := bcrypt.CompareHashAndPassword(existingPassword, providedPasswordBytes)

    // Check to see if the password is correct. If so, set the user id session value.
    if comparison == nil {
        ctx.setSessionUserId(user.Id)
        return "", true
    } else {
        return _invalid, false
    }
}

func (ctx *RequestContext) setSessionUserId(userId int) {
    ctx.setSession("UserId", userId)
}

func (ctx *RequestContext) getUserBy(field string, value string) (User, bool) {
    var user User
    return user, false
}