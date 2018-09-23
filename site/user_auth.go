package main

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

    // Try to get the user by their username or email
    user, err := ctx.getUserBy("username", providedUsername)
    if err != nil {
        user, err = ctx.getUserBy("email", providedUsername)
        if err != nil {
            return _invalid, false
        }
    }

    // Compare existing password hash and provided password
    isCorrectPassword := validatePassword(providedPassword, user.PasswordHash)

    // Check to see if the password is correct. If so, set the user id session value.
    if isCorrectPassword {
        ctx.setSessionUserId(user.Id)
        return "", true
    } else {
        return _invalid, false
    }
}

func (ctx *RequestContext) logoutUser() {
    userId := ctx.getUserId()
    if userId <= 0 {
        // User was never logged in
        return
    }
    ctx.removeSession("UserId")
}

func (ctx *RequestContext) setSessionUserId(userId int) {
    ctx.setSession("UserId", userId)
}

func isUserAuthForTransactionEdit(userId int, transaction Transaction) bool {
    return transaction.UserId == userId
}