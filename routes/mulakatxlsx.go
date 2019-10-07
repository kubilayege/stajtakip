package routes

import (
	"net/http"
  "strconv"

	"stajtakip/database"

	"github.com/sirupsen/logrus"
  "github.com/360EntSecGroup-Skylar/excelize"
)

type MulakatListesiXLSX struct {
	Conn *database.Connection
}

func (sh MulakatListesiXLSX) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Get metodu kullanılmalı!", http.StatusNotFound)
		return
	}

  mulakatlar, err := database.MulakatListesi(sh.Conn)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Mülakat listesi yüklenemedi!")
    http.Error(w, "Mülakat listesi yüklenemedi!", http.StatusInternalServerError)
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
  xlsx.SetCellValue(sheetName, "E1", "Mülakat Tarihi")
  xlsx.SetCellValue(sheetName, "F1", "Mülakat Saati")
  xlsx.SetCellValue(sheetName, "G1", "Komisyon Üyesi")
  xlsx.SetCellValue(sheetName, "H1", "Komisyon Üyesi")

  var ogretimstr string
  for i, mul := range mulakatlar {
    ogretim := mul.Ogretim
    if ogretim == 1 {
      ogretimstr = "I"
    } else {
      ogretimstr = "II"
    }

    xlsx.SetCellValue(sheetName, "A" + strconv.Itoa(i+2), mul.OgrenciNo)
    xlsx.SetCellValue(sheetName, "B" + strconv.Itoa(i+2), mul.Ad)
    xlsx.SetCellValue(sheetName, "C" + strconv.Itoa(i+2), mul.Soyad)
    xlsx.SetCellValue(sheetName, "D" + strconv.Itoa(i+2), ogretimstr)
    xlsx.SetCellValue(sheetName, "E" + strconv.Itoa(i+2), mul.Tarih)
    xlsx.SetCellValue(sheetName, "F" + strconv.Itoa(i+2), mul.Saat)
    xlsx.SetCellValue(sheetName, "G" + strconv.Itoa(i+2), mul.KomisyonUye1)
    xlsx.SetCellValue(sheetName, "H" + strconv.Itoa(i+2), mul.KomisyonUye2)
  }

  w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
  w.Header().Set("Content-Disposition", `attachment; filename="mulakat-listesi.xlsx"`)

  // Çıktı
  if err := xlsx.Write(w); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("XLSX oluşturulamadı!")
    http.Error(w, "XLSX oluşturulamadı!", http.StatusInternalServerError)
    return
  }
}
