using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace ScopusCrawler.Models
{
    public class Affiliation
    {
        public int Id { get; set; }
        public string ScopusID { get; set; }
        public string Name { get; set; }
        public string City { get; set; }
        public string Country { get; set; }
        public List<Author> Authors { get; set; }
    }
}
