package routes

import (
	"net/http"
	"time"

	"stajtakip/database"
	"stajtakip/templates"

	"github.com/sirupsen/logrus"
)

var tpl_staj_ekle = templates.Load("templates/staj-ekle.html")

const zamanFormati = "2006-01-02"

type StajEkleVars struct {
	Konular  []database.Konu
	Kurumlar []string
	Stajlar []database.TumStaj
	DenkStaj bool
	No int
}

type StajEkle struct {
	Conn *database.Connection
}

func (sh StajEkle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Geçersiz metod!", http.StatusNotFound)
		return
	}

	data := templates.NewMain("StajTakip - Öğrenci Ekle")
	vars := StajEkleVars{}
	data.Vars = vars

	if err := r.ParseForm(); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Warn("Staj ekleme formu okunamadı!")
		w.WriteHeader(http.StatusBadRequest)
		sablonHatasi(w, tpl_staj_ekle.ExecuteTemplate(w, "main", data.Error("Parametre hatası!")))
		return
	}

	var no int
	var err error

	nostr := r.FormValue("no")
	if nostr == "" {
		sablonHatasi(w, tpl_staj_ekle.ExecuteTemplate(w, "main", data))
		return
	} else {
		no, err = formSayi(nostr)
		if err != nil || (no < 1) {
			w.WriteHeader(http.StatusBadRequest)
			sablonHatasi(w, tpl_staj_ekle.ExecuteTemplate(w, "main", data.Warning("Öğrenci no eksik veya yanlış!")))
			return
		}
	}

	vars.No = no

	if konular, err := database.KonuListesi(sh.Conn); err == nil {
		vars.Konular = konular
	} else {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Konu listesi yüklenemedi!")
		w.WriteHeader(http.StatusInternalServerError)
		sablonHatasi(w, tpl_staj_ekle.ExecuteTemplate(w, "main", data.Error("Konu listesi yüklenemedi!")))
		return
	}

	if kurumlar, err := database.KurumListesi(sh.Conn); err == nil {
		vars.Kurumlar = kurumlar
	} else {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Kurum listesi yüklenemedi!")
		w.WriteHeader(http.StatusInternalServerError)
		sablonHatasi(w, tpl_staj_ekle.ExecuteTemplate(w, "main", data.Error("Kurum listesi yüklenemedi!")))
		return
	}

	if stjlar, err := database.OgrenciTumStajListesi(sh.Conn, no); err == nil {
		vars.Stajlar = stjlar
	} else {
    logrus.WithFields(logrus.Fields{
      "err": err,
    }).Error("Öğrenci staj bilgileri aranırken veritabanında bir hata oluştu!")
    w.WriteHeader(http.StatusInternalServerError)
    sablonHatasi(w, tpl_staj_ekle.ExecuteTemplate(w, "main", data.Error("Veritabanında bir hata oluştu!")))
    return
  }

	data.Vars = vars

	if r.Method == http.MethodGet {
		sablonHatasi(w, tpl_staj_ekle.ExecuteTemplate(w, "main", data))
		return
	}

	denk, err := formSayi(r.PostFormValue("denk"))
	if err != nil {
		http.Error(w, "Eksik bilgi!", http.StatusBadRequest)
		return
	}

	if denk == 0 {
		sh.NormalStajEkle(data, w, r, no)
	} else {
		sh.DenkStajEkle(data, w, r, no)
	}
}

