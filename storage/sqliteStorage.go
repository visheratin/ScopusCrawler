package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/visheratin/scopus-crawler/crawler"
)

type SQLiteStorage struct {
	Path string
}

const createAuthorsTable = `CREATE TABLE IF NOT EXISTS authors (
	scopus_id VARCHAR PRIMARY KEY,
	affiliation_id VARCHAR,
	initials VARCHAR,
	indexed_name VARCHAR,
	surname VARCHAR,
	name VARCHAR,
	FOREIGN KEY(affiliation_id) REFERENCES affiliations(scopus_id)
)`

const createAffiliationsTable = `CREATE TABLE IF NOT EXISTS affiliations (
	scopus_id VARCHAR PRIMARY KEY,
	title VARCHAR,
	country VARCHAR,
	city VARCHAR,
	state VARCHAR,
	postal_code VARCHAR,
	address VARCHAR
)`

const createArticlesTable = `CREATE TABLE IF NOT EXISTS articles (
	scopus_id VARCHAR PRIMARY KEY,
	title VARCHAR,
	abstracts VARCHAR,
	publication_date VARCHAR,
	citations_count INTEGER,
	publication_type VARCHAR,
	publication_title VARCHAR,
	affiliation_id VARCHAR,
	FOREIGN KEY(affiliation_id) REFERENCES affiliations(scopus_id)
)`

const createSubjectAreasTable = `CREATE TABLE IF NOT EXISTS subject_areas (
	scopus_id VARCHAR PRIMARY KEY,
	title VARCHAR,
	code VARCHAR,
	description VARCHAR
)`

const createKeywordsTable = `CREATE TABLE IF NOT EXISTS keywords (
	id VARCHAR PRIMARY KEY,
	keyword VARCHAR
)`

const createArticleAuthorsTable = `CREATE TABLE IF NOT EXISTS article_author(
	author_id VARCHAR,
	article_id VARCHAR,
	FOREIGN KEY(author_id) REFERENCES authors(scopus_id) ON DELETE CASCADE,
	FOREIGN KEY(article_id) REFERENCES articles(scopus_id) ON DELETE CASCADE
)`

const createArticleArticlesTable = `CREATE TABLE IF NOT EXISTS article_article(
	from_id VARCHAR,
	to_id VARCHAR,
	FOREIGN KEY(from_id) REFERENCES articles(scopus_id) ON DELETE CASCADE,
	FOREIGN KEY(to_id) REFERENCES articles(scopus_id) ON DELETE CASCADE
)`

const createArticleAreasTable = `CREATE TABLE IF NOT EXISTS article_area(
	area_id VARCHAR,
	article_id VARCHAR,
	FOREIGN KEY(article_id) REFERENCES articles(scopus_id) ON DELETE CASCADE,
	FOREIGN KEY(area_id) REFERENCES subject_areas(scopus_id) ON DELETE CASCADE
)`

const createArticleKeywordsTable = `CREATE TABLE IF NOT EXISTS article_keyword(
	keyword_id VARCHAR,
	article_id VARCHAR,
	FOREIGN KEY(article_id) REFERENCES articles(scopus_id) ON DELETE CASCADE,
	FOREIGN KEY(keyword_id) REFERENCES keywords(id) ON DELETE CASCADE
)`

// Init creates new storage or initializes the existing one
func (storage SQLiteStorage) Init() error {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", storage.Path))
	if err != nil {
		return err
	}
	_, err = db.Exec(createAffiliationsTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(createAuthorsTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(createKeywordsTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(createSubjectAreasTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(createArticlesTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(createArticleAreasTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(createArticleArticlesTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(createArticleAuthorsTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(createArticleKeywordsTable)
	if err != nil {
		return err
	}
	return nil
}

func (storage SQLiteStorage) CreateAuthor(author crawler.Author) error {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", storage.Path))
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
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", storage.Path))
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
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", storage.Path))
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

func (storage SQLiteStorage) DeleteAuthor(scopusID string) error {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", storage.Path))
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

func (storage SQLiteStorage) CreateAffiliation(affiliation crawler.Affiliation) error {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", storage.Path))
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
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", storage.Path))
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
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", storage.Path))
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

func (storage SQLiteStorage) DeleteAffiliation(scopusID string) error {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", storage.Path))
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

func (storage SQLiteStorage) CreateKeyword(keyword crawler.Keyword) error {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", storage.Path))
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
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", storage.Path))
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
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", storage.Path))
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

func (storage SQLiteStorage) DeleteKeyword(id string) error {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", storage.Path))
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

func (storage SQLiteStorage) CreateSubjectArea(subjectArea crawler.SubjectArea) error {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", storage.Path))
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
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", storage.Path))
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
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", storage.Path))
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

func (storage SQLiteStorage) DeleteSubjectArea(scopusID string) error {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", storage.Path))
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

func (storage SQLiteStorage) CreateArticle(article crawler.Article) error {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", storage.Path))
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
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", storage.Path))
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
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", storage.Path))
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

func (storage SQLiteStorage) DeleteArticle(scopusID string) error {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", storage.Path))
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
