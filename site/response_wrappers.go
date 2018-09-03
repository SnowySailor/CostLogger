package main

func makeBadRequestPage(msg string) string {
    return "<!DOCTYPE html><html><body>" + msg + "</body></html>"
}