package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"spannerfts/importer"
	"spannerfts/search"
)

func main() {
	// Define a command-line flag to trigger the import functionality
	importFlag := flag.Bool("import", false, "Import data from a CSV file")
	filePath := flag.String("file", "", "Path to the CSV file to import")
	flag.Parse()

	// Debugging statements to check flag values
	// fmt.Printf("importFlag: %v\n", *importFlag)
	// fmt.Printf("filePath: %s\n", *filePath)

	if *importFlag {
		if *filePath == "" {
			fmt.Println("Please provide the path to the CSV file using the -file flag")
			os.Exit(1)
		}
		// Call the import function
		if err := importer.ImportData(*filePath); err != nil {
			fmt.Printf("Failed to import data: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Data imported successfully")
		return
	}

	// Serve static files from the "public" directory
	publicFs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", publicFs))

	// Serve static files from the "src" directory
	srcFs := http.FileServer(http.Dir("src"))
	http.Handle("/src/", http.StripPrefix("/src/", srcFs))

	// Serve index.html file
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})

	http.HandleFunc("/search", handleSearch)

	fmt.Println("Server started on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Failed to start server: %v", err)
	}
}

// Handles both GET and POST requests, creates a new search, and returns the search results as JSON
func handleSearch(w http.ResponseWriter, r *http.Request) {
	var searchText string

	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*") // Adjust as per your requirements
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Determine the request method and extract the search query accordingly
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		searchText = r.FormValue("query")
	} else {
		searchText = r.URL.Query().Get("query")
	}

	// Perform the search
	searcher := &search.TextSearch{}
	results, err := searcher.NewSearch(searchText)
	if err != nil {
		http.Error(w, fmt.Sprintf("Search failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Failed to encode search results", http.StatusInternalServerError)
	}
}
