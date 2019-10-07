package routes

import (
	"net/http"
  "strconv"

	"stajtakip/database"

	"github.com/sirupsen/logrus"
  "github.com/signintech/gopdf"
)

const OgrenciAdiUzunlugu float64 = 80
const KomisyonAdiUzunlugu float64 = 120

type MulakatListesiPDF struct {
	Conn *database.Connection
}

func (sh MulakatListesiPDF) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Geçersiz metod!", http.StatusNotFound)
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
    &gopdf.Rect{ // No
      W: 50,
      H: 16,
    },
    &gopdf.Rect{ // Ad
      W: OgrenciAdiUzunlugu,
      H: 16,
    },
    &gopdf.Rect{ // Soyad
      W: OgrenciAdiUzunlugu,
      H: 16,
    },
    &gopdf.Rect{ // Öğretim
      W: 25,
      H: 16,
    },
    &gopdf.Rect{ // Tarih
      W: 60,
      H: 16,
    },
    &gopdf.Rect{ // Saat
      W: 35,
      H: 16,
    },
    &gopdf.Rect{ // Komisyon
      W: KomisyonAdiUzunlugu,
      H: 16,
    },
    &gopdf.Rect{ // Komisyon
      W: KomisyonAdiUzunlugu,
      H: 16,
    },
    &gopdf.Rect{ // Boşluk
      W: 30,
      H: 16,
    },
  }

  var ogretimstr string
  for i, mul := range mulakatlar {
    if i % 40 == 0 {
      pdf.AddPage()

      // Liste başlığı
      header_opt.Border = 1
      if err := pdf.SetFont("Open Sans Bold", "", 12); err != nil {
        logrus.WithFields(logrus.Fields{
          "err": err,
        }).Error("Fontlar yüklenemedi!")
        http.Error(w, "PDF oluşturulamadı!", http.StatusInternalServerError)
        return
      }

      basliklar := []string{"No", "Ad", "Soyad", "I/II", "Tarih", "Saat", "Komisyon", "Komisyon"}
      pdf.SetX(10)
      for i, baslik := range basliklar {
        if pdf.CellWithOption(rect[i], baslik, header_opt) != nil {
          logrus.WithFields(logrus.Fields{
          	"err": err,
          }).Error("Header oluşturulamadı!")
          http.Error(w, "PDF oluşturulamadı!", http.StatusInternalServerError)
          return
        }
      }
      pdf.Br(20)

      // Liste değerleri
      header_opt.Border = 0
      if err := pdf.SetFont("Open Sans", "", 10); err != nil {
        logrus.WithFields(logrus.Fields{
          "err": err,
        }).Error("Fontlar yüklenemedi!")
        http.Error(w, "PDF oluşturulamadı!", http.StatusInternalServerError)
        return
      }
    }

    ogretim := mul.Ogretim
    if ogretim == 1 {
      ogretimstr = "I"
    } else {
      ogretimstr = "II"
    }

    vals := []string{
      strconv.Itoa(mul.OgrenciNo),
      mul.Ad,
      mul.Soyad,
      ogretimstr,
      mul.Tarih,
      mul.Saat,
      mul.KomisyonUye1,
      mul.KomisyonUye2,
    }

    pdf.SetX(10)
    for i, val := range vals {
      // Değeri yazar
      if pdf.CellWithOption(rect[i], val, header_opt) != nil {
        logrus.WithFields(logrus.Fields{
          "err": err,
        }).Error("Liste oluşturulamadı!")
        http.Error(w, "PDF oluşturulamadı!", http.StatusInternalServerError)
        return
      }

      // Taşan kısımların üstünü örter
      pdf.SetFillColor(255, 255, 255)
      pdf.RectFromUpperLeftWithStyle(pdf.GetX(), pdf.GetY(), rect[i+1].W, rect[i+1].H, "F")
      pdf.SetFillColor(0, 0, 0)
    }
    pdf.Br(20)
  }

  w.Header().Set("Content-Type", "application/pdf")
  w.Header().Set("Content-Disposition", `attachment; filename="mulakat-listesi.pdf"`)

  // Çıktı
  if err := pdf.Write(w); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("PDF oluşturulamadı!")
    http.Error(w, "PDF oluşturulamadı!", http.StatusInternalServerError)
    return
  }
}
