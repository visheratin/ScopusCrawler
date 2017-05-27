package storage

import (
	"errors"

	"github.com/visheratin/scopus-crawler/crawler"
)

func (storage SQLiteStorage) CreateArticle(article crawler.Article) error {
	db, err := storage.getDBConnection()
	if err != nil {
		return err
	}
	req, _ := db.Prepare("INSERT OR REPLACE INTO articles VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	_, err = req.Exec(article.ScopusID, article.Title, article.Abstracts, article.PublicationDate,
		article.CitationsCount, article.PublicationType, article.PublicationTitle, article.AffiliationID)
	if err != nil {
		return err
	}
	return nil
}

func (storage SQLiteStorage) UpdateArticle(article crawler.Article) error {
	db, err := storage.getDBConnection()
	if err != nil {
		return err
	}
	req, _ := db.Prepare(`UPDATE articles 
		SET title = ?, abstracts = ?, publication_date = ?, citations_count = ?, publication_type = ?, 
		publication_title = ?, affiliation_id = ?
		WHERE scopus_id = ?`)
	_, err = req.Exec(article.Title, article.Abstracts, article.PublicationDate, article.CitationsCount,
		article.PublicationType, article.PublicationTitle, article.AffiliationID)
	if err != nil {
		return err
	}
	return nil
}

func (storage SQLiteStorage) GetArticle(scopusID string) (crawler.Article, error) {
	var article crawler.Article
	db, err := storage.getDBConnection()
	if err != nil {
		return article, err
	}
	req, _ := db.Prepare(`SELECT DISTINCT * FROM articles WHERE scopus_id = ?`)
	res, err := req.Query(scopusID)
	defer res.Close()
	if err != nil {
		return article, err
	}
	for res.Next() {
		err = res.Scan(&article.ScopusID, &article.Title, &article.Abstracts,
			&article.PublicationDate, &article.CitationsCount, &article.PublicationType,
			&article.PublicationTitle, &article.AffiliationID)
		if err != nil {
			return article, err
		}
		return article, nil
	}
	return article, errors.New("data was not found in the storage")
}

func (storage SQLiteStorage) SearchArticles(fields map[string]string) ([]crawler.Article, error) {
	var articles []crawler.Article
	db, err := storage.getDBConnection()
	if err != nil {
		return articles, err
	}
	query := "SELECT DISTINCT * FROM articles WHERE "
	for key, value := range fields {
		query += key + "=" + value + " AND "
	}
	query = query[:(len(query) - 4)]
	res, err := db.Query(query)
	if err != nil {
		return articles, err
	}
	for res.Next() {
		var article crawler.Article
		err = res.Scan(&article.ScopusID, &article.Title, &article.Abstracts,
			&article.PublicationDate, &article.CitationsCount, &article.PublicationType,
			&article.PublicationTitle, &article.AffiliationID)
		if err != nil {
			return articles, err
		}
		articles = append(articles, article)
	}
	return articles, nil
}

func (storage SQLiteStorage) DeleteArticle(scopusID string) error {
	db, err := storage.getDBConnection()
	if err != nil {
		return err
	}
	req, _ := db.Prepare(`DELETE FROM articles WHERE scopus_id = ?`)
	_, err = req.Exec(scopusID)
	if err != nil {
		return err
	}
	return nil
}
