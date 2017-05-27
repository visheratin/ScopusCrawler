package storage

import (
	"database/sql"
	"fmt"
)

type SQLiteStorage struct {
	Path        string
	Initialized bool
	DB          *sql.DB
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

const createFinishedRequestsTable = `CREATE TABLE IF NOT EXISTS finished_requests(
	request VARCHAR PRIMARY KEY
)`

// Init creates new storage or initializes the existing one
func (storage SQLiteStorage) Init(keepAlive bool) error {
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
	_, err = db.Exec(createFinishedRequestsTable)
	if err != nil {
		return err
	}
	if keepAlive {
		storage.DB = db
	} else {
		db.Close()
	}
	storage.Initialized = true
	return nil
}

func (storage SQLiteStorage) Close() {
	if storage.DB != nil {
		storage.DB.Close()
	}
	storage.Initialized = false
}

func (storage SQLiteStorage) getDBConnection() (*sql.DB, error) {
	var err error
	if !storage.Initialized {
		err = storage.Init(false)
		if err != nil {
			return nil, err
		}
	}
	var db *sql.DB
	if storage.DB != nil {
		db = storage.DB
	} else {
		db, err = sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", storage.Path))
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}
