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

func validateDatabaseConfig(conf *AppConfig) string {
    if conf.DatabaseConfig.Port == 0 {
        conf.DatabaseConfig.Port = 5432
    }
    if conf.DatabaseConfig.Host == "" {
        return "No database host provided in ./secrets.yaml"
    }
    if conf.DatabaseConfig.Username == "" {
        return "No database user provided in ./secrets.yaml"
    }
    if conf.DatabaseConfig.Database == "" {
        return "No database provided in ./secrets.yaml"
    }
    return ""
}
