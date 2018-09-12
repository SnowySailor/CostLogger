package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

func getDatabaseConnection() (*sql.DB, error) {
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        config.DatabaseConfig.Host, config.DatabaseConfig.Port, config.DatabaseConfig.Username,
        config.DatabaseConfig.Password, config.DatabaseConfig.Database)
    if db, err := sql.Open("postgres", psqlInfo); err != nil {
        return nil, err
    } else {
        if err = db.Ping(); err != nil {
            return nil, err
        } else {
            return db, nil
        }
    }
}

// USER QUERIES
// //////////////////////////////////////////////////////////////

func (ctx *RequestContext) doesUserExist(id int) bool {
    var uId int
    query := "select id from app_user where id = $1"
    row   := ctx.database.QueryRow(query, id)
    err   := row.Scan(&uId)
    return (err == nil)
}

// Get a single user by their id
func (ctx *RequestContext) getUser(id int) (User, error) {
    return ctx.getUserBy("id", id)
}

// Get a single user by a provided field as long as the field is supported
func (ctx *RequestContext) getUserBy(field string, value interface{}) (User, error) {
    var user User
    allowedFields := []string{"id", "username", "display_name", "email"}
    if !strInList(strToLower(field), allowedFields) {
        return user, makeError("Invalid search field.")
    }

    // Query
    query := fmt.Sprintf("select id, username, display_name, email, password_hash from app_user where %s = $1", field)
    row   := ctx.database.QueryRow(query, value)
    err   := row.Scan(&user.Id, &user.Username, &user.DisplayName, &user.Email, &user.PasswordHash)

    return user, err
}

// Insert a new user into the database
func (ctx *RequestContext) insertUser(user User) (int, error) {
    var userId int
    // Query
    query := "insert into app_user (username, display_name, email, password_hash) values ($1, $2, $3, $4) returning id"
    row   := ctx.database.QueryRow(query, user.Username, user.DisplayName, user.Email, user.PasswordHash)
    err   := row.Scan(&userId)

    return userId, err
}

// //////////////////////////////////////////////////////////////


// TRANSACTION QUERIES
// //////////////////////////////////////////////////////////////

// Get a single transaction
func (ctx *RequestContext) getTransaction(transactionId int) (Transaction, error) {
    var transaction Transaction
    // Query
    query := "select id, amount, comments, createdate, userid, lastupdatetime from transaction where id = $1"
    row   := ctx.database.QueryRow(query, transactionId)
    err   := row.Scan(&transaction.Id, &transaction.Amount, &transaction.Comments, &transaction.CreateDate, &transaction.UserId, &transaction.LastUpdateTime)

    if err != nil {
        return transaction, err
    }

    involvedUsers, err        := ctx.getTransactionUsers(transactionId)
    transaction.InvolvedUsers = involvedUsers

    return transaction, err
}

// Insert a new transaction into the database
func (ctx *RequestContext) insertTransaction(transaction Transaction) (int, error) {
    var transactionId int

    // Begin database transaction
    tx, err := ctx.database.Begin()
    if err != nil {
        return transactionId, err
    }

    // Insert main transaction
    query   := "insert into transaction (amount, comments, user_id) values ($1, $2, $3) returning id"
    row     := ctx.database.QueryRow(query, transaction.Amount, transaction.Comments, transaction.UserId)
    err     = row.Scan(&transactionId)

    // Insert all involved users
    for i := 0; i < len(transaction.InvolvedUsers); i = i+1 {
        err = ctx.insertTransactionUser(transaction.InvolvedUsers[i])
        if err != nil {
            err = tx.Rollback()
            if err != nil {
                panic(err)
            }
            return transactionId, err
        }
    }

    // Commit database transaction
    err = tx.Commit()
    if err != nil {
        panic(err)
    }

    return transactionId, err
}

// Get the users that are involved in a transaction
func (ctx *RequestContext) getTransactionUsers(transactionId int) ([]TransactionUser, error) {
    users     := make([]TransactionUser, 0)
    query     := "select user_id, transaction_id, percentage from transaction_user where transaction_id = $1"
    rows, err := ctx.database.Query(query, transactionId)
    if err != nil {
        return users, err
    }
    defer rows.Close()
    for rows.Next() {
        var transactionUser TransactionUser
        err = rows.Scan(&transactionUser.UserId, transactionUser.PercentInvolvement)
        if err != nil {
            return users, err
        }
        users = append(users, transactionUser)
    }
    err = rows.Err()
    if err != nil {
        return users, err
    }
    return users, nil
}

// Insert all the users that are involved in a transaction
func (ctx *RequestContext) insertTransactionUser(user TransactionUser) error {
    query   := "insert into transaction_user (user_id, transaction_id, percentage) values ($1, $2, $3)"
    _, err := ctx.database.Exec(query, user.UserId, user.TransactionId, user.PercentInvolvement)
    return err
}

// //////////////////////////////////////////////////////////////
