package routes

import (
	"net/http"
	"stajtakip/database"
	"stajtakip/templates"

	"github.com/sirupsen/logrus"
)

var tpl_ogrenci_listesi = templates.Load("templates/ogrenci-listesi.html")

type OgrenciListesiVars struct {
	Ogrenciler []database.Ogrenci
}

type OgrenciListesi struct {
	Conn *database.Connection
}

func (sh OgrenciListesi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Get metodu kullanılmalı!", http.StatusNotFound)
		return
	}

	data := templates.NewMain("StajTakip - Stajı Biten Öğrenci Listesi")

	if ogrenciler, err := database.StajiTamamOgrenciler(sh.Conn); err == nil {
		data.Vars = OgrenciListesiVars{ogrenciler}
	} else {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Öğrenci listesi yüklenemedi!")
		w.WriteHeader(http.StatusInternalServerError)
		sablonHatasi(w, tpl_ogrenci_listesi.ExecuteTemplate(w, "main", data.Error("Öğrenci listesi yüklenemedi!")))
		return
	}

	w.WriteHeader(http.StatusOK)
	sablonHatasi(w, tpl_ogrenci_listesi.ExecuteTemplate(w, "main", data))
}
