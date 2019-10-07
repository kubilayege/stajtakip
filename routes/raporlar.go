package routes

import (
	"net/http"
	"stajtakip/database"
	"stajtakip/templates"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

var tpl_404 = templates.Load("templates/404.html")
var tpl_raporlar = templates.Load("templates/raporlar.html")

type Raporlar struct {
	Conn *database.Connection
}

type RaporlarVars struct {
	KonuRapor  []database.RaporKonu
	DagilimRapor  []database.RaporKonuDagilim
	SehirRapor []database.RaporSehir
	Baslangic int
	Bitis int
}

func (sh Raporlar) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)

		err := tpl_404.ExecuteTemplate(w, "main", templates.NewMain("StajTakip - Sayfa Bulunamadı"))
		if err != nil {
			http.Error(w, "Sayfa yüklenemedi!", http.StatusInternalServerError)
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Warn("Şablon çalıştırılamadı!")
		}
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Geçersiz metod!", http.StatusNotFound)
		return
	}

	data := templates.NewMain("StajTakip - Raporlar")
	vars := RaporlarVars{}
	data.Vars = vars

	var baslangic, bitis int
	var err error

	yearstr := r.FormValue("baslangic")
	if yearstr == "" {
		baslangic = time.Now().Year() - 1
	} else {
		baslangic, err = strconv.Atoi(yearstr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			sablonHatasi(w, tpl_raporlar.ExecuteTemplate(w, "main", data.Warning("Yıl eksik!")))
			return
		}
	}

	yearstr = r.FormValue("bitis")
	if yearstr == "" {
		bitis = time.Now().Year() - 1
	} else {
		bitis, err = strconv.Atoi(yearstr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			sablonHatasi(w, tpl_raporlar.ExecuteTemplate(w, "main", data.Warning("Yıl eksik!")))
			return
		}
	}

	if baslangic > bitis {
		w.WriteHeader(http.StatusBadRequest)
		sablonHatasi(w, tpl_raporlar.ExecuteTemplate(w, "main", data.Warning("Yıl yanlış!")))
		return
	}

	vars.Baslangic = baslangic
	vars.Bitis = bitis

	if sehirRapor, err := database.RaporSehirler(sh.Conn, baslangic, bitis); err == nil {
		vars.SehirRapor = sehirRapor
	} else {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Rapor listesi oluşturulamadı!")
		w.WriteHeader(http.StatusInternalServerError)
		sablonHatasi(w, tpl_raporlar.ExecuteTemplate(w, "main", data.Error("Rapor listesi oluşturulamadı!")))
		return
	}

	if konuRapor, err := database.RaporKonular(sh.Conn, baslangic, bitis); err == nil {
		vars.KonuRapor = konuRapor
	} else {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Rapor listesi oluşturulamadı!")
		w.WriteHeader(http.StatusInternalServerError)
		sablonHatasi(w, tpl_raporlar.ExecuteTemplate(w, "main", data.Error("Rapor listesi oluşturulamadı!")))
		return
	}

	if dagilimRapor, err := database.RaporKonularDagilim(sh.Conn, baslangic, bitis); err == nil {
		vars.DagilimRapor = dagilimRapor
	} else {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Rapor listesi oluşturulamadı!")
		w.WriteHeader(http.StatusInternalServerError)
		sablonHatasi(w, tpl_raporlar.ExecuteTemplate(w, "main", data.Error("Rapor listesi oluşturulamadı!")))
		return
	}

	data.Vars = vars
	sablonHatasi(w, tpl_raporlar.ExecuteTemplate(w, "main", data))
}
