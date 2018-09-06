package main

import (
    "net/http"
    "net/url"
    "bytes"
    "html/template"
    "strings"
)

func splitPathRoutes(path string) []string {
    return denullStrList(strings.Split(path, "/"))
}

func getPathRoutes(path string) []string {
    routes := splitPathRoutes(path)
    if len(routes) == 0 {
        return routes
    }
    // Split the query from the last route
    lastRoute := getLastPathRoute(path)
    // "replace" the last route in `routes` with the new lastRoute
    return append(routes[:len(routes)-1], lastRoute)
}

func getLastPathRoute(path string) string {
    // Get all routes
    routes := splitPathRoutes(path)
    if len(routes) == 0 {
        // If there were no routes returned, return empty string
        return ""
    }
    // Get the last route and the location of the beginning of the url query
    lastRoute := routes[len(routes)-1]
    queryIndex   := strings.Index(lastRoute, "?")
    if queryIndex == -1 {
        // If there is no query, we don't have to go further. Just return the last route.
        return lastRoute
    }
    // Return the last route without the query
    return lastRoute[:queryIndex]
}


func makeHtmlWithTemplate(pageData PageData, templateLocation string) string {
    var templateBytes bytes.Buffer
    t := template.Must(template.ParseFiles(templateLocation))

    if err := t.Execute(&templateBytes, pageData); err != nil {
        panic(err)
    } else {
        return templateBytes.String()
    }
}

func makePageData(title string, body string, styleSrc []Link, scriptSrc []Link) PageData {
    return PageData {
        Title:     title,
        Body:      template.HTML(body),
        StyleSrc:  styleSrc,
        ScriptSrc: scriptSrc,
    }
}

func firstOrDefault(l []string) string {
    if len(l) == 0 {
        return ""
    }
    return l[0]
}

func filterStrings(l []string, f func(string) bool) []string {
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
    return filterStrings(l, func(v string) bool {
        return v != ""
    })
}