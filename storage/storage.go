package storage

import "github.com/visheratin/scopus-crawler/models"

// GenericStorage is a general interface for data storage
type GenericStorage interface {
	Init(keepAlive bool) error

	CreateAuthor(author models.Author) error
	UpdateAuthor(author models.Author) error
	GetAuthor(scopusID string) (models.Author, error)
	SearchAuthors(fields map[string]string) ([]models.Author, error)
	DeleteAuthor(scopusID string) error

	CreateAffiliation(affiliation models.Affiliation) error
	UpdateAffiliation(affiliation models.Affiliation) error
	GetAffiliation(scopusID string) (models.Affiliation, error)
	SearchAffiliations(fields map[string]string) ([]models.Affiliation, error)
	DeleteAffiliation(scopusID string) error

	CreateArticle(article models.Article) error
	UpdateArticle(article models.Article) error
	GetArticle(scopusID string) (models.Article, error)
	SearchArticles(fields map[string]string) ([]models.Article, error)
	DeleteArticle(scopusID string) error

	CreateSubjectArea(area models.SubjectArea) error
	UpdateSubjectArea(area models.SubjectArea) error
	GetSubjectArea(scopusID string) (models.SubjectArea, error)
	SearchSubjectAreas(fields map[string]string) ([]models.SubjectArea, error)
	DeleteSubjectArea(scopusID string) error

	CreateKeyword(keyword models.Keyword) error
	UpdateKeyword(article models.Keyword) error
	GetKeyword(id string) (models.Keyword, error)
	SearchKeywords(fields map[string]string) ([]models.Keyword, error)
	DeleteKeyword(id string) error

	CreateFinishedRequest(request string, response string) error
	GetFinishedRequest(request string) (string, error)
}
