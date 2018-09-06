package main

import (
    "net/http"
    "database/sql"
    "html/template"
    "github.com/go-redis/redis"
    "github.com/gorilla/sessions"
)

// Type alias so methods can be defined on non-local types
type Redis redis.Client
type Session sessions.Session

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
    }
    SessionConfig struct {
        SessionSecretKey string `yaml:"sessionsecretkey"`
        SessionName      string `yaml:"sessionname"`
    }
    WebConfig struct {
        MaxUploadSize int64 `yaml:"maxuploadsize"`
    }
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