package main

import (
    "net/http"
    "database/sql"
    "html/template"
    "github.com/go-redis/redis"
    "github.com/gorilla/sessions"
    "time"
)

// Type alias so methods can be defined on non-local types
type Redis        redis.Client
type Session      sessions.Session
type flint        int
type ReadOnlyBool bool

// Special handling of json unmarshaling for ReadOnlyBool type
func (ReadOnlyBool) UnmarshalJSON([]byte) error { return nil }

// Request context
type RequestContext struct {
    request    *http.Request
    response   http.ResponseWriter
    method     string
    userId     int
    routes     []string
    database   *sql.DB
    session    *sessions.Session
    redis      *redis.Client
}

type JSONResponse struct {
    Message     string `json:"Message"`
    RedirectUrl string `json:"RedirectUrl"`
}

// Configuration
type AppConfig struct {
    DatabaseConfig struct {
        Host     string `yaml:"host"`
        Port     int    `yaml:"port"`
        Username string `yaml:"username"`
        Password string `yaml:"password"`
        Database string `yaml:"database"`
    } `yaml:"databaseconfig"`
    RedisConfig struct {
        Host     string `yaml:"host"`
        Port     int    `yaml:"port"`
        Password string `yaml:"password"`
        Database int    `yaml:"database"`
    } `yaml:"redisconfig"`
    SessionConfig struct {
        SecretKey      string `yaml:"secretkey"`
        Name           string `yaml:"name"`
        MaxAge         int    `yaml:"maxage"`
        HttpOnly       string `yaml:"httponly"`
        internHttpOnly bool
    } `yaml:"sessionconfig"`
    WebConfig struct {
        Port              int   `yaml:"port"`
        MaxUploadSize     int64 `yaml:"maxuploadsize"`
        PasswordStrength  int   `yaml:"passwordstrength"`
        MinPasswordLength int   `yaml:"minpasswordlength"`
    } `yaml:"webconfig"`
}

// Types for rendering pages with templates
type PageData struct {
    Title     string
    StyleSrc  []Link
    ScriptSrc []Link
    Body      template.HTML
}

type Link struct {
    Url string
}

type HeaderData struct {
    IsUserLoggedIn bool
    DisplayName    string
}

type FeedData struct {
    Transactions  []PageTransaction
    UsersJSON     string
    CurrentUserId int
    AmountsOwed   map[string]flint
    AmountsOwedToThisUser map[string]flint
}

type PageTransaction struct {
    Id              int
    Amount          flint
    Comments        string
    CreateDate      time.Time
    UserId          int
    DisplayName     string
    Username        string
    InvolvedUsers   []PageTransactionUser
    LastUpdateDate  time.Time
}

type PageTransactionUser struct {
    UserId             int
    PercentInvolvement flint
    AmountInvolvement  flint
    Username           string
    DisplayName        string
    IsPaid             bool
}

// Application data types
type User struct {
    Id           int
    Username     string
    DisplayName  string
    Email        string
    PasswordHash string
}

type MinimalUser struct {
    Id          int
    Username    string
    DisplayName string
}

type Transaction struct {
    Id              int               `json:"id"`
    Amount          flint             `json:"amount"`// Example: 5049 = $50.49
    Comments        string            `json:"comments"`
    CreateDate      time.Time
    UserId          int
    InvolvedUsers   []TransactionUser `json:"involvedusers"`
    LastUpdateDate  time.Time
}

type TransactionUser struct {
    UserId             int           `json:"userid"`
    TransactionId      int           `json:"transactionid"`
    PercentInvolvement flint         `json:"percentinvolvement"` // Example: 5049 = 50.49%
    IsPaid             ReadOnlyBool  `json:"ispaid"`
}
