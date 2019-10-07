package database

type RaporSehir struct {
	Sehir       string
	BasariOrani float32
}

type RaporKonu struct {
	Konu        string
	BasariOrani float32
}

type RaporKonuDagilim struct {
	Konu        string
	DagilimOrani float32
}

func RaporSehirler(conn *Connection, baslangic int, bitis int) ([]RaporSehir, error) {
	const sql string = `SELECT
Sehir,
ROUND((SUM(KabulGun)*1.0/SUM(ToplamGun))*100,2)
as "%BasariOrani"
FROM staj
Where Degerlendirildi AND Year(Bitis) BETWEEN ? AND ?
Group BY Sehir;`

	q, err := conn.db.Query(sql, baslangic, bitis)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	liste := []RaporSehir{}
	for q.Next() {
		var Sehir string
		var BasariOrani float32
		err := q.Scan(&Sehir, &BasariOrani)
		if err != nil {
			return nil, err
		}
		liste = append(liste, RaporSehir{Sehir, BasariOrani})
	}

	return liste, nil
}

func RaporKonular(conn *Connection, baslangic int, bitis int) ([]RaporKonu, error) {
	const sql string = `SELECT
KonuBaslik,
ROUND((SUM(KabulGun)*1.0/SUM(ToplamGun))*100,2)
AS "%BasariOrani"
FROM staj
Where Degerlendirildi AND Year(Bitis) BETWEEN ? AND ?
Group BY KonuBaslik;`

	q, err := conn.db.Query(sql, baslangic, bitis)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	liste := []RaporKonu{}
	for q.Next() {
		var Konu string
		var BasariOrani float32

		err := q.Scan(&Konu, &BasariOrani)
		if err != nil {
			return nil, err
		}
		liste = append(liste, RaporKonu{Konu, BasariOrani})
	}

	return liste, nil
}

func RaporKonularDagilim(conn *Connection, baslangic int, bitis int) ([]RaporKonuDagilim, error) {
	const sql string = `SELECT KonuBaslik, ROUND(Sayi / Toplam * 100.0, 2) AS '%DagilimOrani'
FROM
(SELECT KonuBaslik, COUNT(KonuBaslik) AS Sayi
FROM staj
Where Year(Bitis) BETWEEN ? AND ?
GROUP BY KonuBaslik) AS a
JOIN
(SELECT COUNT(*) AS Toplam
FROM staj
Where Year(Bitis) BETWEEN ? AND ?) AS b
ON true`

	q, err := conn.db.Query(sql, baslangic, bitis, baslangic, bitis)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	liste := []RaporKonuDagilim{}
	for q.Next() {
		var Konu string
		var DagilimOrani float32

		err := q.Scan(&Konu, &DagilimOrani)
		if err != nil {
			return nil, err
		}
		liste = append(liste, RaporKonuDagilim{Konu, DagilimOrani})
	}

	return liste, nil
}
