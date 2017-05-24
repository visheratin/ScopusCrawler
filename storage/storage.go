package storage

import "github.com/visheratin/scopus-crawler/crawler"

// GenericStorage is a general interface for data storage
type GenericStorage interface {
	Init() error

	CreateAuthor(author crawler.Author) error
	UpdateAuthor(author crawler.Author) error
	GetAuthor(scopusID string) (crawler.Author, error)
	DeleteAuthor(scopusID string) error

	CreateAffiliation(affiliation crawler.Affiliation) error
	UpdateAffiliation(affiliation crawler.Affiliation) error
	GetAffiliation(scopusID string) (crawler.Affiliation, error)
	DeleteAffiliation(scopusID string) error

	CreateArticle(article crawler.Article) error
	UpdateArticle(article crawler.Article) error
	GetArticle(scopusID string) (crawler.Article, error)
	DeleteArticle(scopusID string) error

	CreateSubjectArea(article crawler.Article) error
	UpdateSubjectArea(article crawler.Article) error
	GetSubjectArea(scopusID string) (crawler.SubjectArea, error)
	DeleteSubjectArea(scopusID string) error

	CreateKeyword(keyword crawler.Keyword) error
	UpdateKeyword(article crawler.Keyword) error
	GetKeyword(id string) (crawler.Keyword, error)
	DeleteKeyword(id string) error
}
