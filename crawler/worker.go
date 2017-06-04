package crawler

import (
	"errors"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/visheratin/scopus-crawler/config"
	"github.com/visheratin/scopus-crawler/logger"
	"github.com/visheratin/scopus-crawler/models"
	"github.com/visheratin/scopus-crawler/query"
	"github.com/visheratin/scopus-crawler/storage"
)

type Worker struct {
	Config      config.Configuration
	Storage     storage.GenericStorage
	DataSources []DataSource
	Work        chan SearchRequest
	WorkerQueue chan chan SearchRequest
}

func (worker *Worker) Start() {
	go func() {
		for {
			worker.WorkerQueue <- worker.Work
			work := <-worker.Work
			ds := work.Source
			switch work.SourceName {
			case "search":
				data, err := query.MakeQuery(ds.Path, "", work.Fields, worker.Config.RequestTimeout, worker.Storage)
				if err != nil {
					logger.Error.Println(err)
					return
				}
				articles, err := ExtractArticles(data)
				if err != nil {
					logger.Error.Println(err)
					return
				}
				var articleDs DataSource
				for _, v := range worker.DataSources {
					if v.Name == "article" {
						articleDs = v
						break
					}
				}
				for _, article := range articles {
					err := worker.ProceedArticle(article, articleDs, 0)
					if err != nil {
						logger.Error.Println("Error on proceeding article with id=" + article.ScopusID)
						logger.Error.Println(err)
					}
				}
			}
		}
	}()
}

func ExtractArticles(rawResponse map[string]interface{}) ([]models.Article, error) {
	result := []models.Article{}
	searchContainer, searchSucceed := rawResponse["search-results"]
	if searchSucceed {
		rawEntries, ok := searchContainer.(map[string]interface{})["entry"]
		if !ok {
			return result, errors.New("error on parsing search-results element")
		}
		entries, ok := rawEntries.([]interface{})
		if !ok {
			return result, errors.New("error on parsing entries element")
		}
		for _, value := range entries {
			entry, ok := value.(map[string]interface{})
			if !ok {
				continue
			}
			id := entry["dc:identifier"].(string)
			scopusID := strings.Replace(id, "SCOPUS_ID:", "", -1)
			article := models.Article{
				ScopusID: scopusID,
			}
			title, ok := entry["dc:title"]
			if ok {
				article.Title = title.(string)
			}
			abstracts, ok := entry["dc:description"]
			if ok {
				article.Abstracts = abstracts.(string)
			}
			citations, ok := entry["citedby-count"]
			if ok {
				count, err := strconv.Atoi(citations.(string))
				if err == nil {
					article.CitationsCount = count
				}
			}
			pubDate, ok := entry["prism:coverDate"]
			if ok {
				article.PublicationDate = pubDate.(string)
			}
			pubType, ok := entry["prism:aggregationType"]
			if ok {
				article.PublicationType = pubType.(string)
			}
			pubName, ok := entry["prism:publicationName"]
			if ok {
				article.PublicationTitle = pubName.(string)
			}

			article.Authors = []models.Author{}
			authorsElem := entry["author"]
			authors := authorsElem.([]interface{})
			for _, authorVal := range authors {
				authorElem := authorVal.(map[string]interface{})
				author := models.Author{}
				id, ok := authorElem["authid"]
				if ok {
					author.ScopusID = id.(string)
				}
				name, ok := authorElem["authname"]
				if ok {
					author.IndexedName = name.(string)
				}
				surname, ok := authorElem["surname"]
				if ok {
					author.Surname = surname.(string)
				}
				firstname, ok := authorElem["given-name"]
				if ok {
					author.Name = firstname.(string)
				}
				author.Affiliation = models.Affiliation{}
				authAffElem := authorElem["afid"]
				authAff := authAffElem.([]interface{})
				if len(authAff) > 0 {
					affID, ok := authAff[0].(map[string]interface{})["$"]
					if ok {
						author.AffiliationID = affID.(string)
					}
				}
				article.Authors = append(article.Authors, author)
			}

			article.Affiliations = []models.Affiliation{}
			affElem := entry["affiliation"]
			affiliations := affElem.([]interface{})
			for _, affVal := range affiliations {
				affElem := affVal.(map[string]interface{})
				affiliation := models.Affiliation{}
				id, ok := affElem["afid"]
				if ok {
					affiliation.ScopusID = id.(string)
				}
				title, ok := affElem["affilname"]
				if ok {
					affiliation.Title = title.(string)
				}
				city, ok := affElem["affiliation-city"]
				if ok {
					affiliation.City = city.(string)
				}
				country, ok := affElem["affiliation-country"]
				if ok {
					affiliation.Country = country.(string)
				}
				article.Affiliations = append(article.Affiliations, affiliation)
			}

			result = append(result, article)
		}
		return result, nil
	}
	return []models.Article{}, errors.New("empty search response")
}

