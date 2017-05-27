package storage

func (storage SQLiteStorage) CreateFinishedRequest(request string) error {
	db, err := storage.getDBConnection()
	if err != nil {
		return err
	}
	req, _ := db.Prepare("INSERT OR REPLACE INTO finished_requests VALUES (?)")
	_, err = req.Exec(request)
	if err != nil {
		return err
	}
	return nil
}

func (storage SQLiteStorage) CheckFinishedRequest(request string) (bool, error) {
	db, err := storage.getDBConnection()
	if err != nil {
		return false, err
	}
	res, err := db.Query("SELECT DISTINCT * FROM finished_requests WHERE request=" + request)
	if err != nil {
		return false, err
	}
	for res.Next() {
		return true, nil
	}
	return false, nil
}
