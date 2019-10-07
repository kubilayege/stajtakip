package routes

import (
	"net/http"
  "strconv"
  "time"

	"stajtakip/database"

	"github.com/sirupsen/logrus"
  "github.com/360EntSecGroup-Skylar/excelize"
)

type SonucListesiXLSX struct {
	Conn *database.Connection
}

func (sh SonucListesiXLSX) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Get metodu kullanılmalı!", http.StatusNotFound)
		return
	}

  if err := r.ParseForm(); err != nil {
    http.Error(w, "Parametreler okunamadı!", http.StatusBadRequest)
    return
  }

  baslangic, err := formStr( r.FormValue("baslangic"))
  if err != nil {
    http.Error(w, "Başlangıç tarihi eksik veya yanlış!", http.StatusBadRequest)
    return
  }

  bitis, err := formStr(r.FormValue("bitis"))
  if err != nil {
    http.Error(w, "Bitiş tarihi eksik veya yanlış!", http.StatusBadRequest)
    return
  }

  bas, err := time.Parse(database.TarihFormati, baslangic)
  if err != nil {
  	http.Error(w, "Hatalı tarih formatı!", http.StatusBadRequest)
    return
  }

  son, err := time.Parse(database.TarihFormati, bitis)
  if err != nil {
  	http.Error(w, "Hatalı tarih formatı!", http.StatusBadRequest)
    return
  }

  if bas.After(son) {
    http.Error(w, "Başlangıç tarihi, bitiş tarihinden sonra olamaz!", http.StatusBadRequest)
    return
  }

  sonuclar, err := database.MulakatSonucListesi(sh.Conn, baslangic, bitis)
  if err != nil {
    logrus.WithFields(logrus.Fields{
      "err": err,
    }).Error("Mülakat sonuçları listelenirken veritabanında bir hata oluştu!")
    http.Error(w, "Veritabanında bir hata oluştu!", http.StatusInternalServerError)
    return
  }

  xlsx := excelize.NewFile()
  sheetName := xlsx.GetSheetName(1)

  style, err := xlsx.NewStyle(`{"font":{"bold":true}}`)
  if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("XLSX oluşturulamadı!")
    http.Error(w, "XLSX oluşturulamadı!", http.StatusInternalServerError)
		return
  }
  xlsx.SetCellStyle(sheetName, "A1", "H1", style)

  xlsx.SetCellValue(sheetName, "A1", "No")
  xlsx.SetCellValue(sheetName, "B1", "Ad")
  xlsx.SetCellValue(sheetName, "C1", "Soyad")
  xlsx.SetCellValue(sheetName, "D1", "Öğretim")
  xlsx.SetCellValue(sheetName, "E1", "Başlangıç Tarihi")
  xlsx.SetCellValue(sheetName, "F1", "Toplam Gün")
  xlsx.SetCellValue(sheetName, "G1", "Kabul Edilen Gün")

  var ogretimstr string
  for i, snc := range sonuclar {
    ogretim := snc.Ogretim
    if ogretim == 1 {
      ogretimstr = "I"
    } else {
      ogretimstr = "II"
    }

    xlsx.SetCellValue(sheetName, "A" + strconv.Itoa(i+2), snc.OgrenciNo)
    xlsx.SetCellValue(sheetName, "B" + strconv.Itoa(i+2), snc.Ad)
    xlsx.SetCellValue(sheetName, "C" + strconv.Itoa(i+2), snc.Soyad)
    xlsx.SetCellValue(sheetName, "D" + strconv.Itoa(i+2), ogretimstr)
    xlsx.SetCellValue(sheetName, "E" + strconv.Itoa(i+2), snc.StajBaslangic)
    xlsx.SetCellValue(sheetName, "F" + strconv.Itoa(i+2), snc.ToplamGun)
    xlsx.SetCellValue(sheetName, "G" + strconv.Itoa(i+2), snc.KabulGun)
  }

  w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
  w.Header().Set("Content-Disposition", `attachment; filename="mulakat-sonuc.xlsx"`)

  // Çıktı
  if err := xlsx.Write(w); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("XLSX oluşturulamadı!")
    http.Error(w, "XLSX oluşturulamadı!", http.StatusInternalServerError)
    return
  }
}
