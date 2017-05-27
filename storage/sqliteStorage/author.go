package storage

import (
	"errors"

	"github.com/visheratin/scopus-crawler/crawler"
)

func (storage SQLiteStorage) CreateAuthor(author crawler.Author) error {
	db, err := storage.getDBConnection()
	if err != nil {
		return err
	}
	req, _ := db.Prepare("INSERT OR REPLACE INTO authors VALUES (?, ?, ?, ?, ?, ?)")
	_, err = req.Exec(author.ScopusID, author.AffiliationID, author.Initials,
		author.IndexedName, author.Surname, author.Name)
	if err != nil {
		return err
	}
	return nil
}

func (storage SQLiteStorage) UpdateAuthor(author crawler.Author) error {
	db, err := storage.getDBConnection()
	if err != nil {
		return err
	}
	req, _ := db.Prepare(`UPDATE authors 
		SET affiliation_id = ?, initials = ?, indexed_name = ?, surname = ?, name = ?
		WHERE scopus_id = ?`)
	_, err = req.Exec(author.AffiliationID, author.Initials,
		author.IndexedName, author.Surname, author.Name, author.ScopusID)
	if err != nil {
		return err
	}
	return nil
}

func (storage SQLiteStorage) GetAuthor(scopusID string) (crawler.Author, error) {
	var author crawler.Author
	db, err := storage.getDBConnection()
	if err != nil {
		return author, err
	}
	req, _ := db.Prepare(`SELECT DISTINCT * FROM authors WHERE scopus_id = ?`)
	res, err := req.Query(scopusID)
	defer res.Close()
	if err != nil {
		return author, err
	}
	for res.Next() {
		err = res.Scan(&author.ScopusID, &author.AffiliationID, &author.Initials, &author.IndexedName, &author.Surname, &author.Name)
		if err != nil {
			return author, err
		}
		return author, nil
	}
	return author, errors.New("data was not found in the storage")
}

func (storage SQLiteStorage) SearchAuthors(fields map[string]string) ([]crawler.Author, error) {
	var authors []crawler.Author
	db, err := storage.getDBConnection()
	if err != nil {
		return authors, err
	}
	query := "SELECT DISTINCT * FROM authors WHERE "
	for key, value := range fields {
		query = query + key + "=" + value + " AND "
	}
	query = query[:(len(query) - 4)]
	res, err := db.Query(query)
	if err != nil {
		return authors, err
	}
	defer res.Close()
	for res.Next() {
		var author crawler.Author
		err = res.Scan(&author.ScopusID, &author.AffiliationID, &author.Initials, &author.IndexedName, &author.Surname, &author.Name)
		if err != nil {
			return authors, err
		}
		authors = append(authors, author)
	}
	return authors, nil
}

func (storage SQLiteStorage) DeleteAuthor(scopusID string) error {
	db, err := storage.getDBConnection()
	if err != nil {
		return err
	}
	req, _ := db.Prepare(`DELETE FROM authors WHERE scopus_id = ?`)
	_, err = req.Exec(scopusID)
	if err != nil {
		return err
	}
	return nil
}
