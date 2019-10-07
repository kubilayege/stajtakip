package routes

import "net/http"

// Web klasörünü sunar
var Assets http.Handler = http.StripPrefix("/assets/", http.FileServer(http.Dir("assets")))
