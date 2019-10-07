package database

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type DenkStaj struct {
	OgrenciNo       int
	KurumAdi        string
	OncekiOkul           string
	ToplamGun       int
	KabulGun        int
}

func (stj *DenkStaj) Insert(conn *Connection) error {
	const sql string = "INSERT INTO denkstaj (OgrenciNo, KurumAdi, OncekiOkul, ToplamGun, KabulGun) VALUES (?, ?, ?, ?, ?);"

	_, err := conn.db.Exec(sql, stj.OgrenciNo, stj.KurumAdi, stj.OncekiOkul, stj.ToplamGun, stj.KabulGun)
	if err == nil {
		return nil
	}

	me, ok := err.(*mysql.MySQLError)
	if ok {
		if me.Number == 1062 {
			return errors.New("Bu öğrenci için dgs/yatay geçiş staj girişi zaten var!")
		}
		if me.Number == 1264 {
			return errors.New("Sayısal değerlerden biri max değeri aşıyor!")
		}
	}

	logrus.WithFields(logrus.Fields{
		"err": err,
	}).Error("DenkStaj eklenirken veritabanında bir hata oluştu!")

	return errors.New("Veritabanında bir hata oluştu!")
}

func OgrenciDenkStajListesi(conn *Connection, no int) ([]DenkStaj, error) {
	const sql string = `SELECT KurumAdi, OncekiOkul, ToplamGun, KabulGun
FROM denkstaj WHERE OgrenciNo=?`

	q, err := conn.db.Query(sql, no)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	liste := []DenkStaj{}
	for q.Next() {
		var stj DenkStaj

		err = q.Scan(&stj.KurumAdi, &stj.OncekiOkul, &stj.ToplamGun, &stj.KabulGun)
		if err != nil {
			return nil, err
		}

		liste = append(liste, stj)
	}

	return liste, nil
}
