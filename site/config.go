package main

import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "io"
    "fmt"
    "crypto/rand"
    "encoding/base64"
)

func (conf *AppConfig) populateAppConfig() {
    if file, err := ioutil.ReadFile("../secrets.yaml"); err != nil {
        panic(fmt.Sprintf("Error reading config: %v\n", err))
    } else {
        if err := yaml.Unmarshal(file, conf); err != nil {
            panic(fmt.Sprintf("Error parsing secrets.yaml: %v\n", err))
        }
        err := conf.validateDatabaseConfig()
        err = append(err, conf.validateRedisConfig()...)
        err = append(err, conf.validateSessionConfig()...)
        if len(err) > 0 {
            panic(stringJoin(err, ", ") + " in secrets.yaml")
        }
    }
}

func (conf *AppConfig) validateWebConfig() []string {
    var errors []string
    if conf.WebConfig.MaxUploadSize == 0 {
        conf.WebConfig.MaxUploadSize = 32 << 20
    }
    if conf.WebConfig.PasswordStrength <= 0 {
        conf.WebConfig.PasswordStrength = 10
    }
    return errors
}

func (conf *AppConfig) validateDatabaseConfig() []string {
    errors := make([]string, 0)
    if conf.DatabaseConfig.Port == 0 {
        conf.DatabaseConfig.Port = 5432
    }
    if conf.DatabaseConfig.Host == "" {
        errors = append(errors, "No database host provided")
    }
    if conf.DatabaseConfig.Username == "" {
        errors = append(errors, "No database user provided")
    }
    if conf.DatabaseConfig.Database == "" {
        errors = append(errors, "No database provided")
    }
    return errors
}

func (conf *AppConfig) validateRedisConfig() []string {
    var errors []string
    if conf.RedisConfig.Port == 0 {
        conf.RedisConfig.Port = 6379
    }
    if conf.RedisConfig.Host == "" {
        errors = append(errors, "No Redis host provided")
    }
    return errors
}

func (conf *AppConfig) validateSessionConfig() []string {
    var errors []string
    if conf.SessionConfig.Name == "" {
        conf.SessionConfig.Name = "sessionid"
    }
    if conf.SessionConfig.SecretKey == "" {
        randBytes := make([]byte, 256)
        _, err := io.ReadFull(rand.Reader, randBytes)
        if err != nil {
            errors = append(errors, err.Error())
        } else {
            conf.SessionConfig.SecretKey = base64.StdEncoding.EncodeToString(randBytes)
        }
    }
    if conf.SessionConfig.MaxAge == 0 {
        conf.SessionConfig.MaxAge = 3600
    }
    conf.SessionConfig.internHttpOnly = strToLower(conf.SessionConfig.HttpOnly) != "off"
    return errors
}

func stringJoin(strings []string, delimiter string) string {
    res := ""
    last := len(strings) - 1
    for i := 0; i < len(strings); i++ {
        res = res + strings[i]
        if i != last {
            res = res + delimiter
        }
    }
    return res
}