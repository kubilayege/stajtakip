package routes

import (
	"net/http"
	"strconv"
  "path/filepath"
  "io"
  "os"

	"stajtakip/database"
	"stajtakip/templates"
  "stajtakip/cfg"

	"github.com/sirupsen/logrus"
)

var tpl_ogrenci_belge = templates.Load("templates/ogrenci-belge.html")

type OgrenciBelge struct {
	Conn *database.Connection
}

// Verilen parametrelere göre veritabanına bir öğrenci eklemeye çalışır
func (sh OgrenciBelge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := templates.NewMain("StajTakip - Öğrenci Belge Ekle")

	if r.Method == http.MethodGet {
		sablonHatasi(w, tpl_ogrenci_belge.ExecuteTemplate(w, "main", data))
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Post metodu kullanılmalı!", http.StatusNotFound)
		return
	}

	if err := r.ParseMultipartForm(cfg.MaxIstekBoyutu()); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Warn("Öğrenci belge ekleme formu okunamadı!")
		w.WriteHeader(http.StatusBadRequest)
		sablonHatasi(w, tpl_ogrenci_belge.ExecuteTemplate(w, "main", data.Error("Öğrenci belge formunda bir hata var!")))
		return
	}

	var no int
	var err error

	no, err = formSayi(r.PostFormValue("no"))
	if err != nil || (no < 0) {
		w.WriteHeader(http.StatusBadRequest)
		sablonHatasi(w, tpl_ogrenci_belge.ExecuteTemplate(w, "main", data.Warning("Öğrenci no eksik veya yanlış!")))
		return
	}

  istekdosya, handler, err := r.FormFile("dosya")
  if err != nil {
  		logrus.WithFields(logrus.Fields{
  			"err": err,
  		}).Warn("Dosya okunurken bir hata oluştu!")
      w.WriteHeader(http.StatusInternalServerError)
  		sablonHatasi(w, tpl_ogrenci_belge.ExecuteTemplate(w, "main", data.Warning("Dosya ile ilgili bir problem oluştu!")))
     return
  }
  defer istekdosya.Close()

  dosyaadi := strconv.Itoa(no) + "-" + filepath.Base(handler.Filename)
  belgeyolu := filepath.Join("./uploads/", dosyaadi) // TODO upload yolunu cfg dosyasına ekle?

	ogr := database.OgrenciEk{no, dosyaadi}
	if err := ogr.Insert(sh.Conn); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sablonHatasi(w, tpl_ogrenci_belge.ExecuteTemplate(w, "main", data.Error(err.Error())))
		return
	}

  if filepath.Ext(dosyaadi) != ".pdf" {
      w.WriteHeader(http.StatusBadRequest)
  		sablonHatasi(w, tpl_ogrenci_belge.ExecuteTemplate(w, "main", data.Warning("Sadece pdf dosyası yüklenebilir!")))
     return
  }

  diskdosya, err := os.OpenFile(belgeyolu, os.O_WRONLY|os.O_CREATE, 0666)
  if err != nil {
  		logrus.WithFields(logrus.Fields{
  			"err": err,
  		}).Warn("Dosya yazılırken bir hata oluştu!")
      w.WriteHeader(http.StatusInternalServerError)
  		sablonHatasi(w, tpl_ogrenci_belge.ExecuteTemplate(w, "main", data.Warning("Dosya ile ilgili bir problem oluştu!")))
     return
  }
  defer diskdosya.Close()

  if _, err := io.Copy(diskdosya, istekdosya); err != nil {
  		logrus.WithFields(logrus.Fields{
  			"err": err,
  		}).Warn("Dosya kopyalanırken bir hata oluştu!")
      w.WriteHeader(http.StatusInternalServerError)
  		sablonHatasi(w, tpl_ogrenci_belge.ExecuteTemplate(w, "main", data.Warning("Dosya ile ilgili bir problem oluştu!")))
     return
  }

	w.WriteHeader(http.StatusOK)
	sablonHatasi(w, tpl_ogrenci_belge.ExecuteTemplate(w, "main", data.Info("Öğrenci belgesi veritabanına başarıyla eklendi!")))
}
