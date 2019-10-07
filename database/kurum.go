package database

func KurumListesi(conn *Connection) ([]string, error) {
	const sql string = `SELECT kurumadi FROM kurum`

	q, err := conn.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	liste := []string{}
	for q.Next() {
		var baslik string

		err = q.Scan(&baslik)
		if err != nil {
			return nil, err
		}

		liste = append(liste, baslik)
	}

	return liste, nil
}
