package storage

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/visheratin/scopus-crawler/logger"
	"github.com/visheratin/scopus-crawler/models"
)

type MySQLStorage struct {
	Address     string
	User        string
	Password    string
	DbName      string
	Initialized bool
	DB          *sql.DB
}

const createAuthorsTable = `CREATE TABLE IF NOT EXISTS authors (
	scopus_id VARCHAR(20),
	affiliation_id VARCHAR(20),
	initials TEXT,
	indexed_name TEXT,
	surname TEXT,
	name TEXT,
	PRIMARY KEY (scopus_id),
	FOREIGN KEY(affiliation_id) REFERENCES affiliations(scopus_id)
)`

const createAffiliationsTable = `CREATE TABLE IF NOT EXISTS affiliations (
	scopus_id VARCHAR(20),
	title TEXT,
	country TEXT,
	city TEXT,
	state TEXT,
	postal_code TEXT,
	address TEXT,
    PRIMARY KEY (scopus_id)
)`

const createArticlesTable = `CREATE TABLE IF NOT EXISTS articles (
	scopus_id VARCHAR(20),
	title TEXT,
	abstracts TEXT,
	publication_date TEXT,
	citations_count INTEGER,
	publication_type TEXT,
	publication_title TEXT,
	PRIMARY KEY (scopus_id)
)`

const createSubjectAreasTable = `CREATE TABLE IF NOT EXISTS subject_areas (
	scopus_id VARCHAR(20),
	title TEXT,
	code TEXT,
	description TEXT,
	PRIMARY KEY (scopus_id)
)`

const createKeywordsTable = `CREATE TABLE IF NOT EXISTS keywords (
	id VARCHAR(20),
	keyword TEXT,
	PRIMARY KEY (id)
)`

const createArticleAuthorsTable = `CREATE TABLE IF NOT EXISTS article_author(
	author_id VARCHAR(20),
	article_id VARCHAR(20),
	FOREIGN KEY(author_id) REFERENCES authors(scopus_id) ON DELETE CASCADE,
	FOREIGN KEY(article_id) REFERENCES articles(scopus_id) ON DELETE CASCADE
)`

const createArticleArticlesTable = `CREATE TABLE IF NOT EXISTS article_article(
	from_id VARCHAR(20),
	to_id VARCHAR(20),
	FOREIGN KEY(from_id) REFERENCES articles(scopus_id) ON DELETE CASCADE,
	FOREIGN KEY(to_id) REFERENCES articles(scopus_id) ON DELETE CASCADE
)`

const createArticleAreasTable = `CREATE TABLE IF NOT EXISTS article_area(
	area_id VARCHAR(20),
	article_id VARCHAR(20),
	FOREIGN KEY(article_id) REFERENCES articles(scopus_id) ON DELETE CASCADE,
	FOREIGN KEY(area_id) REFERENCES subject_areas(scopus_id) ON DELETE CASCADE
)`

const createArticleKeywordsTable = `CREATE TABLE IF NOT EXISTS article_keyword(
	keyword_id VARCHAR(20),
	article_id VARCHAR(20),
	FOREIGN KEY(article_id) REFERENCES articles(scopus_id) ON DELETE CASCADE,
	FOREIGN KEY(keyword_id) REFERENCES keywords(id) ON DELETE CASCADE
)`

const createFinishedRequestsTable = `CREATE TABLE IF NOT EXISTS finished_requests(
	request VARCHAR(256) PRIMARY KEY,
	response TEXT
)`

