namespace ScopusCrawler.Migrations
{
    using System;
    using System.Data.Entity.Migrations;
    
    public partial class Init : DbMigration
    {
        public override void Up()
        {
            CreateTable(
                "dbo.Affiliations",
                c => new
                    {
                        Id = c.Int(nullable: false, identity: true),
                        ScopusID = c.String(),
                        Name = c.String(),
                        City = c.String(),
                        Country = c.String(),
                    })
                .PrimaryKey(t => t.Id);
            
            CreateTable(
                "dbo.Authors",
                c => new
                    {
                        Id = c.Int(nullable: false, identity: true),
                        ScopusID = c.String(),
                        FullName = c.String(),
                        Initials = c.String(),
                        Surname = c.String(),
                        AffiliationId = c.Int(nullable: false),
                    })
                .PrimaryKey(t => t.Id)
                .ForeignKey("dbo.Affiliations", t => t.AffiliationId, cascadeDelete: true)
                .Index(t => t.AffiliationId);
            
            CreateTable(
                "dbo.Articles",
                c => new
                    {
                        Id = c.Int(nullable: false, identity: true),
                        ScopusID = c.String(),
                        Title = c.String(),
                        PublicationDate = c.String(),
                        Abstracts = c.String(),
                        CitationsCount = c.Int(nullable: false),
                        DOI = c.String(),
                    })
                .PrimaryKey(t => t.Id);
            
            CreateTable(
                "dbo.SubjectAreas",
                c => new
                    {
                        Id = c.Int(nullable: false, identity: true),
                        ScopusID = c.String(),
                        Name = c.String(),
                    })
                .PrimaryKey(t => t.Id);
            
            CreateTable(
                "dbo.ArticleAuthors",
                c => new
                    {
                        Article_Id = c.Int(nullable: false),
                        Author_Id = c.Int(nullable: false),
                    })
                .PrimaryKey(t => new { t.Article_Id, t.Author_Id })
                .ForeignKey("dbo.Articles", t => t.Article_Id, cascadeDelete: true)
                .ForeignKey("dbo.Authors", t => t.Author_Id, cascadeDelete: true)
                .Index(t => t.Article_Id)
                .Index(t => t.Author_Id);
            
            CreateTable(
                "dbo.ArticleArticles",
                c => new
                    {
                        Article_Id = c.Int(nullable: false),
                        Article_Id1 = c.Int(nullable: false),
                    })
                .PrimaryKey(t => new { t.Article_Id, t.Article_Id1 })
                .ForeignKey("dbo.Articles", t => t.Article_Id)
                .ForeignKey("dbo.Articles", t => t.Article_Id1)
                .Index(t => t.Article_Id)
                .Index(t => t.Article_Id1);
            
            CreateTable(
                "dbo.SubjectAreaArticles",
                c => new
                    {
                        SubjectArea_Id = c.Int(nullable: false),
                        Article_Id = c.Int(nullable: false),
                    })
                .PrimaryKey(t => new { t.SubjectArea_Id, t.Article_Id })
                .ForeignKey("dbo.SubjectAreas", t => t.SubjectArea_Id, cascadeDelete: true)
                .ForeignKey("dbo.Articles", t => t.Article_Id, cascadeDelete: true)
                .Index(t => t.SubjectArea_Id)
                .Index(t => t.Article_Id);
            
        }
        
        public override void Down()
        {
            DropForeignKey("dbo.SubjectAreaArticles", "Article_Id", "dbo.Articles");
            DropForeignKey("dbo.SubjectAreaArticles", "SubjectArea_Id", "dbo.SubjectAreas");
            DropForeignKey("dbo.ArticleArticles", "Article_Id1", "dbo.Articles");
            DropForeignKey("dbo.ArticleArticles", "Article_Id", "dbo.Articles");
            DropForeignKey("dbo.ArticleAuthors", "Author_Id", "dbo.Authors");
            DropForeignKey("dbo.ArticleAuthors", "Article_Id", "dbo.Articles");
            DropForeignKey("dbo.Authors", "AffiliationId", "dbo.Affiliations");
            DropIndex("dbo.SubjectAreaArticles", new[] { "Article_Id" });
            DropIndex("dbo.SubjectAreaArticles", new[] { "SubjectArea_Id" });
            DropIndex("dbo.ArticleArticles", new[] { "Article_Id1" });
            DropIndex("dbo.ArticleArticles", new[] { "Article_Id" });
            DropIndex("dbo.ArticleAuthors", new[] { "Author_Id" });
            DropIndex("dbo.ArticleAuthors", new[] { "Article_Id" });
            DropIndex("dbo.Authors", new[] { "AffiliationId" });
            DropTable("dbo.SubjectAreaArticles");
            DropTable("dbo.ArticleArticles");
            DropTable("dbo.ArticleAuthors");
            DropTable("dbo.SubjectAreas");
            DropTable("dbo.Articles");
            DropTable("dbo.Authors");
            DropTable("dbo.Affiliations");
        }
    }
}
