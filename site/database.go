package main

import (
    "database/sql"
    "fmt"
    "sort"
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

func (ctx RequestContext) getAllUsers() ([]User, error) {
    var users []User

    query     := "select id, username, display_name, email, password_hash from app_user"
    rows, err := ctx.database.Query(query)

    if err != nil {
        return users, err
    }

    defer rows.Close()
    for rows.Next() {
        var user User
        err := rows.Scan(&user.Id, &user.Username, &user.DisplayName, &user.Email, &user.PasswordHash)
        if err != nil {
            return users, err
        }
        users = append(users, user)
    }
    err = rows.Err()
    if err != nil {
        return users, err
    }
    return users, nil
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
    query := "select id, amount, comments, create_date, user_id, last_update_time from transaction where id = $1"
    row   := ctx.database.QueryRow(query, transactionId)
    err   := row.Scan(&transaction.Id, &transaction.Amount, &transaction.Comments, &transaction.CreateDate, &transaction.UserId, &transaction.LastUpdateDate)

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
    query     := "select user_id, percentage from transaction_user where transaction_id = $1"
    rows, err := ctx.database.Query(query, transactionId)
    if err != nil {
        return users, err
    }
    defer rows.Close()
    for rows.Next() {
        var transactionUser TransactionUser
        transactionUser.TransactionId = transactionId
        err = rows.Scan(&transactionUser.UserId, &transactionUser.PercentInvolvement)
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

// Get all transactions that a user is involved in
func (ctx *RequestContext) getUserTransactions(userId int) ([]Transaction, error) {
    var transactions []Transaction
    var commentsNull sql.NullString

    // Queries
    userTransactionsQ :=
        `select id, amount, comments, create_date, user_id, last_update_date from transaction
        where user_id = $1
        
        union
        
        select T.id, T.amount, T.comments, T.create_date, T.user_id, T.last_update_date from transaction T
        inner join transaction_user TU on TU.transaction_id = T.id
        where TU.user_id = $1`

    // Get the transactions that the user has created
    rows, err := ctx.database.Query(userTransactionsQ, userId)
    if err != nil {
        return make([]Transaction, 0), err
    }

    defer rows.Close()
    for rows.Next() {
        var transaction Transaction
        // Get the transaction
        err = rows.Scan(&transaction.Id, &transaction.Amount, &commentsNull, &transaction.CreateDate, &transaction.UserId, &transaction.LastUpdateDate)
        if err != nil {
            return make([]Transaction, 0), err
        }
        if commentsNull.Valid {
            transaction.Comments = commentsNull.String
        }

        // Get involved users
        transactionUsers, err := ctx.getTransactionUsers(transaction.Id)
        if err != nil {
            return make([]Transaction, 0), err
        }
        transaction.InvolvedUsers = transactionUsers
        
        // Update our list
        transactions = append(transactions, transaction)
    }

    // Sort by CreateDate
    sort.Slice(transactions, func(i int, j int) bool {
        return transactions[i].CreateDate.Before(transactions[j].CreateDate)
    })
    return transactions, nil
}

// //////////////////////////////////////////////////////////////
