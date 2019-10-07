package database

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type Komisyon struct {
	AdSoyad string
	Dahil   bool
}

func (kom *Komisyon) Insert(conn *Connection) error {
	const sql string = "INSERT INTO komisyon (AdSoyad) VALUES (?);"

	_, err := conn.db.Exec(sql, kom.AdSoyad)
	if err == nil {
		return nil
	}

	me, ok := err.(*mysql.MySQLError)
	if ok {
		if me.Number == 1062 {
			return errors.New("Bu komisyon üyesi zaten mevcut!")
		}
		if me.Number == 1406 {
			return errors.New("Komisyon üye ismi çok uzun!")
		}
	}

	logrus.WithFields(logrus.Fields{
		"err": err,
	}).Error("Komisyon üyesi eklenirken veritabanında bir hata oluştu!")

	return errors.New("Veritabanında bir hata oluştu!")
}

func (kom *Komisyon) Update(conn *Connection) error {
	const sql string = "UPDATE komisyon SET Dahil = (?) WHERE AdSoyad = (?);"

	_, err := conn.db.Exec(sql, kom.Dahil, kom.AdSoyad)
	if err != nil {
		return err
	}

	return nil
}

func (kom *Komisyon) Delete(conn *Connection) error {
	const sql string = "DELETE FROM komisyon WHERE AdSoyad = (?);"

	_, err := conn.db.Exec(sql, kom.AdSoyad)
	if err != nil {
		return err
	}

	return nil
}

func KomisyonListesi(conn *Connection) ([]Komisyon, error) {
	const sql string = `SELECT AdSoyad,Dahil FROM komisyon`

	q, err := conn.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	liste := []Komisyon{}
	for q.Next() {
		var AdSoyad string
		var Dahil bool
		err = q.Scan(&AdSoyad, &Dahil)
		if err != nil {
			return nil, err
		}

		liste = append(liste, Komisyon{AdSoyad, Dahil})
	}

	return liste, nil
}
