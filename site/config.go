package main

import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "fmt"
)

func (conf *AppConfig) getAppConfig() {
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