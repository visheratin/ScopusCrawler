package crawler

type DataSource struct {
	Name string
	Path string
	Keys []string
}

type SearchRequest struct {
	SourceName string
	Source     DataSource
	ID         string
	Fields     map[string]string
}