func (sh StajEkle) NormalStajEkle(data templates.Main, w http.ResponseWriter, r *http.Request, ogrno int) {
	var sinif, toplamgun int
	var kurum, sehir, konu, baslangic, bitis string
	var err error

	kurum, err = formStr(r.PostFormValue("kurum"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		sablonHatasi(w, tpl_staj_ekle.ExecuteTemplate(w, "main", data.Warning("Kurum eksik veya yanlış!")))
		return
	}

	sehir, err = formStr(r.PostFormValue("sehir"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		sablonHatasi(w, tpl_staj_ekle.ExecuteTemplate(w, "main", data.Warning("Şehir eksik veya yanlış!")))
		return
	}

	konu, err = formStr(r.PostFormValue("konu"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		sablonHatasi(w, tpl_staj_ekle.ExecuteTemplate(w, "main", data.Warning("Konu eksik veya yanlış!")))
		return
	}

	var t1, t2 time.Time

	baslangic, err = formStr(r.PostFormValue("baslangic"))
	if err == nil {
		t1, err = time.Parse(zamanFormati, baslangic)
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		sablonHatasi(w, tpl_staj_ekle.ExecuteTemplate(w, "main", data.Warning("Başlangıç tarihi eksik veya yanlış!")))
		return
	}

	bitis, err = formStr(r.PostFormValue("bitis"))
	if err == nil {
		t2, err = time.Parse(zamanFormati, bitis)
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		sablonHatasi(w, tpl_staj_ekle.ExecuteTemplate(w, "main", data.Warning("Bitiş tarihi eksik veya yanlış!")))
		return
	}

	if t2.Sub(t1).Hours() < 1 {
		w.WriteHeader(http.StatusBadRequest)
		sablonHatasi(w, tpl_staj_ekle.ExecuteTemplate(w, "main", data.Warning("Başlangıç tarihi, bitiş tarihinden önce olmalıdır!")))
		return
	}

	sinif, err = formSayi(r.PostFormValue("sinif"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		sablonHatasi(w, tpl_staj_ekle.ExecuteTemplate(w, "main", data.Warning("Sınıf eksik veya yanlış!")))
		return
	}

	toplamgun, err = formSayi(r.PostFormValue("toplamgun"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		sablonHatasi(w, tpl_staj_ekle.ExecuteTemplate(w, "main", data.Warning("Toplam gün eksik veya yanlış!")))
		return
	}

	ogr := database.Staj{ogrno, kurum, sehir, konu, baslangic, bitis, sinif, toplamgun, 0, false}
	if err := ogr.Insert(sh.Conn); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sablonHatasi(w, tpl_staj_ekle.ExecuteTemplate(w, "main", data.Error(err.Error())))
		return
	}

	w.WriteHeader(http.StatusOK)
	sablonHatasi(w, tpl_staj_ekle.ExecuteTemplate(w, "main", data.Info("Staj bilgisi veritabanına başarıyla eklendi!")))

	// TODO Eski değerleri inputlara ata
}

func (sh StajEkle) DenkStajEkle(data templates.Main, w http.ResponseWriter, r *http.Request, ogrno int) {
	var toplamgun int
	var kurum, okul string
	var err error

	var dat StajEkleVars = data.Vars.(StajEkleVars)
	dat.DenkStaj = true
	data.Vars = dat

	kurum, err = formStr(r.PostFormValue("kurum"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		sablonHatasi(w, tpl_staj_ekle.ExecuteTemplate(w, "main", data.Warning("Kurum eksik veya yanlış!")))
		return
	}

	okul, err = formStr(r.PostFormValue("okul"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		sablonHatasi(w, tpl_staj_ekle.ExecuteTemplate(w, "main", data.Warning("Önceki okul eksik!")))
		return
	}

	toplamgun, err = formSayi(r.PostFormValue("toplamgun"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		sablonHatasi(w, tpl_staj_ekle.ExecuteTemplate(w, "main", data.Warning("Toplam gün eksik veya yanlış!")))
		return
	}

	ogr := database.DenkStaj{ogrno, kurum, okul, toplamgun, toplamgun / 2}
	if err := ogr.Insert(sh.Conn); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sablonHatasi(w, tpl_staj_ekle.ExecuteTemplate(w, "main", data.Error(err.Error())))
		return
	}

	w.WriteHeader(http.StatusOK)
	sablonHatasi(w, tpl_staj_ekle.ExecuteTemplate(w, "main", data.Info("DGS/Yatay Geçiş Staj bilgisi veritabanına başarıyla eklendi!")))

	// TODO Eski değerleri inputlara ata
}