// Init creates new storage or initializes the existing one
func (storage MySQLStorage) Init(keepAlive bool) error {
	path := fmt.Sprintf("%s:%s@(%s)/%s", storage.User, storage.Password, storage.Address, storage.DbName)
	db, err := sql.Open("mysql", path)
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

func (storage MySQLStorage) Close() {
	if storage.DB != nil {
		storage.DB.Close()
	}
	storage.Initialized = false
}

func (storage MySQLStorage) getDBConnection() (*sql.DB, error) {
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
		path := fmt.Sprintf("%s:%s@(%s)/%s", storage.User, storage.Password, storage.Address, storage.DbName)
		db, err = sql.Open("mysql", path)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

func (storage MySQLStorage) CreateAffiliation(affiliation models.Affiliation) error {
	db, err := storage.getDBConnection()
	if err != nil {
		return err
	}
	req, _ := db.Prepare("REPLACE INTO affiliations VALUES (?, ?, ?, ?, ?, ?, ?)")
	_, err = req.Exec(affiliation.ScopusID, affiliation.Title, affiliation.Country, affiliation.City,
		affiliation.State, affiliation.PostalCode, affiliation.Address)
	if err != nil {
		return err
	}
	return nil
}

func (storage MySQLStorage) UpdateAffiliation(affiliation models.Affiliation) error {
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

func (storage MySQLStorage) GetAffiliation(scopusID string) (models.Affiliation, error) {
	var affiliation models.Affiliation
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

func (storage MySQLStorage) SearchAffiliations(fields map[string]string) ([]models.Affiliation, error) {
	var affiliations []models.Affiliation
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
		var affiliation models.Affiliation
		err = res.Scan(&affiliation.ScopusID, &affiliation.Title, &affiliation.Country,
			&affiliation.City, &affiliation.State, &affiliation.PostalCode, &affiliation.Address)
		if err != nil {
			return affiliations, err
		}
		affiliations = append(affiliations, affiliation)
	}
	return affiliations, nil
}

func (storage MySQLStorage) DeleteAffiliation(scopusID string) error {
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
func (storage MySQLStorage) CreateArticle(article models.Article) error {
	db, err := storage.getDBConnection()
	if err != nil {
		return err
	}
	req, _ := db.Prepare("REPLACE INTO articles VALUES (?, ?, ?, ?, ?, ?, ?)")
	_, err = req.Exec(article.ScopusID, article.Title, article.Abstracts, article.PublicationDate,
		article.CitationsCount, article.PublicationType, article.PublicationTitle)
	if err != nil {
		return err
	}
	for _, affiliation := range article.Affiliations {
		err = storage.CreateAffiliation(affiliation)
		if err != nil {
			logger.Error.Println("Unable to add affiliation " + affiliation.ScopusID + " to storage")
			logger.Error.Println(err)
		}
	}
	for _, area := range article.SubjectAreas {
		err = storage.CreateSubjectArea(area)
		if err != nil {
			logger.Error.Println("Unable to add subject area " + area.ScopusID + " to storage")
			logger.Error.Println(err)
		} else {
			req, _ := db.Prepare("INSERT INTO article_area VALUES(?, ?)")
			_, err = req.Exec(area.ScopusID, article.ScopusID)
			if err != nil {
				logger.Error.Println("Unable to connect article " + article.ScopusID + " with area " + area.ScopusID)
				logger.Error.Println(err)
			}
		}
	}
	for _, author := range article.Authors {
		err = storage.CreateAuthor(author)
		if err != nil {
			logger.Error.Println("Unable to add author " + author.ScopusID + " to storage")
			logger.Error.Println(err)
		} else {
			req, _ := db.Prepare("INSERT INTO article_author VALUES(?, ?)")
			_, err = req.Exec(author.ScopusID, article.ScopusID)
			if err != nil {
				logger.Error.Println("Unable to connect article " + article.ScopusID + " with author " + author.ScopusID)
				logger.Error.Println(err)
			}
		}
	}
	for _, keyword := range article.Keywords {
		err = storage.CreateKeyword(keyword)
		if err != nil {
			logger.Error.Println("Unable to add keyword " + keyword.ID + " to storage")
			logger.Error.Println(err)
		} else {
			req, _ := db.Prepare("INSERT INTO article_keyword VALUES(?, ?)")
			_, err = req.Exec(keyword.ID, article.ScopusID)
			if err != nil {
				logger.Error.Println("Unable to connect article " + article.ScopusID + " with keyword " + keyword.ID)
				logger.Error.Println(err)
			}
		}
	}
	for _, reference := range article.References {
		req, _ := db.Prepare("INSERT INTO article_article VALUES(?, ?)")
		_, err = req.Exec(article.ScopusID, reference.ScopusID)
		if err != nil {
			logger.Error.Println("Unable to connect article " + article.ScopusID + " with reference " + reference.ScopusID)
			logger.Error.Println(err)
		}
	}
	return nil
}

func (storage MySQLStorage) UpdateArticle(article models.Article) error {
	db, err := storage.getDBConnection()
	if err != nil {
		return err
	}
	req, _ := db.Prepare(`UPDATE articles 
		SET title = ?, abstracts = ?, publication_date = ?, citations_count = ?, publication_type = ?, 
		publication_title = ?
		WHERE scopus_id = ?`)
	_, err = req.Exec(article.Title, article.Abstracts, article.PublicationDate, article.CitationsCount,
		article.PublicationType, article.PublicationTitle)
	if err != nil {
		return err
	}
	return nil
}

func (storage MySQLStorage) GetArticle(scopusID string) (models.Article, error) {
	var article models.Article
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
			&article.PublicationTitle)
		if err != nil {
			return article, err
		}
		return article, nil
	}
	return article, errors.New("data was not found in the storage")
}

func (storage MySQLStorage) SearchArticles(fields map[string]string) ([]models.Article, error) {
	var articles []models.Article
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
		var article models.Article
		err = res.Scan(&article.ScopusID, &article.Title, &article.Abstracts,
			&article.PublicationDate, &article.CitationsCount, &article.PublicationType,
			&article.PublicationTitle)
		if err != nil {
			return articles, err
		}
		articles = append(articles, article)
	}
	return articles, nil
}

func (storage MySQLStorage) DeleteArticle(scopusID string) error {
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

func (storage MySQLStorage) CreateAuthor(author models.Author) error {
	db, err := storage.getDBConnection()
	if err != nil {
		return err
	}
	req, _ := db.Prepare("REPLACE INTO authors VALUES (?, ?, ?, ?, ?, ?)")
	_, err = req.Exec(author.ScopusID, author.AffiliationID, author.Initials,
		author.IndexedName, author.Surname, author.Name)
	if err != nil {
		return err
	}
	return nil
}

func (storage MySQLStorage) UpdateAuthor(author models.Author) error {
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

func (storage MySQLStorage) GetAuthor(scopusID string) (models.Author, error) {
	var author models.Author
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

func (storage MySQLStorage) SearchAuthors(fields map[string]string) ([]models.Author, error) {
	var authors []models.Author
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
		var author models.Author
		err = res.Scan(&author.ScopusID, &author.AffiliationID, &author.Initials, &author.IndexedName, &author.Surname, &author.Name)
		if err != nil {
			return authors, err
		}
		authors = append(authors, author)
	}
	return authors, nil
}

func (storage MySQLStorage) DeleteAuthor(scopusID string) error {
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

func (storage MySQLStorage) CreateFinishedRequest(request string, response string) error {
	db, err := storage.getDBConnection()
	if err != nil {
		return err
	}
	req, _ := db.Prepare("REPLACE INTO finished_requests VALUES (?, ?)")
	_, err = req.Exec(request, response)
	if err != nil {
		return err
	}
	return nil
}

func (storage MySQLStorage) GetFinishedRequest(request string) (string, error) {
	db, err := storage.getDBConnection()
	if err != nil {
		return "", err
	}
	res, err := db.Query("SELECT DISTINCT response FROM finished_requests WHERE request=" + request)
	if err != nil {
		return "", err
	}
	for res.Next() {
		var response string
		err = res.Scan(&response)
		if err != nil {
			return "", err
		}
		return response, nil
	}
	return "", nil
}

func (storage MySQLStorage) CreateKeyword(keyword models.Keyword) error {
	db, err := storage.getDBConnection()
	if err != nil {
		return err
	}
	req, _ := db.Prepare("REPLACE INTO keywords VALUES (?, ?)")
	_, err = req.Exec(keyword.ID, keyword.Value)
	if err != nil {
		return err
	}
	return nil
}

func (storage MySQLStorage) UpdateKeyword(keyword models.Keyword) error {
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

func (storage MySQLStorage) GetKeyword(id string) (models.Keyword, error) {
	var keyword models.Keyword
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

func (storage MySQLStorage) SearchKeywords(fields map[string]string) ([]models.Keyword, error) {
	var keywords []models.Keyword
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
		var keyword models.Keyword
		err = res.Scan(&keyword.ID, &keyword.Value)
		if err != nil {
			return keywords, err
		}
		keywords = append(keywords, keyword)
	}
	return keywords, nil
}

func (storage MySQLStorage) DeleteKeyword(id string) error {
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

func (storage MySQLStorage) CreateSubjectArea(subjectArea models.SubjectArea) error {
	db, err := storage.getDBConnection()
	if err != nil {
		return err
	}
	req, _ := db.Prepare("REPLACE INTO subject_areas VALUES (?, ?, ?, ?)")
	_, err = req.Exec(subjectArea.ScopusID, subjectArea.Title, subjectArea.Code, subjectArea.Description)
	if err != nil {
		return err
	}
	return nil
}

func (storage MySQLStorage) UpdateSubjectArea(subjectArea models.SubjectArea) error {
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

func (storage MySQLStorage) GetSubjectArea(scopusID string) (models.SubjectArea, error) {
	var subjectArea models.SubjectArea
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

func (storage MySQLStorage) SearchSubjectAreas(fields map[string]string) ([]models.SubjectArea, error) {
	var subjectAreas []models.SubjectArea
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
		var subjectArea models.SubjectArea
		err = res.Scan(&subjectArea.ScopusID, &subjectArea.Title, &subjectArea.Code,
			&subjectArea.Description)
		if err != nil {
			return subjectAreas, err
		}
		subjectAreas = append(subjectAreas, subjectArea)
	}
	return subjectAreas, nil
}

func (storage MySQLStorage) DeleteSubjectArea(scopusID string) error {
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
