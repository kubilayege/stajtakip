package routes

import (
	"net/http"
  "strconv"
  "time"

	"stajtakip/database"

	"github.com/sirupsen/logrus"
  "github.com/signintech/gopdf"
)

type SonucListesiPDF struct {
	Conn *database.Connection
}

func (sh SonucListesiPDF) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Geçersiz metod!", http.StatusNotFound)
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
      W: 110,
      H: 16,
    },
    &gopdf.Rect{ // Soyad
      W: 110,
      H: 16,
    },
    &gopdf.Rect{ // Öğretim
      W: 25,
      H: 16,
    },
    &gopdf.Rect{ // Başlangıç Tarihi
      W: 100,
      H: 16,
    },
    &gopdf.Rect{ // Toplam Gün
      W: 75,
      H: 16,
    },
    &gopdf.Rect{ // Kabul Edilen Gün
      W: 100,
      H: 16,
    },
    &gopdf.Rect{ // Boşluk
      W: 30,
      H: 16,
    },
  }

  var ogretimstr string
  for i, snc := range sonuclar {
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

      basliklar := []string{"No", "Ad", "Soyad", "I/II", "Başlangıç Tarihi", "Toplam Gün", "Kabul Edilen Gün"}
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

    ogretim := snc.Ogretim
    if ogretim == 1 {
      ogretimstr = "I"
    } else {
      ogretimstr = "II"
    }

    vals := []string{
      strconv.Itoa(snc.OgrenciNo),
      snc.Ad,
      snc.Soyad,
      ogretimstr,
      snc.StajBaslangic,
      strconv.Itoa(snc.ToplamGun),
      strconv.Itoa(snc.KabulGun),
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
  w.Header().Set("Content-Disposition", `attachment; filename="mulakat-sonuc.pdf"`)

  // Çıktı
  if err := pdf.Write(w); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("PDF oluşturulamadı!")
    http.Error(w, "PDF oluşturulamadı!", http.StatusInternalServerError)
    return
  }
}
