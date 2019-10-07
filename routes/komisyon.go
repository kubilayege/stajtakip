package routes

import (
	"net/http"
	"stajtakip/database"
	"stajtakip/templates"

	"github.com/sirupsen/logrus"
)

var tpl_komisyon = templates.Load("templates/komisyon.html")

type KomisyonVars struct {
	Uyeler []database.Komisyon
}

type KomisyonListesi struct {
	Conn *database.Connection
}

func (sh KomisyonListesi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Geçersiz metod!", http.StatusNotFound)
		return
	}

	code := http.StatusOK
	data := templates.NewMain("StajTakip - Komisyon")

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Warn("Form okunamadı!")
			w.WriteHeader(http.StatusBadRequest)
			sablonHatasi(w, tpl_komisyon.ExecuteTemplate(w, "main", data.Error("Form okunamadı!")))
			return
		}

		gorev, err := formStr(r.PostFormValue("gorev"))
		if err != nil {
			http.Error(w, "Eksik parametre!", http.StatusBadRequest)
			return
		}

		if gorev == "ekle" {
			code, data = sh.Ekle(w, r, data)
		} else if gorev == "guncelle" {
			code, data = sh.Guncelle(w, r, data)
		} else if gorev == "sil" {
			code, data = sh.Sil(w, r, data)
		}
	}

	if uyeler, err := database.KomisyonListesi(sh.Conn); err == nil {
		data.Vars = KomisyonVars{uyeler}
	} else {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Komisyon listesi yüklenemedi!")
		w.WriteHeader(http.StatusInternalServerError)
		sablonHatasi(w, tpl_komisyon.ExecuteTemplate(w, "main", data.Error("Komisyon listesi yüklenemedi!")))
		return
	}

	w.WriteHeader(code)
	sablonHatasi(w, tpl_komisyon.ExecuteTemplate(w, "main", data))
}

func (sh KomisyonListesi) Ekle(w http.ResponseWriter, r *http.Request, data templates.Main) (int, templates.Main) {
	var adSoyad string
	var err error

	adSoyad, err = formStr(r.PostFormValue("adSoyad"))
	if err != nil {
		return http.StatusBadRequest, data.Warning("AdSoyad eksik!")
	}

	kom := database.Komisyon{adSoyad, true}

	if err := kom.Insert(sh.Conn); err != nil {
		return http.StatusInternalServerError, data.Error(err.Error())
	}

	return http.StatusOK, data.Info("Komisyon üyesi veritabanına başarıyla eklendi!")
}

func (sh KomisyonListesi) Guncelle(w http.ResponseWriter, r *http.Request, data templates.Main) (int, templates.Main) {
	adSoyad, err := formStr(r.PostFormValue("adSoyad"))
	if err != nil {
		return http.StatusBadRequest, data.Warning("AdSoyad eksik!")
	}

	dahilMi := r.PostFormValue("dahil") != ""
	kom := database.Komisyon{adSoyad, dahilMi}

	if err := kom.Update(sh.Conn); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Komisyon üyesi güncellenirken veritabanında bir hata oluştu!")
		return http.StatusInternalServerError, data.Error("Veritabanında bir hata oluştu!")
	}

	return http.StatusOK, data.Info("Komisyon üyesi güncellendi!")
}

func (sh KomisyonListesi) Sil(w http.ResponseWriter, r *http.Request, data templates.Main) (int, templates.Main) {
	adSoyad, err := formStr(r.PostFormValue("adSoyad"))
	if err != nil {
		return http.StatusBadRequest, data.Warning("AdSoyad eksik!")
	}

	kom := database.Komisyon{adSoyad, true}

	if err := kom.Delete(sh.Conn); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Komisyon üyesi silinirken veritabanında bir hata oluştu!")
		return http.StatusInternalServerError, data.Error("Veritabanında bir hata oluştu!")
	}

	return http.StatusOK, data.Info("Komisyon üyesi veritabanından başarıyla silindi!")
}
