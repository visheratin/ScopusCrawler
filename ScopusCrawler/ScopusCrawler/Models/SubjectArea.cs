using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace ScopusCrawler.Models
{
    public class SubjectArea
    {
        public int Id { get; set; }
        public string ScopusID { get; set; }
        public string Abbreviation { get; set; }
        public string Name { get; set; }
        public List<Article> Articles { get; set; }
    }
}
