package main

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

func LogBaslat(logpath string) {
	// Dosya belirtilmemişse sadece stdout'a yazdır
	if logpath == "" {
		logrus.SetOutput(os.Stdout)
		return
	}

	// Hem dosyaya hem stdout'a yazdır
	logfile, err := os.OpenFile(logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.SetOutput(os.Stdout)
		logrus.WithFields(logrus.Fields{
			"filename": logpath,
			"error":    err,
		}).Error("Log dosyası açılamadı!")
		return
	}

	confile := io.MultiWriter(os.Stdout, logfile)
	logrus.SetOutput(confile)
}
