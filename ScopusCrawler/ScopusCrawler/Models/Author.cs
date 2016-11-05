using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace ScopusCrawler.Models
{
    public class Author
    {
        public int Id { get; set; }
        public string ScopusID { get; set; }
        public string FullName { get; set; }
        public string Initials { get; set; }
        public string Surname { get; set; }
        public int AffiliationId { get; set; }
        public Affiliation Affiliation { get; set; }
        public List<Article> Articles { get; set; }
    }
}
