package database

type TumStaj struct {
	KurumAdi        string
	ToplamGun       int
}

func OgrenciTumStajListesi(conn *Connection, no int) ([]TumStaj, error) {
	const sql string = `SELECT KurumAdi, SUM(ToplamGun) AS ToplamGun
FROM denkstaj
WHERE OgrenciNo=?
GROUP BY KurumAdi
UNION
SELECT KurumAdi, SUM(ToplamGun) AS ToplamGun
FROM staj
WHERE OgrenciNo=?
GROUP BY KurumAdi`

	q, err := conn.db.Query(sql, no, no)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	liste := []TumStaj{}
	for q.Next() {
		var stj TumStaj

		err = q.Scan(&stj.KurumAdi, &stj.ToplamGun)
		if err != nil {
			return nil, err
		}

		liste = append(liste, stj)
	}

	return liste, nil
}
