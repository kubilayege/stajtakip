package routes

import (
	"net/http"
  "strconv"
  "fmt"

	"stajtakip/database"

	"github.com/sirupsen/logrus"
  "github.com/360EntSecGroup-Skylar/excelize"
)

type OgrenciListesiXLS struct {
	Conn *database.Connection
}

func (sh OgrenciListesiXLS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Get metodu kullanılmalı!", http.StatusNotFound)
		return
	}

  ogrenciler, err := database.StajiTamamOgrenciler(sh.Conn)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Öğrenci listesi yüklenemedi!")
    http.Error(w, "Öğrenci listesi yüklenemedi!", http.StatusInternalServerError)
		return
	}

  xlsx := excelize.NewFile()
  sheetName := xlsx.GetSheetName(1)
  fmt.Println(sheetName)

  style, err := xlsx.NewStyle(`{"font":{"bold":true}}`)
  if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("XLSX oluşturulamadı!")
    http.Error(w, "XLSX oluşturulamadı!", http.StatusInternalServerError)
		return
  }
  xlsx.SetCellStyle(sheetName, "A1", "D1", style)

  xlsx.SetCellValue(sheetName, "A1", "No")
  xlsx.SetCellValue(sheetName, "B1", "Ad")
  xlsx.SetCellValue(sheetName, "C1", "Soyad")
  xlsx.SetCellValue(sheetName, "D1", "Öğretim")

  var ogretimstr string
  for i, ogr := range ogrenciler {
    ogretim := ogr.Ogretim
    if ogretim == 1 {
      ogretimstr = "I"
    } else {
      ogretimstr = "II"
    }

    xlsx.SetCellValue(sheetName, "A" + strconv.Itoa(i+2), ogr.No)
    xlsx.SetCellValue(sheetName, "B" + strconv.Itoa(i+2), ogr.Ad)
    xlsx.SetCellValue(sheetName, "C" + strconv.Itoa(i+2), ogr.Soyad)
    xlsx.SetCellValue(sheetName, "D" + strconv.Itoa(i+2), ogretimstr)
  }

  w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
  w.Header().Set("Content-Disposition", `attachment; filename="staji-bitenler.xlsx"`)

  // Çıktı
  if err := xlsx.Write(w); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("XLSX oluşturulamadı!")
    http.Error(w, "XLSX oluşturulamadı!", http.StatusInternalServerError)
    return
  }
}
