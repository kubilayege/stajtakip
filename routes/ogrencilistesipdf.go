package routes

import (
	"net/http"
  "strconv"

	"stajtakip/database"

	"github.com/sirupsen/logrus"
  "github.com/signintech/gopdf"
)

type OgrenciListesiPDF struct {
	Conn *database.Connection
}

func (sh OgrenciListesiPDF) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

  pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{ PageSize: gopdf.PageSizeA4 })

  if err := pdf.AddTTFFont("Open Sans", "./fonts/OpenSans-Regular.ttf"); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Fontlar yüklenemedi!")
    http.Error(w, "PDF oluşturulamadı!", http.StatusInternalServerError)
    return
  }

  if err := pdf.AddTTFFont("Open Sans Bold", "./fonts/OpenSans-Bold.ttf"); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Fontlar yüklenemedi!")
    http.Error(w, "PDF oluşturulamadı!", http.StatusInternalServerError)
    return
  }

  header_opt :=  gopdf.CellOption{
    Align: 0,
    Float: 0,
    Border: 1,
  }
  rect := []*gopdf.Rect{
    &gopdf.Rect{
      W: 60,
      H: 16,
    },
    &gopdf.Rect{
      W: 200,
      H: 16,
    },
    &gopdf.Rect{
      W: 200,
      H: 16,
    },
    &gopdf.Rect{
      W: 60,
      H: 16,
    },
  }

  var ogretimstr string
  for i, ogr := range ogrenciler {
    if i % 40 == 0 {
      pdf.AddPage()

      // Liste başlığı
      header_opt.Border = 1
      if err := pdf.SetFont("Open Sans Bold", "", 14); err != nil {
        logrus.WithFields(logrus.Fields{
          "err": err,
        }).Error("Fontlar yüklenemedi!")
        http.Error(w, "PDF oluşturulamadı!", http.StatusInternalServerError)
        return
      }

      pdf.SetX(20)
      if pdf.CellWithOption(rect[0], "No", header_opt) != nil || pdf.CellWithOption(rect[1], "Ad", header_opt) != nil || pdf.CellWithOption(rect[2], "Soyad", header_opt) != nil || pdf.CellWithOption(rect[3], "Öğretim", header_opt) != nil {
        logrus.WithFields(logrus.Fields{
        	"err": err,
        }).Error("Header oluşturulamadı!")
        http.Error(w, "PDF oluşturulamadı!", http.StatusInternalServerError)
        return
      }

      pdf.Br(20)

      // Liste değerleri
      header_opt.Border = 0
      if err := pdf.SetFont("Open Sans", "", 12); err != nil {
        logrus.WithFields(logrus.Fields{
          "err": err,
        }).Error("Fontlar yüklenemedi!")
        http.Error(w, "PDF oluşturulamadı!", http.StatusInternalServerError)
        return
      }
    }

    ogretim := ogr.Ogretim
    if ogretim == 1 {
      ogretimstr = "I"
    } else {
      ogretimstr = "II"
    }

    pdf.SetX(20)
    if pdf.CellWithOption(rect[0], strconv.Itoa(ogr.No), header_opt) != nil || pdf.CellWithOption(rect[1], ogr.Ad, header_opt) != nil || pdf.CellWithOption(rect[2], ogr.Soyad, header_opt) != nil || pdf.CellWithOption(rect[3], ogretimstr, header_opt) != nil {
      logrus.WithFields(logrus.Fields{
        "err": err,
      }).Error("Liste oluşturulamadı!")
      http.Error(w, "PDF oluşturulamadı!", http.StatusInternalServerError)
      return
    }
    pdf.Br(20)
  }

  // Çıktı
  if err := pdf.Write(w); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("PDF oluşturulamadı!")
    http.Error(w, "PDF oluşturulamadı!", http.StatusInternalServerError)
    return
  }
}
