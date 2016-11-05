using ScopusCrawler.Models;
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;
using System.Text.RegularExpressions;
using System.Threading;
using System.Threading.Tasks;

namespace ScopusCrawler.Processing
{
    public class Manager
    {
        public void StartCrawling()
        {
            var journalsISSN = new List<string>();
            using (var stream = new FileStream("data//scopus-journals.csv", FileMode.Open))
            using (var reader = new StreamReader(stream))
            {
                var regex = new Regex("\\d{8}");
                var line = string.Empty;
                while ((line = reader.ReadLine()) != null)
                {
                    var match = regex.Match(line);
                    if (match.Success)
                    {
                        journalsISSN.Add(match.Groups[0].Value);
                        //match = match.NextMatch();
                    }
                }
            }
            var startYear = 1990;
            var finishYear = 2016;
            var workersCount = 10;
            var tasks = new Task[workersCount];
            System.Collections.Concurrent.BlockingCollection<string> processedArticles = new System.Collections.Concurrent.BlockingCollection<string>();
            System.Collections.Concurrent.BlockingCollection<string> brokenArticles = new System.Collections.Concurrent.BlockingCollection<string>();
            using (ScopusDbContext dbContext = new ScopusDbContext())
            {
                foreach (var article in dbContext.Articles.AsNoTracking())
                {
                    processedArticles.Add(article.ScopusID);
                }
            }
            //for (int i = 0; i < workersCount; i++)
            {
                //tasks[i] = Task.Factory.StartNew(() =>
                //{
                var counter = 0;
                var processedEntries = new System.Collections.Concurrent.BlockingCollection<Tuple<string, int>>();
                while (journalsISSN.Count > 0)
                {
                    var issn = journalsISSN.First();
                    journalsISSN.RemoveAt(0);
                    for (int year = finishYear; year >= startYear; year--)
                    {
                        var toStart = false;
                        using (ScopusDbContext dbContext = new ScopusDbContext())
                        {
                            if (!dbContext.EntriesDone.Any(e => e.Issn == issn && e.Year == year))
                                toStart = true;
                        }
                        if (toStart)
                        {
                            var startParameters = new Tuple<string, int>(issn, year);
                            tasks[counter] = Task.Factory.StartNew((parameters) =>
                            {
                                var p = (Tuple<string, int>)parameters;
                                var worker = new Worker(processedArticles, brokenArticles);
                                var result = worker.Start(p.Item2, p.Item1);
                                processedEntries.Add(p);
                                return result;
                            }, startParameters);
                            counter++;
                        }
                        if(counter == workersCount)
                        {
                            Task.WaitAll(tasks);
                            for (int i = 0; i < tasks.Length; i++)
                            {
                                var task = (Task<List<Article>>)tasks[i];
                                if(task.Result != null)
                                    UploadArticles(task.Result);
                                using (ScopusDbContext dbContext = new ScopusDbContext())
                                {
                                    var entry = new EntriesDone()
                                    {
                                        Issn = processedEntries.ElementAt(i).Item1,
                                        Year = processedEntries.ElementAt(i).Item2
                                    };
                                    dbContext.EntriesDone.Add(entry);
                                    dbContext.SaveChanges();
                                    Console.WriteLine("---" + processedEntries.ElementAt(i).Item1 + " - " +
                                        processedEntries.ElementAt(i).Item2.ToString() + "---");
                                }
                            }
                            counter = 0;
                        }
                    }
                }
                //});
            }
            
        }

        public void UploadArticles(List<Article> articles)
        {
            using (var dbContext = new ScopusDbContext())
            {
                foreach (var article in articles)
                {
                    if (!dbContext.Articles.Any(a => a.ScopusID == article.ScopusID))
                    {
                        for (int i = 0; i < article.Authors.Count; i++)
                        {
                            var scopusID = article.Authors[i].ScopusID;
                            var dbAuthor = dbContext.Authors.FirstOrDefault(a => a.ScopusID == scopusID);
                            if (dbAuthor != null)
                                article.Authors[i] = dbAuthor;
                            else
                            {
                                if (article.Authors[i].Affiliation != null)
                                {
                                    scopusID = article.Authors[i].Affiliation.ScopusID;
                                    var dbAffiliation = dbContext.Affiliations.FirstOrDefault(a => a.ScopusID == scopusID);
                                    if (dbAffiliation != null)
                                        article.Authors[i].Affiliation = dbAffiliation;
                                }
                            }
                        }
                        for (int i = 0; i < article.SubjectAreas.Count; i++)
                        {
                            var scopusID = article.SubjectAreas[i].ScopusID;
                            var dbArea = dbContext.SubjectAreas.FirstOrDefault(a => a.ScopusID == scopusID);
                            if (dbArea != null)
                            {
                                article.SubjectAreas[i] = dbArea;
                            }
                        }
                        for (int i = 0; i < article.References.Count; i++)
                        {
                            var scopusID = article.References[i].ScopusID;
                            var dbReference = dbContext.Articles.FirstOrDefault(a => a.ScopusID == scopusID);
                            if (dbReference != null)
                            {
                                article.References[i] = dbReference;
                            }
                        }
                        dbContext.Articles.Add(article);
                        dbContext.SaveChanges();
                        Console.WriteLine("Added article: " + article.Title);
                    }
                }
            }
        }
    }
}
