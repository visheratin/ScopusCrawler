using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Xml.Linq;

namespace ScopusCrawler.Models
{
    public class Article
    {
        public int Id { get; set; }
        public string ScopusID { get; set; }
        public string Title { get; set; }
        public List<Author> Authors { get; set; }
        public string PublicationDate { get; set; }
        public string Abstracts { get; set; }
        public int CitationsCount { get; set; }
        public string DOI { get; set; }
        public List<Article> References { get; set; }
        public List<Article> Citations { get; set; }
        public List<string> Keywords { get; set; }
        public List<SubjectArea> SubjectAreas { get; set; }

        public Article()
        {
            Keywords = new List<string>();
            References = new List<Article>();
            Authors = new List<Author>();
        }
    }
}
