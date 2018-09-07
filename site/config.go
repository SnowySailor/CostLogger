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
    if conf.WebConfig.MaxUploadSize == 0 {
        conf.WebConfig.MaxUploadSize = 32 << 20
    }
    return make([]string, 0)
}

func (conf *AppConfig) validateDatabaseConfig() []string {
    if conf.DatabaseConfig.Port == 0 {
        conf.DatabaseConfig.Port = 5432
    }
    var err []string
    if conf.DatabaseConfig.Host == "" {
        err = append(err, "No database host provided")
    }
    if conf.DatabaseConfig.Username == "" {
        err = append(err, "No database user provided")
    }
    if conf.DatabaseConfig.Database == "" {
        err = append(err, "No database provided")
    }
    return err
}

func (conf *AppConfig) validateRedisConfig() []string {
    if conf.RedisConfig.Port == 0 {
        conf.RedisConfig.Port = 6379
    }
    var err []string
    if conf.RedisConfig.Host == "" {
        err = append(err, "No Redis host provided")
    }
    return err
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