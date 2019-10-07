package database

import (
	dsql "database/sql"
	"github.com/go-sql-driver/mysql"
)

const TarihSaatFormati = "2006-01-02 15:04"
const TarihFormati = "2006-01-02"
const SaatFormati = "15:04"

type MulakatOgrenciStaj struct {
	OgrenciNo      int
	Ad string
	Soyad string
	Ogretim int
	StajBaslangic      string
	KabulGun int
	ToplamGun int
}

type MulakatOgrenci struct {
	OgrenciNo      int
	Ad string
	Soyad string
	Ogretim int
	StajBaslangic      string
	Tarih      string
	Saat      string
	KomisyonUye1      string
	KomisyonUye2      string
}

type Mulakat struct {
	OgrenciNo      int
	StajBaslangic      string
	TarihSaat      string
	KomisyonUye1      string
	KomisyonUye2      string
}

func (mul *Mulakat) Update(conn *Connection) error {
	const sql string = `UPDATE mulakat
SET TarihSaat=?, KomisyonUye1=?, KomisyonUye2=?
WHERE OgrenciNo=? AND StajBaslangic=?`

	result, err := conn.db.Exec(sql, mul.TarihSaat, mul.KomisyonUye1, mul.KomisyonUye2, mul.OgrenciNo, mul.StajBaslangic)
	if err != nil {
		return err
	}

	_, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

type MulakatSonuc struct {
	OgrenciNo      int
	StajBaslangic string
	PuanDevam      int
	PuanCaba      int
	PuanVakit      int
	PuanAmireDavranis int
	PuanIsArkadasaDavranis int
	PuanProje      int
	PuanDuzen      int
	PuanSunum      int
	PuanIcerik      int
	PuanMulakat      int
}

func (mul *MulakatSonuc) Update(conn *Connection) error {
	const sql string = `UPDATE mulakat
SET PuanDevam=?, PuanCaba=?, PuanVakit=?, PuanAmireDavranis=?,
PuanIsArkadasaDavranis=?, PuanProje=?, PuanDuzen=?, PuanSunum=?, PuanIcerik=?, PuanMulakat=?
WHERE OgrenciNo=? AND StajBaslangic=?`
	const sql2 string = `UPDATE staj
SET Degerlendirildi=true, KabulGun=ToplamGun*?
WHERE OgrenciNo=? AND Baslangic=?`

	tx, err := conn.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(
    sql, mul.PuanDevam, mul.PuanCaba, mul.PuanVakit, mul.PuanAmireDavranis,
		mul.PuanIsArkadasaDavranis, mul.PuanProje, mul.PuanDuzen, mul.PuanSunum,
		mul.PuanIcerik, mul.PuanMulakat, mul.OgrenciNo, mul.StajBaslangic)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	var kabulGun float64

	kabulGun = float64(mul.PuanDevam) / 5.0 * 4
	kabulGun += float64(mul.PuanCaba) / 5.0 * 4
	kabulGun += float64(mul.PuanVakit) / 5.0 * 4
	kabulGun += float64(mul.PuanAmireDavranis) / 5.0 * 4
	kabulGun += float64(mul.PuanIsArkadasaDavranis) / 5.0 * 4
	kabulGun += float64(mul.PuanProje) / 100.0 * 15
	kabulGun += float64(mul.PuanDuzen) / 100.0 * 5
	kabulGun += float64(mul.PuanSunum) / 100.0 * 5
	kabulGun += float64(mul.PuanIcerik) / 100.0 * 15
	kabulGun += float64(mul.PuanMulakat) / 100.0 * 40

	_, err = tx.Exec(sql2, kabulGun/100.0, mul.OgrenciNo, mul.StajBaslangic)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func MulakatOgrenciBul(conn *Connection, no int, baslangic string) (*MulakatOgrenci, error) {
	const sql string = `SELECT m.OgrenciNo, o.Ad, o.Soyad, o.Ogretim,
m.StajBaslangic, m.TarihSaat, m.KomisyonUye1, m.KomisyonUye2
FROM mulakat AS m, ogrenci AS o WHERE m.OgrenciNo=o.No AND o.No=? AND m.StajBaslangic=?`

	q, err := conn.db.Query(sql, no, baslangic)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	if !q.Next() {
		return nil, ErrVeriBulunamadi
	}

	if mul, err := MulakatOgrenciScan(q); err != nil {
		return nil, err
	} else {
		return &mul, nil
	}
}

func MulakatListesi(conn *Connection) ([]MulakatOgrenci, error) {
	const sql string = `SELECT m.OgrenciNo, o.Ad, o.Soyad, o.Ogretim,
m.StajBaslangic, m.TarihSaat, m.KomisyonUye1, m.KomisyonUye2
FROM mulakat AS m, ogrenci AS o WHERE m.OgrenciNo=o.No AND m.PuanMulakat IS NULL`

	q, err := conn.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	liste := []MulakatOgrenci{}
	for q.Next() {
		if mul, err := MulakatOgrenciScan(q); err != nil {
			return nil, err
		} else {
			liste = append(liste, mul)
		}
	}

	return liste, nil
}

func MulakatSonucListesi(conn *Connection, baslangic string, bitis string) ([]MulakatOgrenciStaj, error) {
	const sql string = `SELECT m.OgrenciNo, o.Ad, o.Soyad, o.Ogretim,
m.StajBaslangic, s.ToplamGun, s.KabulGun
FROM mulakat AS m, ogrenci AS o, staj AS s WHERE m.OgrenciNo=o.No
AND o.No=s.OgrenciNo AND m.StajBaslangic=s.Baslangic AND s.Degerlendirildi
AND m.TarihSaat BETWEEN ? AND ?`

	q, err := conn.db.Query(sql, baslangic, bitis)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	liste := []MulakatOgrenciStaj{}
	for q.Next() {
		var mul MulakatOgrenciStaj

		err := q.Scan(&mul.OgrenciNo, &mul.Ad, &mul.Soyad, &mul.Ogretim, &mul.StajBaslangic,
								 &mul.ToplamGun, &mul.KabulGun)
		if err != nil {
			return nil, err
		}

		liste = append(liste, mul)
	}

	return liste, nil
}

func MulakatOgrenciScan(q *dsql.Rows) (MulakatOgrenci, error) {
	var mul MulakatOgrenci
	var tarihSaat mysql.NullTime
	var kom1, kom2 dsql.NullString

	err := q.Scan(&mul.OgrenciNo, &mul.Ad, &mul.Soyad, &mul.Ogretim, &mul.StajBaslangic,
							 &tarihSaat, &kom1, &kom2)
	if err != nil {
		return MulakatOgrenci{}, err
	}

	if tarihSaat.Valid {
		mul.Tarih = tarihSaat.Time.Format(TarihFormati)
		mul.Saat = tarihSaat.Time.Format(SaatFormati)
	} else {
		mul.Tarih = ""
		mul.Saat = ""
	}

	if kom1.Valid {
		mul.KomisyonUye1 = kom1.String
	} else {
		mul.KomisyonUye1 = "-"
	}

	if kom2.Valid {
		mul.KomisyonUye2 = kom2.String
	} else {
		mul.KomisyonUye2 = "-"
	}

	return mul, nil
}

func MulakatListesiOlustur(conn *Connection) error {
	const sql string = `INSERT INTO mulakat (OgrenciNo, StajBaslangic)
SELECT s.OgrenciNo, s.Baslangic
FROM staj AS s LEFT JOIN mulakat AS m
ON s.OgrenciNo=m.OgrenciNo AND s.Baslangic=m.StajBaslangic
WHERE m.StajBaslangic IS NULL AND NOT s.Degerlendirildi`

	result, err := conn.db.Exec(sql)
	if err != nil {
		return err
	}

	_, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}