func (worker *Worker) ProceedArticle(article models.Article, articleDs DataSource, depth int) error {
	articleData, err := query.MakeQuery(articleDs.Path, article.ScopusID, map[string]string{}, worker.Config.RequestTimeout, worker.Storage)
	if err != nil {
		logger.Error.Println("Error on requesting data for id=" + article.ScopusID)
		logger.Error.Println(err)
	}
	articleContainer, articleRetrieved := articleData["abstracts-retrieval-response"]
	if articleRetrieved {
		article.Keywords = []models.Keyword{}
		keywordsContainer, ok := articleContainer.(map[string]interface{})["authkeywords"]
		if !ok {
			return errors.New("error on parsing authkeywords element")
		}
		keywordsVal, ok := keywordsContainer.(map[string]interface{})["author-keyword"]
		if !ok {
			return errors.New("error on parsing author-keyword element")
		}
		keywords, ok := keywordsVal.([]interface{})
		if !ok {
			return errors.New("error on parsing keywords element")
		}
		for _, keywordVal := range keywords {
			keyword := models.Keyword{}
			keyword.ID = uuid.NewV4().String()
			keywordElem := keywordVal.(map[string]interface{})
			value, ok := keywordElem["$"]
			if ok {
				keyword.Value = value.(string)
				article.Keywords = append(article.Keywords, keyword)
			}
		}

		article.SubjectAreas = []models.SubjectArea{}
		areaContainer, ok := articleContainer.(map[string]interface{})["subject-areas"]
		if !ok {
			return errors.New("error on parsing subject-areas element")
		}
		areasVal, ok := areaContainer.(map[string]interface{})["subject-area"]
		if !ok {
			return errors.New("error on parsing subject-area element")
		}
		areas, ok := areasVal.([]interface{})
		if !ok {
			return errors.New("error on parsing subject areas element")
		}
		for _, areaVal := range areas {
			area := models.SubjectArea{}
			areaElem := areaVal.(map[string]interface{})
			name, ok := areaElem["$"]
			if ok {
				area.Description = name.(string)
			}
			code, ok := areaElem["@code"]
			if ok {
				area.Code = code.(string)
			}
			abbrev, ok := areaElem["@abbrev"]
			if ok {
				area.Title = abbrev.(string)
			}
			article.SubjectAreas = append(article.SubjectAreas, area)
		}

		article.References = []models.Article{}
		refContainer, ok := articleContainer.(map[string]interface{})["item"]
		if !ok {
			return errors.New("error on parsing item element")
		}
		refContainer, ok = refContainer.(map[string]interface{})["bibrecord"]
		if !ok {
			return errors.New("error on parsing bibrecord element")
		}
		refContainer, ok = refContainer.(map[string]interface{})["tail"]
		if !ok {
			return errors.New("error on parsing tail element")
		}
		refContainer, ok = refContainer.(map[string]interface{})["bibliography"]
		if !ok {
			return errors.New("error on parsing bibliography element")
		}
		refContainer, ok = refContainer.(map[string]interface{})["reference"]
		if !ok {
			return errors.New("error on parsing reference element")
		}
		references, ok := refContainer.([]interface{})
		if !ok {
			return errors.New("error on parsing references element")
		}
		for _, refVal := range references {
			reference := models.Article{}
			refElem := refVal.(map[string]interface{})
			refData, ok := refElem["ref-info"]
			if !ok {
				continue
			}
			refData, ok = refData.(map[string]interface{})["refd-itemidlist"]
			if !ok {
				continue
			}
			refData, ok = refData.(map[string]interface{})["itemid"]
			if !ok {
				continue
			}
			refElem = refData.(map[string]interface{})
			idType, ok := refElem["@idtype"]
			if ok && idType == "SGR" {
				id, ok := refElem["$"]
				if ok {
					reference.ScopusID = id.(string)
					if depth < worker.Config.ReferencesDepth {
						worker.ProceedArticle(reference, articleDs, depth+1)
					}
					article.References = append(article.References, reference)
				}
			}
		}
		for _, reference := range article.References {
			err := worker.Storage.CreateArticle(reference)
			if err != nil {
				return err
			}
		}
		err := worker.Storage.CreateArticle(article)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("empty article response")
}
