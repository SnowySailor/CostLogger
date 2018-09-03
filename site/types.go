package main

import (
    "net/http"
    "database/sql"
)

type RequestContext struct {
    request   *http.Request
    response  http.ResponseWriter
    method    string
    userId    int
    routes    []string
    database  *sql.DB
}

type AppConfig struct {
    DatabaseConfig struct {
        Host     string `yaml:"host"`
        Port     int    `yaml:"port"`
        Username string `yaml:"username"`
        Password string `yaml:"password"`
        Database string `yaml:"database"`
    } `yaml:"databaseconfig"`
}

type PageData struct {
    Title     string
    StyleSrc  []Link
    ScriptSrc []Link
    Body      string
}

type Link struct {
    Url string
}