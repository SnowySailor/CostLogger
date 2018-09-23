package main

import (
    "net/http"
    "net/url"
    "bytes"
    "html/template"
    "strings"
    "errors"
    "golang.org/x/crypto/bcrypt"
    "encoding/json"
    "strconv"
)

func getValueString(i int, l []string) (string, bool) {
    if i >= len(l) {
        return "", false
    }
    return l[i], true
}

func toMinimalUsers(users []User) []MinimalUser {
    var minimalUsers []MinimalUser
    for _, user := range users {
        minimalUsers = append(minimalUsers, MinimalUser{Id: user.Id, Username: user.Username, DisplayName: user.DisplayName})
    }
    return minimalUsers
}

func (f flint) FlintToString(baseOffset int, decimalPlaces int) string {
    strVal := strconv.Itoa(int(f))
    major  := ""
    minor  := ""
    if len(strVal) > baseOffset {
        minor = strVal[(len(strVal)-baseOffset):]
        major = strVal[:(len(strVal)-baseOffset)]
    } else {
        major = "0"
        minor = strVal
    }
    minor = trim(padString(minor, decimalPlaces, "0", true), decimalPlaces)
    ret := major
    if decimalPlaces == 0 {
        return ret
    }
    return ret + "." + minor
}

func trim(s string, l int) string {
    if l <= 0 {
        return ""
    }
    if len(s) < l {
        return s
    }
    return s[:l]
}

func padString(s string, l int, pad string, beginning bool) string {
    if len(s) >= l {
        return s
    }
    if len(pad) == 0 {
        panic("padString: Argument `pad` is empty; avoiding infinite loop")
    }
    remaining := int((l - len(s))/len(pad))
    for i := 0; i < remaining; i = i+1 {
        if beginning {
            s = pad + s
        } else {
            s = s + pad
        }
    }
    return s
}

// Takes a string (URL) and removes the leading slash if it exists
func removeLeadingSlash(str string) string {
    if len(str) == 0 {
        return str
    }
    if str[0] == '/' {
        return str[1:]
    }
    return str
}

// Takes the password the user provided and the hash to verify it with
// Returns whether the password matches the hash
func validatePassword(provided string, storedHash string) bool {
    storedHashBytes := []byte(storedHash)
    providedBytes   := []byte(provided)
    err := bcrypt.CompareHashAndPassword(storedHashBytes, providedBytes)
    return (err == nil)
}

// Takes a string password and returns the base64-encoded bcrypt-hashed version
func hashPassword(provided string) string {
    providedBytes    := []byte(provided)
    hashedBytes, err := bcrypt.GenerateFromPassword(providedBytes, config.WebConfig.PasswordStrength)
    if err != nil {
        panic(err)
    }
    return string(hashedBytes[:])
}

func marshalJSON(inter interface{}) (string, error) {
    bytes, err := json.Marshal(inter)
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}

func makeJSONResponse(msg string) JSONResponse {
    return JSONResponse {
        Message: msg,
    }
}

func makeError(msg string) error {
    return errors.New(msg)
}

func strToLower(value string) string {
    return strings.ToLower(value)
}

func intInList(value int, list []int) bool {
    for i := 0; i < len(list); i = i+1 {
        if value == list[i] {
            return true
        }
    }
    return false
}

func strInList(value string, list []string) bool {
    for i := 0; i < len(list); i = i+1 {
        if value == list[i] {
            return true
        }
    }
    return false
}

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

func (ctx *RequestContext) getUserDisplayName() string {
    if !ctx.isUserLoggedIn() {
        return ""
    }
    user, err := ctx.getUser(ctx.userId)
    if err != nil {
        return ""
    }
    return user.DisplayName
}

func (ctx *RequestContext) makeHtmlWithHeader(templateLocation string, data interface{}) (string, error) {
    var templateBytes bytes.Buffer
    t := template.Must(template.ParseFiles("../templates/header.template"))

    headerData := HeaderData {
        IsUserLoggedIn: ctx.isUserLoggedIn(),
        DisplayName: ctx.getUserDisplayName(),
    }
    if err := t.Execute(&templateBytes, headerData); err != nil {
        return "", err
    } else {
        rest, err := makeHtml(templateLocation, data)
        if err != nil {
            return "", err
        }
        return templateBytes.String() + rest, nil
    }
}

func makeHtml(templateLocation string, data interface{}) (string, error) {
    var templateBytes bytes.Buffer
    t := template.Must(template.ParseFiles(templateLocation))

    if err := t.Execute(&templateBytes, data); err != nil {
        return "", err
    } else {
        return templateBytes.String(), err
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