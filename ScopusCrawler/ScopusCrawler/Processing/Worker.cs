using ScopusCrawler.Models;
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Net;
using System.Threading;
using System.Xml.Linq;
using System.Collections.Concurrent;

namespace ScopusCrawler.Processing
{
    public class Worker
    {
        string ApiUrl = "http://api.elsevier.com/content/search/scopus?sort=citedby-count&httpAccept=application/xml&";
        int ResultsPerPage = 25;
        BlockingCollection<string> brokenArticles;
        BlockingCollection<string> processedArticles;

        public Worker(BlockingCollection<string> articlesProcessed, BlockingCollection<string> articlesBroken)
        {
            processedArticles = articlesProcessed;
            brokenArticles = articlesBroken;
        }

        public List<Article> Start(int year, string journalISSN)
        {
            var query = "query=ISSN(" + journalISSN + ")";
            var initRequest = ApiUrl + "&" + query + "&date=" + year;
            var initData = GetData(initRequest);
            var doc = XElement.Parse(initData);
            var totalResultsNumber = int.Parse(doc.Elements().First().Value);
            if (totalResultsNumber > 4975)
                totalResultsNumber = 4975;
            var splitsCount = totalResultsNumber / ResultsPerPage;
            if (totalResultsNumber % ResultsPerPage != 0)
            {
                splitsCount++;
            }
            var result = new List<Article>();
            for (int i = 0; i < splitsCount; i++)
            {
                var request = ApiUrl + "&" + query + "&date=" + year + "&start=" + (i * 25).ToString();
                var data = GetData(request);
                if (data != string.Empty)
                {
                    var element = XElement.Parse(data);
                    var articleElements = element.Elements(XName.Get("{http://www.w3.org/2005/Atom}entry"));
                    var articles = CollectArticles(articleElements);
                    result.AddRange(articles);
                    //foreach (var articleElement in articleElements)
                    //{
                    //    try
                    //    {
                    //        var article = new Article();
                    //        var id = articleElement.Element(XName.Get("{http://www.w3.org/2005/Atom}identifier")).Value;
                    //        article.ScopusID = id.Replace("SCOPUS_ID:", "");
                    //        GetArticleData(ref article);
                    //    }
                    //    catch (Exception ex)
                    //    {
                    //        Console.WriteLine(ex.Message);
                    //    }
                    //}
                }
            }
            return result;
        }

        List<Article> CollectArticles(IEnumerable<XElement> articleElements)
        {
            var result = new List<Article>();
            foreach (var articleElement in articleElements)
            {
                try
                {
                    var article = new Article();
                    var id = articleElement.Element(XName.Get("{http://www.w3.org/2005/Atom}identifier")).Value;
                    article.ScopusID = id.Replace("SCOPUS_ID:", "");
                    GetArticleData(ref article);
                    result.Add(article);
                }
                catch (Exception ex)
                {
                    Console.WriteLine(ex.Message);
                }
            }
            return result;
        }

