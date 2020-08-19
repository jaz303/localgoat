package main

import (
	"net/http"
	"os"
)

type Action interface{}

type messageAction struct {
	StatusCode int
	Message    string
}

var NotFound = &messageAction{http.StatusNotFound, "Not Found"}
var InternalServerError = &messageAction{http.StatusInternalServerError, "Internal Server Error"}

type serveFileAction struct {
	TargetFile   string
	Stat         os.FileInfo
	ExtraHeaders map[string]string
}

type proxyAction struct {
	ProxyName string
}
