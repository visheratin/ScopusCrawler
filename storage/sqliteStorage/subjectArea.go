package storage

import (
	"errors"

	"github.com/visheratin/scopus-crawler/crawler"
)

func (storage SQLiteStorage) CreateSubjectArea(subjectArea crawler.SubjectArea) error {
	db, err := storage.getDBConnection()
	if err != nil {
		return err
	}
	req, _ := db.Prepare("INSERT OR REPLACE INTO subject_areas VALUES (?, ?, ?, ?)")
	_, err = req.Exec(subjectArea.ScopusID, subjectArea.Title, subjectArea.Code, subjectArea.Description)
	if err != nil {
		return err
	}
	return nil
}

func (storage SQLiteStorage) UpdateSubjectArea(subjectArea crawler.SubjectArea) error {
	db, err := storage.getDBConnection()
	if err != nil {
		return err
	}
	req, _ := db.Prepare(`UPDATE subject_areas 
		SET title = ?, code = ?, description = ?
		WHERE scopus_id = ?`)
	_, err = req.Exec(subjectArea.Title, subjectArea.Code, subjectArea.Description, subjectArea.ScopusID)
	if err != nil {
		return err
	}
	return nil
}

func (storage SQLiteStorage) GetSubjectArea(scopusID string) (crawler.SubjectArea, error) {
	var subjectArea crawler.SubjectArea
	db, err := storage.getDBConnection()
	if err != nil {
		return subjectArea, err
	}
	req, _ := db.Prepare(`SELECT DISTINCT * FROM subject_areas WHERE scopus_id = ?`)
	res, err := req.Query(scopusID)
	defer res.Close()
	if err != nil {
		return subjectArea, err
	}
	for res.Next() {
		err = res.Scan(&subjectArea.ScopusID, &subjectArea.Title, &subjectArea.Code,
			&subjectArea.Description)
		if err != nil {
			return subjectArea, err
		}
		return subjectArea, nil
	}
	return subjectArea, errors.New("data was not found in the storage")
}

func (storage SQLiteStorage) SearchSubjectAreas(fields map[string]string) ([]crawler.SubjectArea, error) {
	var subjectAreas []crawler.SubjectArea
	db, err := storage.getDBConnection()
	if err != nil {
		return subjectAreas, err
	}
	query := "SELECT DISTINCT * FROM subjectAreas WHERE "
	for key, value := range fields {
		query = query + key + "=" + value + " AND "
	}
	query = query[:(len(query) - 4)]
	res, err := db.Query(query)
	if err != nil {
		return subjectAreas, err
	}
	for res.Next() {
		var subjectArea crawler.SubjectArea
		err = res.Scan(&subjectArea.ScopusID, &subjectArea.Title, &subjectArea.Code,
			&subjectArea.Description)
		if err != nil {
			return subjectAreas, err
		}
		subjectAreas = append(subjectAreas, subjectArea)
	}
	return subjectAreas, nil
}

func (storage SQLiteStorage) DeleteSubjectArea(scopusID string) error {
	db, err := storage.getDBConnection()
	if err != nil {
		return err
	}
	req, _ := db.Prepare(`DELETE FROM subject_areas WHERE scopus_id = ?`)
	_, err = req.Exec(scopusID)
	if err != nil {
		return err
	}
	return nil
}
