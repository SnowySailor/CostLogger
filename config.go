package main

import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "fmt"
)

func (conf *AppConfig) getAppConfig() {
    if file, err := ioutil.ReadFile("./secrets.yaml"); err != nil {
        panic(fmt.Sprintf("Error reading config: %v\n", err))
    } else {
        if err := yaml.Unmarshal(file, conf); err != nil {
            panic(fmt.Sprintf("Error parsing ./secrets.yaml: %v\n", err))
        }
        err := validateDatabaseConfig(conf)
        if err != "" {
            panic(err + "\n")
        }
    }
}