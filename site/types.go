package main

import (
    "net/http"
    "database/sql"
    "html/template"
    "github.com/go-redis/redis"
    "github.com/gorilla/sessions"
)

// Type alias for redis so that I can define methods on the redis.Client type
type Redis redis.Client

// Request context
type RequestContext struct {
    request   *http.Request
    response  http.ResponseWriter
    method    string
    userId    int
    routes    []string
    database  *sql.DB
    session   *sessions.Session
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