using ScopusCrawler.Models;
using System;
using System.Collections.Generic;
using System.Data.Entity;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace ScopusCrawler
{
    public class ScopusDbContext : DbContext
    {
        public ScopusDbContext() : base("DefaultConnection")
        {
        }

        public virtual DbSet<Article> Articles { get; set; }
        public virtual DbSet<Affiliation> Affiliations { get; set; }
        public virtual DbSet<Author> Authors { get; set; }
        public virtual DbSet<SubjectArea> SubjectAreas { get; set; }
        public virtual DbSet<RequestDone> Requests { get; set; }
        public virtual DbSet<EntriesDone> EntriesDone { get; set; }
    }
}
