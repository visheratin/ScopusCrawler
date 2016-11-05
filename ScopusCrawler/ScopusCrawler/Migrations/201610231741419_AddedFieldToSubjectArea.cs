namespace ScopusCrawler.Migrations
{
    using System;
    using System.Data.Entity.Migrations;
    
    public partial class AddedFieldToSubjectArea : DbMigration
    {
        public override void Up()
        {
            AddColumn("dbo.SubjectAreas", "Abbreviation", c => c.String());
        }
        
        public override void Down()
        {
            DropColumn("dbo.SubjectAreas", "Abbreviation");
        }
    }
}
