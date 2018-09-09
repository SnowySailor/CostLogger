package main

import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "fmt"
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