package main

import (
    "net/http"
    "net/url"
    "bytes"
)

func firstOrDefault(l []string) string {
    if len(l) == 0 {
        return ""
    }
    return l[0]
}

func filter(l []string, f func(string) bool) []string {
    rl := make([]string, 0)
    for _, v := range l {
        if f(v) {
            rl = append(rl, v)
        }
    }
    return rl
}

func getQueryParams(r http.Request) url.Values {
    return r.URL.Query()
}

func printMap(m map[string]string) {
    for i, j := range m {
        print("\"" + i + "\"")
        println(": \"" + j + "\"")
    }
}

func printStrLStrMap(m map[string][]string) {
    for i, j := range m {
        print("\"" + i + "\"")
        println(": " + strListToStr(j))
    }
}

func printList(l []string) {
    println(strListToStr(l))
}

func strListToStr(l []string) string {
    var buffer bytes.Buffer
    if len(l) == 0 {
        buffer.WriteString("[]")
        return buffer.String()
    }
    buffer.WriteString("[")
    for i := 0; i < len(l)-1; i++ {
        buffer.WriteString("\"" + l[i] + "\", ")
    }
    buffer.WriteString("\"" + l[len(l)-1] + "\"]")
    return buffer.String()
}

func denullStrList(l []string) []string {
    return filter(l, func(v string) bool {
        return v != ""
    })
}