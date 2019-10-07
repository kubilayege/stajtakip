package routes

import (
	"net/http"
	"stajtakip/database"
)

type TODO struct {
	Conn *database.Connection
}

func (sh TODO) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Bu özellik henüz eklenmedi!", http.StatusInternalServerError)
}
