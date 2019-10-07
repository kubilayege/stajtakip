package routes

import "net/http"

// Web klasörünü sunar
var Uploads http.Handler = http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads")))
