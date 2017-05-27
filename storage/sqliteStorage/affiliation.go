package storage

import (
	"errors"

	"github.com/visheratin/scopus-crawler/crawler"
)

func (storage SQLiteStorage) CreateAffiliation(affiliation crawler.Affiliation) error {
	db, err := storage.getDBConnection()
	if err != nil {
		return err
	}
	req, _ := db.Prepare("INSERT OR REPLACE INTO affiliations VALUES (?, ?, ?, ?, ?, ?, ?)")
	_, err = req.Exec(affiliation.ScopusID, affiliation.Title, affiliation.Country, affiliation.City,
		affiliation.State, affiliation.PostalCode, affiliation.Address)
	if err != nil {
		return err
	}
	return nil
}

func (storage SQLiteStorage) UpdateAffiliation(affiliation crawler.Affiliation) error {
	db, err := storage.getDBConnection()
	if err != nil {
		return err
	}
	req, _ := db.Prepare(`UPDATE affiliations 
		SET title = ?, country = ?, city = ?, state = ?, postal_code = ?, address = ?
		WHERE scopus_id = ?`)
	_, err = req.Exec(affiliation.Title, affiliation.Country, affiliation.City, affiliation.State,
		affiliation.PostalCode, affiliation.Address, affiliation.ScopusID)
	if err != nil {
		return err
	}
	return nil
}

func (storage SQLiteStorage) GetAffiliation(scopusID string) (crawler.Affiliation, error) {
	var affiliation crawler.Affiliation
	db, err := storage.getDBConnection()
	if err != nil {
		return affiliation, err
	}
	req, _ := db.Prepare(`SELECT DISTINCT * FROM affiliations WHERE scopus_id = ?`)
	res, err := req.Query(scopusID)
	defer res.Close()
	if err != nil {
		return affiliation, err
	}
	for res.Next() {
		err = res.Scan(&affiliation.ScopusID, &affiliation.Title, &affiliation.Country,
			&affiliation.City, &affiliation.State, &affiliation.PostalCode, &affiliation.Address)
		if err != nil {
			return affiliation, err
		}
		return affiliation, nil
	}
	return affiliation, errors.New("data was not found in the storage")
}

func (storage SQLiteStorage) SearchAffiliations(fields map[string]string) ([]crawler.Affiliation, error) {
	var affiliations []crawler.Affiliation
	db, err := storage.getDBConnection()
	if err != nil {
		return affiliations, err
	}
	query := "SELECT DISTINCT * FROM affiliations WHERE "
	for key, value := range fields {
		query = query + key + "=" + value + " AND "
	}
	query = query[:(len(query) - 4)]
	res, err := db.Query(query)
	if err != nil {
		return affiliations, err
	}
	for res.Next() {
		var affiliation crawler.Affiliation
		err = res.Scan(&affiliation.ScopusID, &affiliation.Title, &affiliation.Country,
			&affiliation.City, &affiliation.State, &affiliation.PostalCode, &affiliation.Address)
		if err != nil {
			return affiliations, err
		}
		affiliations = append(affiliations, affiliation)
	}
	return affiliations, nil
}

func (storage SQLiteStorage) DeleteAffiliation(scopusID string) error {
	db, err := storage.getDBConnection()
	if err != nil {
		return err
	}
	req, _ := db.Prepare(`DELETE FROM affiliations WHERE scopus_id = ?`)
	_, err = req.Exec(scopusID)
	if err != nil {
		return err
	}
	return nil
}
