package models

type Author struct {
	ScopusID      string
	Initials      string
	IndexedName   string
	Surname       string
	Name          string
	AffiliationID string
	Affiliation   Affiliation
}

type Affiliation struct {
	ScopusID   string
	Title      string
	Country    string
	City       string
	State      string
	PostalCode string
	Address    string
}

type Article struct {
	ScopusID         string
	Title            string
	Abstracts        string
	PublicationDate  string
	CitationsCount   int
	PublicationType  string
	PublicationTitle string
	Affiliations     []Affiliation
	Authors          []Author
	Keywords         []Keyword
	SubjectAreas     []SubjectArea
	References       []Article
}

type SubjectArea struct {
	ScopusID    string
	Title       string
	Code        string
	Description string
}

type Keyword struct {
	ID    string
	Value string
}