        void GetArticleData(ref Article article)
        {
            //using (ScopusDbContext dbContext = new ScopusDbContext())
            {
                {
                    var articleId = article.ScopusID;
                    if (processedArticles.Contains(articleId))
                    {
                        return;
                    }
                    Console.WriteLine("Working on article " + article.ScopusID);
                    var apiUrl = "http://api.elsevier.com/content/abstract/scopus_id/{0}?httpAccept=application/xml";
                    var request = string.Format(apiUrl, article.ScopusID);
                    var data = GetData(request);
                    if (data != string.Empty)
                    {
                        var element = XElement.Parse(data);
                        var coreElement = element.Element(XName.Get("{http://www.elsevier.com/xml/svapi/abstract/dtd}coredata"));
                        try
                        {
                            article.Title = coreElement.Element(XName.Get("{http://www.elsevier.com/xml/svapi/abstract/dtd}title")).Value;
                            article.PublicationDate = coreElement.Element(XName.Get("{http://www.elsevier.com/xml/svapi/abstract/dtd}coverDate")).Value;
                        }
                        catch
                        {
                            brokenArticles.Add(article.ScopusID);
                            return;
                        }
                        article.CitationsCount = int.Parse(coreElement.Element(XName.Get("{http://www.elsevier.com/xml/svapi/abstract/dtd}citedby-count")).Value);
                        var descriptionElement = coreElement.Elements(XName.Get("{http://www.elsevier.com/xml/svapi/abstract/dtd}description"));
                        var abstractsElement = descriptionElement.Elements(XName.Get("{http://www.elsevier.com/xml/svapi/abstract/dtd}abstract"));
                        if (abstractsElement.Count() > 0)
                        {
                            article.Abstracts = abstractsElement.First().Value;
                        }
                        var affiliationsElement = element.Elements(XName.Get("{http://www.elsevier.com/xml/svapi/abstract/dtd}affiliation"));
                        var affiliations = GetAffiliations(affiliationsElement);
                        var authorsElement = element.Element(XName.Get("{http://www.elsevier.com/xml/svapi/abstract/dtd}authors"));
                        article.Authors = GetAuthors(authorsElement, affiliations);
                        var doiElement = coreElement.Elements(XName.Get("{http://www.elsevier.com/xml/svapi/abstract/dtd}doi"));
                        if (doiElement.Count() > 0)
                        {
                            article.DOI = doiElement.First().Value;
                        }
                        var keywordsElement = element.Elements(XName.Get("{http://www.elsevier.com/xml/svapi/abstract/dtd}authkeywords"));
                        foreach (var keyword in keywordsElement.Elements())
                        {
                            article.Keywords.Add(keyword.Value);
                        }
                        var recordsItem = element.Element(XName.Get("item"));
                        var bibrecordItem = recordsItem.Element(XName.Get("bibrecord"));
                        var tailItem = bibrecordItem.Element(XName.Get("tail"));
                        var subjectAreasElement = element.Element(XName.Get("{http://www.elsevier.com/xml/svapi/abstract/dtd}subject-areas"));
                        article.SubjectAreas = GetSubjectAreas(subjectAreasElement);
                        var referencesItem = tailItem.Element(XName.Get("bibliography"));
                        if (referencesItem != null)
                        {
                            article.References = GetReferences(referencesItem.Elements());
                        }
                        //if (!dbContext.Articles.Any(a => a.ScopusID == articleId))
                        //    dbContext.Articles.Add(article);
                        //dbContext.SaveChanges();
                        //Console.WriteLine("Added article: " + article.Title);
                        //processedArticles.Add(article.ScopusID);
                    }
                }
            }
        }

        private List<SubjectArea> GetSubjectAreas(XElement subjectAreasElement)
        {
            var result = new List<SubjectArea>();
            //using (ScopusDbContext dbContext = new ScopusDbContext())
            {
                foreach (var element in subjectAreasElement.Elements())
                {
                    var area = new SubjectArea()
                    {
                        ScopusID = element.Attribute(XName.Get("code")).Value,
                        Abbreviation = element.Attribute(XName.Get("abbrev")).Value,
                        Name = element.Value
                    };
                    //var subjectArea = dbContext.SubjectAreas.FirstOrDefault(a => a.ScopusID == area.ScopusID);
                    //if (subjectArea == null)
                    //{
                    //    dbContext.SubjectAreas.Add(area);
                    //    dbContext.SaveChanges();
                    //    subjectArea = area;
                    //    Console.WriteLine("Added subject area: " + area.Name);
                    //}
                    //result.Add(subjectArea);
                    result.Add(area);
                }
            }
            return result;
        }

        private List<Article> GetReferences(IEnumerable<XElement> referencesElements)
        {
            var result = new List<Article>();
            //using (ScopusDbContext dbContext = new ScopusDbContext())
            {
                foreach (var element in referencesElements)
                {
                    var refInfoElement = element.Element(XName.Get("ref-info"));
                    var refIdListElement = refInfoElement.Element(XName.Get("refd-itemidlist"));
                    foreach (var refIdElement in refIdListElement.Elements())
                    {
                        var refType = refIdElement.Attribute(XName.Get("idtype")).Value;
                        if (refType == "SGR")
                        {
                            var refId = refIdElement.Value;
                            var refArticle = new Article();
                            refArticle.ScopusID = refId;
                            //if (!brokenArticles.Contains(refId) && !processedArticles.Contains(refId))
                            //{
                            //    if (!dbContext.Articles.Any(a => a.ScopusID == refId))
                            //    {
                            //        var refArticle = new Article();
                            //        refArticle.ScopusID = refId;
                            //        dbContext.Articles.Add(refArticle);
                            //        dbContext.SaveChanges();
                            //        //GetArticleData(ref refArticle);
                            //    }
                            //    try
                            //    {
                            //        result.Add(dbContext.Articles.First(a => a.ScopusID == refId));
                            //    }
                            //    catch
                            //    {
                            //        Console.WriteLine("--- Broken reference! ---");
                            //    }
                            //}
                            result.Add(refArticle);
                        }
                    }
                }
            }
            return result;
        }

