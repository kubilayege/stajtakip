package main

import (
	"log"
	"stajtakip/cfg"
	"stajtakip/database"
	"stajtakip/routes"

	"github.com/sirupsen/logrus"
)

const AYARLAR_DOSYASI string = "ayarlar.json"

func main() {
	// Ayarları oku
	err := cfg.AyarlariOku(AYARLAR_DOSYASI)
	if err != nil {
		log.Fatalln("Ayarlar dosyası okunamadı:", err)
	}

	// Logrus'u ayarla
	LogBaslat(cfg.LogDosyasi())

	// Max istek boyutunu ayarla
	routes.Ayarla(cfg.MaxIstekBoyutu())

	srvAddr := cfg.SunucuAdresi()
	srv := NewStajServer(srvAddr)

	dbAddr := cfg.VeritabaniAdresi()
	if dbAddr == "" {
		logrus.Error("Veritabanı adresi belirtilmemiş!")
		return
	}

	conn := database.NewConnection(dbAddr)
	dbOk := make(chan bool)
	errChan := make(chan error)

	go func() {
		logrus.WithFields(logrus.Fields{
			"database-address": dbAddr,
		}).Info("Veritabanına bağlanılıyor...")

		// Veritabanına bağlantıyı sağla
		if err := conn.Connect(dbOk); err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Error("Veritabanı bağlantısında bir hata meydana geldi!")
			errChan <- err
		}
	}()

	go func() {
		<-dbOk // Veritabanı bağlantısını bekle

		srv.SetHandlers(conn)

		logrus.WithFields(logrus.Fields{
			"addr": srvAddr,
		}).Info("Sunucu başlatılıyor...")

		// Sunucuyu başlat
		if err := srv.Run(); err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Error("Sunucuyu çalıştırırken bir hata meydana geldi!")
			errChan <- err
		}
	}()

	<-errChan // Bir hata oluşmasını bekle
}
