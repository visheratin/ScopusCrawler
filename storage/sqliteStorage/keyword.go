package storage

import (
	"errors"

	"github.com/visheratin/scopus-crawler/crawler"
)

func (storage SQLiteStorage) CreateKeyword(keyword crawler.Keyword) error {
	db, err := storage.getDBConnection()
	if err != nil {
		return err
	}
	req, _ := db.Prepare("INSERT OR REPLACE INTO keywords VALUES (?, ?)")
	_, err = req.Exec(keyword.ID, keyword.Value)
	if err != nil {
		return err
	}
	return nil
}

func (storage SQLiteStorage) UpdateKeyword(keyword crawler.Keyword) error {
	db, err := storage.getDBConnection()
	if err != nil {
		return err
	}
	req, _ := db.Prepare(`UPDATE keywords SET value = ? WHERE id = ?`)
	_, err = req.Exec(keyword.Value, keyword.ID)
	if err != nil {
		return err
	}
	return nil
}

func (storage SQLiteStorage) GetKeyword(id string) (crawler.Keyword, error) {
	var keyword crawler.Keyword
	db, err := storage.getDBConnection()
	if err != nil {
		return keyword, err
	}
	req, _ := db.Prepare(`SELECT DISTINCT * FROM keywords WHERE id = ?`)
	res, err := req.Query(id)
	defer res.Close()
	if err != nil {
		return keyword, err
	}
	for res.Next() {
		err = res.Scan(&keyword.ID, &keyword.Value)
		if err != nil {
			return keyword, err
		}
		return keyword, nil
	}
	return keyword, errors.New("data was not found in the storage")
}

func (storage SQLiteStorage) SearchKeywords(fields map[string]string) ([]crawler.Keyword, error) {
	var keywords []crawler.Keyword
	db, err := storage.getDBConnection()
	if err != nil {
		return keywords, err
	}
	query := "SELECT DISTINCT * FROM keywords WHERE "
	for key, value := range fields {
		query = query + key + "=" + value + " AND "
	}
	query = query[:(len(query) - 4)]
	res, err := db.Query(query)
	if err != nil {
		return keywords, err
	}
	for res.Next() {
		var keyword crawler.Keyword
		err = res.Scan(&keyword.ID, &keyword.Value)
		if err != nil {
			return keywords, err
		}
		keywords = append(keywords, keyword)
	}
	return keywords, nil
}

func (storage SQLiteStorage) DeleteKeyword(id string) error {
	db, err := storage.getDBConnection()
	if err != nil {
		return err
	}
	req, _ := db.Prepare(`DELETE FROM keywords WHERE id = ?`)
	_, err = req.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
