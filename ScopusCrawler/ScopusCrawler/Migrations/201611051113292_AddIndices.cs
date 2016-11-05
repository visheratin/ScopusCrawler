namespace ScopusCrawler.Migrations
{
    using System;
    using System.Data.Entity.Migrations;
    
    public partial class AddIndices : DbMigration
    {
        public override void Up()
        {
            CreateIndex("Articles", "ScopusID");
            CreateIndex("Authors", "ScopusID");
            CreateIndex("Affiliations", "ScopusID");
            CreateIndex("SubjectAreas", "ScopusID");
        }
        
        public override void Down()
        {
        }
    }
}
