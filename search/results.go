package search

type SearchResults struct {
	Rows []sqlRow
}

type sqlRow struct {
	Country    string `json:"country"`
	City       string `json:"city"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	Websites   string `json:"websites"`
	Categories string `json:"categories"`
}
