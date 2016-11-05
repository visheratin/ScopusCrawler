namespace ScopusCrawler.Migrations
{
    using System;
    using System.Data.Entity.Migrations;
    
    public partial class AddedEntries : DbMigration
    {
        public override void Up()
        {
            CreateTable(
                "dbo.EntriesDones",
                c => new
                    {
                        Id = c.Int(nullable: false, identity: true),
                        Issn = c.String(),
                        Year = c.Int(nullable: false),
                    })
                .PrimaryKey(t => t.Id);
            
        }
        
        public override void Down()
        {
            DropTable("dbo.EntriesDones");
        }
    }
}
