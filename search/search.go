package search

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/googleapis/go-sql-spanner"
)

const fullTextSearchQuery = `
select 
  country, city, name, address, websites,categories
from Restaurants
  where search_substring(categories_token, '%s') OR search(name_token, '%s')
  ORDER BY SCORE(name_token, '%s') DESC;`

const (
	projectId  = "test"
	instanceId = "test"
	databaseId = "testdb"
)

type TextSearch struct{}

func (s *TextSearch) NewSearch(searchText string) (*SearchResults, error) {

	ctx := context.Background()
	db, err := sql.Open("spanner", fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectId, instanceId, databaseId))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %s", err)
	}
	defer db.Close()

	query := fmt.Sprintf(fullTextSearchQuery, searchText, searchText, searchText)
	fmt.Printf("Executing query: %s\n", query) // Debug feature to print the query

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	return s.formatSearchResults(rows)
}

// Format the search results
func (s *TextSearch) formatSearchResults(rows *sql.Rows) (*SearchResults, error) {
	defer rows.Close()

	// Format the returned rows from the query
	var rowData []sqlRow
	for rows.Next() {
		var r sqlRow
		if err := rows.Scan(&r.Country, &r.City, &r.Name, &r.Address, &r.Websites, &r.Categories); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		rowData = append(rowData, r)

		// Print the search result
		fmt.Printf("Country: %s, City: %s, Name: %s, Address: %s, Websites: %s, Categories: %s\n",
			r.Country, r.City, r.Name, r.Address, r.Websites, r.Categories)
	}

	return &SearchResults{rowData}, nil
}
