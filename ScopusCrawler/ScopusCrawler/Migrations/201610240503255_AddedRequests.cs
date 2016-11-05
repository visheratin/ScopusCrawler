namespace ScopusCrawler.Migrations
{
    using System;
    using System.Data.Entity.Migrations;
    
    public partial class AddedRequests : DbMigration
    {
        public override void Up()
        {
            CreateTable(
                "dbo.RequestDones",
                c => new
                    {
                        Id = c.Int(nullable: false, identity: true),
                        Request = c.String(),
                        Response = c.String(),
                    })
                .PrimaryKey(t => t.Id);
            
        }
        
        public override void Down()
        {
            DropTable("dbo.RequestDones");
        }
    }
}
