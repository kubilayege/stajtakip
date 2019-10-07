package database

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type Konu struct {
	Baslik string
	Aktif   bool
}

func (konu *Konu) Insert(conn *Connection) error {
	const sql string = "INSERT INTO konu (Baslik) VALUES (?);"

	_, err := conn.db.Exec(sql, konu.Baslik)
	if err == nil {
		return nil
	}

	me, ok := err.(*mysql.MySQLError)
	if ok {
		if me.Number == 1062 {
			return errors.New("Bu konu zaten mevcut!")
		}
		if me.Number == 1406 {
			return errors.New("Konu ismi çok uzun!")
		}
	}

	logrus.WithFields(logrus.Fields{
		"err": err,
	}).Error("Konu eklenirken veritabanında bir hata oluştu!")

	return errors.New("Veritabanında bir hata oluştu!")
}

func(konu *Konu) Update(conn *Connection) error {
	const sql string = "UPDATE konu SET Aktif = (?) WHERE Baslik = (?);"

	_, err := conn.db.Exec(sql, konu.Aktif, konu.Baslik)
	if err != nil {
		return err
	}

	return nil
}

func(konu *Konu) Delete(conn *Connection) error {
	const sql string = "DELETE FROM konu WHERE Baslik = (?);"

	_, err := conn.db.Exec(sql, konu.Baslik)
	if err != nil {
		return err
	}

	return nil
}

func KonuListesi(conn *Connection) ([]Konu, error) {
	const sql string = `SELECT Baslik, Aktif FROM konu`

	q, err := conn.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	liste := []Konu{}
	for q.Next() {
		var Baslik string
		var Aktif bool
		err = q.Scan(&Baslik, &Aktif)
		if err != nil {
			return nil, err
		}

		liste = append(liste, Konu{Baslik,Aktif})
	}

	return liste, nil
}
