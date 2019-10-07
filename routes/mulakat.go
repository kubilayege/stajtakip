package routes

import (
	"net/http"
	"stajtakip/database"
	"stajtakip/templates"

	"time"

	"github.com/sirupsen/logrus"
)

var tpl_mulakat_listesi = templates.Load("templates/mulakat-listesi.html")

type MulakatListesiVars struct {
	Mulakatlar []database.MulakatOgrenci
		Komisyon []database.Komisyon
}

type MulakatListesi struct {
	Conn *database.Connection
}

func (sh MulakatListesi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		http.Error(w, "Geçersiz metod!", http.StatusNotFound)
		return
	}

	code := http.StatusOK
	data := templates.NewMain("StajTakip - Mülakat Listesi")

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Warn("Mülakat formu okunamadı!")
			http.Error(w, "Hatalı istek!", http.StatusBadRequest)
			return
		}

		gorev, err := formStr(r.PostFormValue("gorev"))
		if err != nil {
			http.Error(w, "Hatalı istek!", http.StatusBadRequest)
			return
		}

		if gorev == "yenile" {
			code, data = sh.Yenile(w, r, data)
		} else if gorev == "guncelle" {
			code, data = sh.Guncelle(w, r, data)
		}
	}

	vars := MulakatListesiVars{nil, nil}
	data.Vars = vars

	if mulakatlar, err := database.MulakatListesi(sh.Conn); err == nil {
		vars.Mulakatlar = mulakatlar
	} else {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Mülakat listesi yüklenemedi!")
		w.WriteHeader(http.StatusInternalServerError)
		sablonHatasi(w, tpl_mulakat_listesi.ExecuteTemplate(w, "main", data.Error("Mülakat listesi yüklenemedi!")))
		return
	}

	if uyeler, err := database.KomisyonListesi(sh.Conn); err == nil {
		vars.Komisyon = uyeler
	} else {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Komisyon listesi yüklenemedi!")
		w.WriteHeader(http.StatusInternalServerError)
		sablonHatasi(w, tpl_mulakat_listesi.ExecuteTemplate(w, "main", data.Error("Komisyon listesi yüklenemedi!")))
		return
	}

	data.Vars = vars

	w.WriteHeader(code)
	sablonHatasi(w, tpl_mulakat_listesi.ExecuteTemplate(w, "main", data))
}

func (sh MulakatListesi) Yenile(w http.ResponseWriter, r *http.Request, data templates.Main) (int, templates.Main) {
	if err := database.MulakatListesiOlustur(sh.Conn); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Mülakat listesi yenilenemedi!")
		return http.StatusInternalServerError, data.Error("Mülakat listesi yenilenemedi!")
	}

	return http.StatusOK, data.Info("Mülakat listesi yenilendi!")
}

func (sh MulakatListesi) Guncelle(w http.ResponseWriter, r *http.Request, data templates.Main) (int, templates.Main) {
	no, err := formSayi(r.PostFormValue("no"))
	if err != nil {
		return http.StatusBadRequest, data.Warning("Hatalı istek!")
	}

	baslangic, err := formStr(r.PostFormValue("baslangic"))
	if err != nil {
		return http.StatusBadRequest, data.Warning("Hatalı istek!")
	}

	tarih, err := formStr(r.PostFormValue("tarih"))
	if err != nil {
		return http.StatusBadRequest, data.Warning("Tarih eksik veya yanlış!")
	}

	saat, err := formStr(r.PostFormValue("saat"))
	if err != nil {
		return http.StatusBadRequest, data.Warning("Saat eksik veya yanlış!")
	}

	komisyon1, err := formStr(r.PostFormValue("komisyon1"))
	if err != nil {
		return http.StatusBadRequest, data.Warning("Komisyon üyesi (1) eksik veya yanlış!")
	}

	komisyon2, err := formStr(r.PostFormValue("komisyon2"))
	if err != nil {
		return http.StatusBadRequest, data.Warning("Komisyon üyesi (2) eksik veya yanlış!")
	}

	tarihSaat, err := time.Parse(database.TarihSaatFormati, tarih + " " + saat)
	if err != nil {
		return http.StatusBadRequest, data.Warning("Tarih/saat formatı uygun değil!")
	}

	mul := database.Mulakat{no, baslangic, tarihSaat.Format(database.TarihSaatFormati), komisyon1, komisyon2}
	if err := mul.Update(sh.Conn); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Mülakat güncellenemedi!")
		return http.StatusInternalServerError, data.Error("Mülakat güncellenemedi!")
	}

	return http.StatusOK, data.Info("Mülakat listesi güncellendi!")
}
