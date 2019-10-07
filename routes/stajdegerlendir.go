package routes

import (
  "net/http"
  "stajtakip/database"
  "stajtakip/templates"

  "github.com/sirupsen/logrus"
)

var tpl_staj_degerlendir = templates.Load("templates/staj-degerlendir.html")

type StajDegerlendirVars struct {
  Mul *database.MulakatOgrenci
}

type StajDegerlendir struct {
	Conn *database.Connection
}

func (sh StajDegerlendir) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost && r.Method != http.MethodGet {
  	http.Error(w, "Geçersiz metod!", http.StatusNotFound)
  	return
  }

  data := templates.NewMain("StajTakip - Staj Değerlendir")

  if err := r.ParseForm(); err != nil {
    logrus.WithFields(logrus.Fields{
      "err": err,
    }).Warn("Sayfa parametreleri okunamadı!")
    w.WriteHeader(http.StatusBadRequest)
    sablonHatasi(w, tpl_staj_degerlendir.ExecuteTemplate(w, "main", data.Error("Parametre hatası!")))
    return
  }

  var no int
  var baslangic string
  var err error
  var code int = http.StatusOK

  vars := StajDegerlendirVars{
    Mul: nil,
  }

  if r.Method == http.MethodPost {
    code, data = sh.Post(w, r, data)
  }

  nostr := r.FormValue("no")
  baslangicstr := r.FormValue("baslangic")

  if nostr == "" || baslangicstr == "" {
    data.Vars = vars
  	w.WriteHeader(code)
  	sablonHatasi(w, tpl_staj_degerlendir.ExecuteTemplate(w, "main", data))
    return
  }

	no, err = formSayi(nostr)
	if err != nil || (no < 0) {
		w.WriteHeader(http.StatusBadRequest)
		sablonHatasi(w, tpl_staj_degerlendir.ExecuteTemplate(w, "main", data.Warning("Öğrenci no eksik veya yanlış!")))
		return
	}

  baslangic, err = formStr(baslangicstr)
  if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		sablonHatasi(w, tpl_staj_degerlendir.ExecuteTemplate(w, "main", data.Warning("Başlangıç tarihi eksik veya yanlış!")))
    return
  }

  vars.Mul, err = database.MulakatOgrenciBul(sh.Conn, no, baslangic)
  data.Vars = vars
  if err == database.ErrVeriBulunamadi {
    w.WriteHeader(http.StatusBadRequest)
    sablonHatasi(w, tpl_staj_degerlendir.ExecuteTemplate(w, "main", data.Error("Mülakat bulunamadı!")))
    return
  } else if err != nil {
    logrus.WithFields(logrus.Fields{
      "err": err,
    }).Error("Mülakat aranırken veritabanında bir hata oluştu!")
    w.WriteHeader(http.StatusInternalServerError)
    sablonHatasi(w, tpl_staj_degerlendir.ExecuteTemplate(w, "main", data.Error("Veritabanında bir hata oluştu!")))
    return
  }

  if vars.Mul.Tarih == "" {
    data = data.Warning("Bu mülakat için tarih bulunmuyor!")
  }

	w.WriteHeader(code)
	sablonHatasi(w, tpl_staj_degerlendir.ExecuteTemplate(w, "main", data))
}

func (sh StajDegerlendir) Post(w http.ResponseWriter, r *http.Request, data templates.Main) (int, templates.Main) {
  var no int
  var baslangic string
  var err error

  no, err = formSayi(r.PostFormValue("no"))
  if err != nil || (no < 0) {
    return http.StatusBadRequest, data.Warning("Öğrenci no eksik veya yanlış!")
  }

  baslangic, err = formStr(r.PostFormValue("baslangic"))
  if err != nil {
    return http.StatusBadRequest, data.Warning("Başlangıç tarihi eksik veya yanlış!")
  }

  var PuanDevam, PuanCaba, PuanVakit, PuanAmireDavranis, PuanIsArkadasaDavranis,
      PuanProje, PuanDuzen, PuanSunum, PuanIcerik, PuanMulakat int

  PuanDevam, err = formSayi(r.PostFormValue("PuanDevam"))
  if err != nil {
    return http.StatusBadRequest, data.Warning("Puan eksik veya yanlış!")
  }

  PuanCaba, err = formSayi(r.PostFormValue("PuanCaba"))
  if err != nil {
    return http.StatusBadRequest, data.Warning("Puan eksik veya yanlış!")
  }

  PuanVakit, err = formSayi(r.PostFormValue("PuanVakit"))
  if err != nil {
    return http.StatusBadRequest, data.Warning("Puan eksik veya yanlış!")
  }

  PuanAmireDavranis, err = formSayi(r.PostFormValue("PuanAmireDavranis"))
  if err != nil {
    return http.StatusBadRequest, data.Warning("Puan eksik veya yanlış!")
  }

  PuanIsArkadasaDavranis, err = formSayi(r.PostFormValue("PuanIsArkadasaDavranis"))
  if err != nil {
    return http.StatusBadRequest, data.Warning("Puan eksik veya yanlış!")
  }

  PuanProje, err = formSayi(r.PostFormValue("PuanProje"))
  if err != nil {
    return http.StatusBadRequest, data.Warning("Puan eksik veya yanlış!")
  }

  PuanDuzen, err = formSayi(r.PostFormValue("PuanDuzen"))
  if err != nil {
    return http.StatusBadRequest, data.Warning("Puan eksik veya yanlış!")
  }

  PuanSunum, err = formSayi(r.PostFormValue("PuanSunum"))
  if err != nil {
    return http.StatusBadRequest, data.Warning("Puan eksik veya yanlış!")
  }

  PuanIcerik, err = formSayi(r.PostFormValue("PuanIcerik"))
  if err != nil {
    return http.StatusBadRequest, data.Warning("Puan eksik veya yanlış!")
  }

  PuanMulakat, err = formSayi(r.PostFormValue("PuanMulakat"))
  if err != nil {
    return http.StatusBadRequest, data.Warning("Puan eksik veya yanlış!")
  }

  mul, err := database.MulakatOgrenciBul(sh.Conn, no, baslangic)
  if err == database.ErrVeriBulunamadi {
    return http.StatusBadRequest, data.Error("Mülakat bulunamadı!")
  } else if err != nil {
    logrus.WithFields(logrus.Fields{
      "err": err,
    }).Error("Mülakat aranırken veritabanında bir hata oluştu!")
    return http.StatusInternalServerError, data.Error("Mülakat sonucu eklenirken bir hata oluştu!")
  }

  if mul.Tarih == "" {
    return http.StatusBadRequest, data.Warning("Bu mülakat için tarih girilmemiş!")
  }

  sonuc := database.MulakatSonuc{
    no, baslangic, PuanDevam, PuanCaba, PuanVakit,
    PuanAmireDavranis, PuanIsArkadasaDavranis,
    PuanProje, PuanDuzen, PuanSunum, PuanIcerik, PuanMulakat,
  }

  if err := sonuc.Update(sh.Conn); err != nil {
    logrus.WithFields(logrus.Fields{
      "err": err,
    }).Warn("Mülakat sonucu eklenirken bir hata oluştu!")
    return http.StatusInternalServerError, data.Error("Mülakat sonucu eklenirken bir hata oluştu!")
  }

  return http.StatusOK, data.Info("Mülakat sonucu eklendi!")
}
