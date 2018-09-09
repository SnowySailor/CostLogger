package main

import (
    "github.com/gorilla/sessions"
)

var config AppConfig
var store *sessions.CookieStore
