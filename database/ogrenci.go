package database

import (
	"errors"

	"stajtakip/cfg"

	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type Ogrenci struct {
	No      int
	Ad      string
	Soyad   string
	Ogretim int
}

func (ogr *Ogrenci) Insert(conn *Connection) error {
	const sql string = "INSERT INTO ogrenci (No, Ad, Soyad, Ogretim) VALUES (?, ?, ?, ?);"

	_, err := conn.db.Exec(sql, ogr.No, ogr.Ad, ogr.Soyad, ogr.Ogretim)
	if err == nil {
		return nil
	}

	me, ok := err.(*mysql.MySQLError)
	if ok {
		if me.Number == 1062 {
			return errors.New("Bu öğrenci zaten var!")
		}
	}

	logrus.WithFields(logrus.Fields{
		"err": err,
	}).Error("Öğrenci eklenirken veritabanında bir hata oluştu!")

	return errors.New("Veritabanında bir hata oluştu!")
}

func OgrenciBul(conn *Connection, no int) (*Ogrenci, error) {
	const sql string = `SELECT No, Ad, Soyad, Ogretim FROM ogrenci WHERE No=?`

	q, err := conn.db.Query(sql, no)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	if !q.Next() {
		return nil, ErrVeriBulunamadi
	}

	var ogr Ogrenci
	if err = q.Scan(&ogr.No, &ogr.Ad, &ogr.Soyad, &ogr.Ogretim); err != nil {
		return nil, err
	}

	return &ogr, nil
}

func StajiTamamOgrenciler(conn *Connection) ([]Ogrenci, error) {
	const sql string = `SELECT No, Ad, Soyad, Ogretim FROM
(SELECT o.No, o.Ad, o.Soyad, o.Ogretim, SUM(s.Kabulgun) as KabulGun, SUM(s.ToplamGun) as ToplamGun
FROM ogrenci AS o, staj AS s
WHERE o.No = s.OgrenciNo
GROUP BY No, Ad, Soyad, Ogretim
UNION
SELECT o.No, o.Ad, o.Soyad, o.Ogretim, SUM(d.Kabulgun) as KabulGun, SUM(d.ToplamGun) as ToplamGun
FROM ogrenci AS o, denkstaj AS d
WHERE o.No = d.OgrenciNo
GROUP BY No, Ad, Soyad, Ogretim) AS GenelStaj
GROUP BY No, Ad, Soyad, Ogretim
HAVING SUM(Kabulgun) >= ?
AND SUM(ToplamGun) >= 60`

	q, err := conn.db.Query(sql, cfg.GerekenStajGunu())
	if err != nil {
		return nil, err
	}
	defer q.Close()

	liste := []Ogrenci{}
	for q.Next() {
		var ogr Ogrenci

		err = q.Scan(&ogr.No, &ogr.Ad, &ogr.Soyad, &ogr.Ogretim)
		if err != nil {
			return nil, err
		}

		liste = append(liste, ogr)
	}

	return liste, nil
}

// TODO update, delete fonksiyonları
