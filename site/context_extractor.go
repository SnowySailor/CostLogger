package main

import (
    "fmt"
    "encoding/json"
    "io/ioutil"
    "time"
    "strconv"
)

func (ctx *RequestContext) extractNewUser() (User, error) {
    var user User
    username, ok := ctx.getFormValue("username")
    if !ok || username == "" {
        return user, makeError("Username not provided")
    }

    email, ok := ctx.getFormValue("email")
    if !ok || email == "" {
        return user, makeError("Email not provided")
    }

    displayName, ok := ctx.getFormValue("displayName")
    if !ok || displayName == "" {
        return user, makeError("Display name not provided")
    }

    password, ok := ctx.getFormValue("password")
    if !ok || password == "" {
        return user, makeError("Password not provided")
    }

    passwordVerify, ok := ctx.getFormValue("passwordVerify")
    if !ok || passwordVerify == "" {
        return user, makeError("Password verification not provided")
    }

    if password != passwordVerify {
        return user, makeError("Passwords do not match")
    }

    if len(password) < config.WebConfig.MinPasswordLength {
        return user, makeError(fmt.Sprintf("Password is not long enough (minimum %v characters)", config.WebConfig.MinPasswordLength))
    }

    user = User {
        Username:     username,
        DisplayName:  displayName,
        Email:        email,
        PasswordHash: hashPassword(password),
    }
    return user, nil
}

func (ctx *RequestContext) extractTransaction() (Transaction, error) {
    var t Transaction

    // Get/parse transactionId from url
    transactionIdStr, ok := getValueString(1, ctx.routes)
    if !ok {
        return t, makeError("Invalid transaction id")
    }

    transactionId, err := strconv.Atoi(transactionIdStr)
    if err != nil {
        return t, makeError("Invalid transaction id")
    }

    now := time.Now().UTC()
    err  = json.Unmarshal(ctx.getRequestBody(), &t)
    t.LastUpdateDate = now
    t.CreateDate     = now
    t.UserId         = ctx.userId
    t.Id             = transactionId
    return t, err
}

// If existing:
    // This user id must be the same person that created the transaction
// Always:
    // Amount must be greater than $0.00
    // All userids in involved users must exist
    // All userids in involved users must be unique
    // Percentages must add up to 100
func (ctx *RequestContext) validateTransaction(transaction Transaction) error {
    // Transaction ID provided (it probably exists)
    if transaction.Id > 0 {
        existingT, err := ctx.getTransaction(transaction.Id)
        // If we get back an existing transaction, it definitely exists
        if err == nil {
            if existingT.UserId != ctx.userId {
                return makeError("Not authorized")
            } 
        }
    }

    if transaction.Amount <= 0 {
        return makeError("Amount must be greater than $0.00")
    }

    if len(transaction.InvolvedUsers) == 0 {
        return makeError("Transaction must have at least one involved user")
    }

    var userIds []int
    var percentSum int
    for _, tUser := range transaction.InvolvedUsers {
        if !ctx.doesUserExist(tUser.UserId) {
            return makeError(fmt.Sprintf("User %v does not exist", tUser.UserId))
        }
        if intInList(tUser.UserId, userIds) {
            return makeError("Involved users contains duplicate users")
        }
        userIds = append(userIds, tUser.UserId)
        percentSum = percentSum + int(tUser.PercentInvolvement)
    }
    // if percentSum != 100.00%
    if percentSum != 10000 {
        return makeError("All percent involvements must sum to 100.00%%")
    }
    return nil
}

func (ctx *RequestContext) validateNewUser(user User) error {
    if user, _ := ctx.getUserBy("username", user.Username); (user != User{}) {
        return makeError("Username already in use")
    }
    if user, _ := ctx.getUserBy("username", user.Email); (user != User{}) {
        return makeError("Email address already in use")
    }
    return nil
}

func (ctx *RequestContext) getRequestBody() []byte {
    bytes, err := ioutil.ReadAll(ctx.request.Body)
    if err != nil {
        return make([]byte, 0)
    }
    return bytes
}