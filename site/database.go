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

    if err != nil {
        return user, err
    } else {
        return user, err
    }
}

func (ctx *RequestContext) insertUser(user User) (int, error) {
    var userId int
    // Query
    query := "insert into app_user (username, display_name, email, password_hash) values ($1, $2, $3, $4) returning id"
    row   := ctx.database.QueryRow(query, user.Username, user.DisplayName, user.Email, user.PasswordHash)
    err   := row.Scan(&userId)

    return userId, err
}