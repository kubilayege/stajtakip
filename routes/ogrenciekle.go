package routes

import (
	"net/http"
	"stajtakip/database"
	"stajtakip/templates"

	"github.com/sirupsen/logrus"
)

var tpl_ogrenci_ekle = templates.Load("templates/ogrenci-ekle.html")

type OgrenciEkleVars struct {
	No int
}

type OgrenciEkle struct {
	Conn *database.Connection
}

func (sh OgrenciEkle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		http.Error(w, "Geçersiz metod!", http.StatusNotFound)
		return
	}

	data := templates.NewMain("StajTakip - Öğrenci Ekle")
	data.Vars = OgrenciEkleVars{}

	if r.Method == http.MethodGet {
		sablonHatasi(w, tpl_ogrenci_ekle.ExecuteTemplate(w, "main", data))
		return
	}

	if err := r.ParseForm(); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Warn("Öğrenci ekleme formu okunamadı!")
		w.WriteHeader(http.StatusBadRequest)
		sablonHatasi(w, tpl_ogrenci_ekle.ExecuteTemplate(w, "main", data.Error("Öğrenci bilgileri formunda bir hata var!")))
		return
	}

	var no, ogretim int
	var ad, soyad string
	var err error

	ad, err = formStr(r.PostFormValue("ad"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		sablonHatasi(w, tpl_ogrenci_ekle.ExecuteTemplate(w, "main", data.Warning("Öğrenci adı eksik veya yanlış!")))
		return
	}

	soyad, err = formStr(r.PostFormValue("soyad"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		sablonHatasi(w, tpl_ogrenci_ekle.ExecuteTemplate(w, "main", data.Warning("Öğrenci soyadı eksik veya yanlış!")))
		return
	}

	no, err = formSayi(r.PostFormValue("no"))
	if err != nil || (no < 0) {
		w.WriteHeader(http.StatusBadRequest)
		sablonHatasi(w, tpl_ogrenci_ekle.ExecuteTemplate(w, "main", data.Warning("Öğrenci no eksik veya yanlış!")))
		return
	}

	ogretim, err = formSayi(r.PostFormValue("ogretim"))
	if err != nil || (ogretim != 2 && ogretim != 1) {
		w.WriteHeader(http.StatusBadRequest)
		sablonHatasi(w, tpl_ogrenci_ekle.ExecuteTemplate(w, "main", data.Warning("Öğretim eksik veya yanlış!")))
		return
	}

	ogr := database.Ogrenci{no, ad, soyad, ogretim}
	if err := ogr.Insert(sh.Conn); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sablonHatasi(w, tpl_ogrenci_ekle.ExecuteTemplate(w, "main", data.Error(err.Error())))
		return
	}

	data.Vars = OgrenciEkleVars{no}

	w.WriteHeader(http.StatusOK)
	sablonHatasi(w, tpl_ogrenci_ekle.ExecuteTemplate(w, "main", data.Info("Öğrenci veritabanına başarıyla eklendi!")))
}
