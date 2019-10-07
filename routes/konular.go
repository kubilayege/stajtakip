package routes

import (
	"net/http"
	"stajtakip/database"
	"stajtakip/templates"

	"github.com/sirupsen/logrus"
)

var tpl_konu_listesi = templates.Load("templates/konular.html")

type KonularVars struct {
	Konular []database.Konu
}

type KonuListesi struct {
	Conn *database.Connection
}

func (sh KonuListesi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Geçersiz metod!", http.StatusNotFound)
		return
	}

	code := http.StatusOK
	data := templates.NewMain("StajTakip - Konular")

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Warn("Form okunamadı!")
			w.WriteHeader(http.StatusBadRequest)
			sablonHatasi(w, tpl_konu_listesi.ExecuteTemplate(w, "main", data.Error("Form okunamadı!")))
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

	if konular, err := database.KonuListesi(sh.Conn); err == nil {
		data.Vars = KonularVars{konular}
	} else {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Konu listesi yüklenemedi!")
		w.WriteHeader(http.StatusInternalServerError)
		sablonHatasi(w, tpl_konu_listesi.ExecuteTemplate(w, "main", data.Error("Konu listesi yüklenemedi!")))
		return
	}

	w.WriteHeader(code)
	sablonHatasi(w, tpl_konu_listesi.ExecuteTemplate(w, "main", data))
}

func (sh KonuListesi) Ekle(w http.ResponseWriter, r *http.Request, data templates.Main) (int, templates.Main) {
	baslik, err := formStr(r.PostFormValue("baslik"))
	if err != nil {
		return http.StatusBadRequest, data.Warning("Baslik eksik!")
	}

	konu := database.Konu{baslik, true}
	if err := konu.Insert(sh.Conn); err != nil {
		return http.StatusInternalServerError, data.Error(err.Error())
	}

	return http.StatusOK, data.Info("Konu veritabanına başarıyla eklendi!")
}

func (sh KonuListesi) Guncelle(w http.ResponseWriter, r *http.Request, data templates.Main) (int, templates.Main) {
	baslik, err := formStr(r.PostFormValue("baslik"))
	if err != nil {
		return http.StatusBadRequest, data.Warning("Baslik eksik!")
	}

	aktifMidir := r.PostFormValue("aktif") != ""
	konu := database.Konu{baslik, aktifMidir}

	if err := konu.Update(sh.Conn); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Konu guncellenirken veritabanında bir hata oluştu!")
		return http.StatusInternalServerError, data.Error("Veritabanında bir hata oluştu!")
	}

	return http.StatusOK, data.Info("Konu güncellendi!")
}

func (sh KonuListesi) Sil(w http.ResponseWriter, r *http.Request, data templates.Main) (int, templates.Main) {
	baslik, err := formStr(r.PostFormValue("baslik"))
	if err != nil {
		return http.StatusBadRequest, data.Warning("Baslik eksik!")
	}

	konu := database.Konu{baslik, true}

	if err := konu.Delete(sh.Conn); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Konu silinirken veritabanında bir hata oluştu!")
		return http.StatusInternalServerError, data.Error("Veritabanında bir hata oluştu!")
	}

	return http.StatusOK, data.Info("Konu veritabanından başarıyla silindi!")
}