        private List<Affiliation> GetAffiliations(IEnumerable<XElement> affiliationElements)
        {
            var result = new List<Affiliation>();
            //using (ScopusDbContext dbContext = new ScopusDbContext())
            {
                foreach (var element in affiliationElements)
                {
                    var scopusId = element.Attribute(XName.Get("id")).Value;
                    //var affiliation = dbContext.Affiliations.FirstOrDefault(a => a.ScopusID == scopusId);
                    //if (affiliation == null)
                    //{
                        var affiliation = new Affiliation()
                        {
                            ScopusID = scopusId,
                            Name = element.Element(XName.Get("{http://www.elsevier.com/xml/svapi/abstract/dtd}affilname")).Value,
                            City = element.Element(XName.Get("{http://www.elsevier.com/xml/svapi/abstract/dtd}affiliation-city")).Value,
                            Country = element.Element(XName.Get("{http://www.elsevier.com/xml/svapi/abstract/dtd}affiliation-country")).Value
                        };
                        //dbContext.Affiliations.Add(affiliation);
                        //dbContext.SaveChanges();
                    //    Console.WriteLine("Added affiliation: " + affiliation.Name);
                    //}
                    result.Add(affiliation);
                }
            }
            return result;
        }

        private List<Author> GetAuthors(XElement authorsElement, List<Affiliation> affiliations)
        {
            var result = new List<Author>();
            //using (ScopusDbContext dbContext = new ScopusDbContext())
            {
                foreach (var element in authorsElement.Elements())
                {
                    var scopusId = element.Attribute(XName.Get("auid")).Value;
                    //var author = dbContext.Authors.FirstOrDefault(a => a.ScopusID == scopusId);
                    //if (author == null)
                    //{
                        var author = new Author();
                        author.ScopusID = scopusId;
                        var fullNameElement = element.Element(XName.Get("{http://www.elsevier.com/xml/svapi/abstract/dtd}indexed-name"));
                        if (fullNameElement != null)
                            author.FullName = fullNameElement.Value;
                        var initialsElement = element.Element(XName.Get("{http://www.elsevier.com/xml/svapi/abstract/dtd}initials"));
                        if (initialsElement != null)
                            author.Initials = initialsElement.Value;
                        var surnameElement = element.Element(XName.Get("{http://www.elsevier.com/xml/svapi/abstract/dtd}surname"));
                        if (surnameElement != null)
                            author.Surname = surnameElement.Value;
                        var affiliationElement = element.Element(XName.Get("{http://www.elsevier.com/xml/svapi/abstract/dtd}affiliation"));
                        if (affiliationElement != null)
                        {
                            var affiliationId = affiliationElement.Attribute(XName.Get("id")).Value;
                            var affiliation = affiliations.FirstOrDefault(a => a.ScopusID == affiliationId);
                            if (affiliation != null)
                                author.Affiliation = affiliations.First(a => a.ScopusID == affiliationId);
                            else
                                author.AffiliationId = 240;
                        }
                        else
                        {
                            author.AffiliationId = 240;
                        }
                        //dbContext.Authors.Add(author);
                        //dbContext.SaveChanges();
                        //Console.WriteLine("Added author: " + author.FullName);
                    //}
                    result.Add(author);
                }
            }
            return result;
        }

        string GetData(string url)
        {
            try
            {
                using (ScopusDbContext dbContext = new ScopusDbContext())
                {
                    var requestDone = dbContext.Requests.FirstOrDefault(r => r.Request == url);
                    if (requestDone == null || requestDone.Response == string.Empty)
                    {
                        var key = KeysStorage.GetKey();
                        var requestUrl = url + "&apiKey=" + key;
                        var request = (HttpWebRequest)WebRequest.Create(requestUrl);
                        request.Proxy = new WebProxy("proxy.ifmo.ru", 3128);
                        var response = request.GetResponse();
                        Stream dataStream = response.GetResponseStream();
                        StreamReader reader = new StreamReader(dataStream);
                        string responseString = reader.ReadToEnd();
                        responseString = responseString.Replace("prism:", "").Replace("dc:", "").Replace("opensearch:", "").Replace("ce:", "").Replace("atom:", "");
                        dbContext.Requests.Add(new RequestDone()
                        {
                            Request = url,
                            Response = responseString
                        });
                        dbContext.ChangeTracker.DetectChanges();
                        dbContext.SaveChanges();
                        dbContext.Dispose();
                        //Thread.Sleep(1000);
                        return responseString;
                    }
                    else
                    {
                        return requestDone.Response;
                    }
                }
            }
            catch (Exception ex)
            {
                Console.WriteLine("Request failed.");
                return string.Empty;
            }
        }
    }
}
