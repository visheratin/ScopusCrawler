using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Configuration;
using System.Net;
using System.IO;
using System.Xml.Linq;
using System.Collections.Concurrent;
using System.Text.RegularExpressions;
using ScopusCrawler.Processing;

namespace ScopusCrawler
{
    class Program
    {
        static string ApiUrl = "http://api.elsevier.com/content/search/scopus?sort=citedby-count&httpAccept=application/xml&";
        static int ResultsPerPage = 25;
        
        static void Main(string[] args)
        {
            var manager = new Manager();
            manager.StartCrawling();
            //var titles = new List<string>();
            //using (var stream = new FileStream("output.tsv", FileMode.Open))
            //using (var reader = new StreamReader(stream))
            //{
            //    var line = string.Empty;
            //    while ((line = reader.ReadLine()) != null)
            //    {
            //        if (!titles.Contains(line))
            //        {
            //            titles.Add(line);
            //        }
            //    }
            //}
            //using (var stream = new FileStream("output.tsv", FileMode.Create))
            //using (var writer = new StreamWriter(stream))
            //{
            //    writer.WriteLine("author\ttitle\tabstracts\tcitedby\tpublication_date\tcitations\tkeywords");
            //    foreach (var line in titles)
            //    {
            //        writer.WriteLine(line);
            //    }
            //}
            //return;

            //var files = Directory.GetFiles("dump");
            //var lines = new List<string>();
            //foreach (var file in files)
            //{
            //    using (var stream = new FileStream(file, FileMode.Open))
            //    using (var reader = new StreamReader(stream))
            //    {
            //        var line = string.Empty;
            //        while ((line = reader.ReadLine()) != null)
            //        {
            //            line = line.Replace("\"", "").Replace("\n", "");
            //            lines.Add(line);
            //        }
            //    }
            //}
            //using (var stream = new FileStream("output.tsv", FileMode.Create))
            //using (var writer = new StreamWriter(stream))
            //{
            //    writer.WriteLine("author\ttitle\tabstracts\tcitedby\tpublication_date\tcitations\tkeywords");
            //    foreach (var line in lines)
            //    {
            //        writer.WriteLine(line);
            //    }
            //}
            //return;
            //var journalsISSN = new List<string>();
            //using (var stream = new FileStream("data//scopus-journals.csv", FileMode.Open))
            //using (var reader = new StreamReader(stream))
            //{
            //    var regex = new Regex("\\d{8}");
            //    var line = string.Empty;
            //    while((line = reader.ReadLine()) != null)
            //    {
            //        var match = regex.Match(line);
            //        while (match.Success)
            //        {
            //            journalsISSN.Add(match.Groups[0].Value);
            //            match = match.NextMatch();
            //        }
            //    }
            //}
            //if (!Directory.Exists("dump"))
            //    Directory.CreateDirectory("dump");
            //var apiKeys = GetApiKeys();
            //var startYear = 1990;
            //var finishYear = 2016;
            //foreach (var issn in journalsISSN)
            //{
            //    if (!Directory.Exists("dump\\" + issn))
            //        Directory.CreateDirectory("dump" + issn);
            //    var query = "query=ISSN(" + issn + ")";
            //    for (var year = finishYear; year >= startYear; year--)
            //    {
            //        var initRequest = ApiUrl + "&" + query + "&apiKey=" + apiKeys[0] + "&date=" + year;
            //        var initData = GetData(initRequest);
            //        var doc = XElement.Parse(initData);
            //        var totalResultsNumber = int.Parse(doc.Elements().First().Value);
            //        if (totalResultsNumber > 4975)
            //            totalResultsNumber = 4975;
            //        var splitsCount = totalResultsNumber / ResultsPerPage;
            //        if (totalResultsNumber % ResultsPerPage != 0)
            //        {
            //            splitsCount++;
            //        }
            //        var splitsPerThread = splitsCount / apiKeys.Count;
            //        var lastSplit = splitsPerThread + splitsCount % apiKeys.Count;
            //        var tasks = new Task[apiKeys.Count];
            //        for (int k = 0; k < apiKeys.Count; k++)
            //        {
            //            var counter = k;
            //            var task = new Task(() =>
            //            {
            //                var currentSplitsCount = counter < apiKeys.Count - 1 ? splitsPerThread : lastSplit;
            //                Console.WriteLine("Thread " + counter + " was started. Total number of splits: " + currentSplitsCount);
            //                for (int i = 0; i < currentSplitsCount; i++)
            //                {
            //                    var articles = new List<Article>();
            //                    var request = ApiUrl + "&" + query + "&apiKey=" + apiKeys[counter] + "&date=" + year + "&start=" + ((splitsPerThread * counter + i) * 25);
            //                    var data = GetData(request);
            //                    if (data != string.Empty)
            //                    {
            //                        var element = XElement.Parse(data);
            //                        var articleElements = element.Elements(XName.Get("{http://www.w3.org/2005/Atom}entry"));
            //                        foreach (var articleElement in articleElements)
            //                        {
            //                            var article = ParseArticle(articleElement);
            //                            if (article != null)
            //                            {
            //                                article.Citations = GetCitations(article.ReferenceID, apiKeys[counter], article.CitationsCount);
            //                                GetKeywordsAndAbstracts(article.Url, apiKeys[0], ref article);
            //                                articles.Add(article);
            //                            }
            //                        }

            //                        using (var stream = new FileStream("dump\\" + year.ToString() + "-" + counter.ToString() + "-" + i.ToString() + ".tsv", FileMode.Create))
            //                        using (var writer = new StreamWriter(stream))
            //                        {
            //                            foreach (var article in articles)
            //                            {
            //                                var stringBuilder = new StringBuilder();
            //                                stringBuilder.Append(article.Author).Append("\t").
            //                                    Append(article.Title).Append("\t").
            //                                    Append(article.Abstracts).Append("\t").
            //                                    Append(article.CitationsCount).Append("\t").
            //                                    Append(article.PublicationDate).Append("\t");
            //                                foreach (var citation in article.Citations)
            //                                {
            //                                    stringBuilder.Append(citation).Append(";");
            //                                }
            //                                if (article.Citations.Count == 0)
            //                                    stringBuilder.Append("none");
            //                                stringBuilder.Append("\t");
            //                                foreach (var keyword in article.Keywords)
            //                                {
            //                                    stringBuilder.Append(keyword).Append(";");
            //                                }
            //                                writer.WriteLine(stringBuilder.ToString());
            //                            }
            //                        }

            //                        Console.WriteLine("Request " + i + " at thread " + counter + " finished.");
            //                    }
            //                }
            //            });
            //            tasks[k] = task;
            //            task.Start();
            //        }
            //        Task.WaitAll(tasks);
            //    }
            //}
        }
    }
}
