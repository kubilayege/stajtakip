package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

const DB_PING_TIMEOUT = 10 * time.Second
const DB_PING_TIMER = 1 * time.Minute

var ErrAlreadyConnected error = errors.New("Zaten veritabanı sunucusuna bağlı!")
var ErrVeriBulunamadi error = errors.New("Aranan veri bulunamadı!")

type Connection struct {
	db      *sql.DB
	datasrc string
}

func NewConnection(datasrc string) *Connection {
	return &Connection{
		db:      nil,
		datasrc: datasrc,
	}
}

func (sv *Connection) Connect(ok chan<- bool) error {
	// Veritabanına zaten bağlı
	if sv.db != nil {
		return ErrAlreadyConnected
	}

	// Veritabanına bağlan
	db, err := sql.Open("mysql", sv.datasrc)
	if err != nil {
		return err
	}
	sv.db = db

	// Fonksiyondan çıkarken bağlantıyı kapat
	defer func() {
		db.Close()
		sv.db = nil
	}()

	first := true

	// Veritabanın bağlı gözüktüğü sürece
	for sv.db != nil {
		// Ping isteği belirtilen süreden fazla sürerse zaman aşımına uğrat
		ctx, cancel := context.WithTimeout(context.Background(), DB_PING_TIMEOUT)

		logrus.Debug("Ping isteği gönderiliyor...")

		// Veritabanı sunucusuna ping isteği yolla
		if err := db.PingContext(ctx); err != nil {
			cancel()
			return err
		}
		cancel() // Zaman aşımını bellekten serbest bırak

		// İlk ping geçerliyse onay gönder
		if first {
			first = false
			ok <- true
		}

		// Tekrar ping isteği göndermek için bekle
		time.Sleep(DB_PING_TIMER)
	}

	return nil
}
