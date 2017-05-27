package crawler

type Manager struct {
	DataSources []DataSource
	Workers     map[string]Worker
}

func (manager manager) Init() {
	
}
