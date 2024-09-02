package importer

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/spanner"
)

const (
	projectId  = "test"
	instanceId = "test"
	databaseId = "testdb"
)

func ImportData(filePath string) error {
	// 1. Set up Spanner client
	ctx := context.Background()
	client, err := spanner.NewClient(ctx,
		fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectId, instanceId, databaseId))
	if err != nil {
		return fmt.Errorf("error creating Spanner client: %v", err)

	}
	defer client.Close()

	// 2. Open CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Error opening CSV file: %v", err)
	}
	defer file.Close()

	// 3. Read and parse the CSV file
	reader := csv.NewReader(file)
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("error reading header row: %v", err)
	}
	columnNames := header

	batchSize := 500
	mutations := []*spanner.Mutation{}

	// 3.1. Read and process data rows
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("Error reading row: %v\n", err)
			continue // Skip this row and continue with the next
		}
		// Convert row from []string to []interface{}
		rowInterfaces := make([]interface{}, len(row))
		for i, v := range row {
			switch i {
			case 1: // TIMESTAMP
				addTimestamp, convErr := time.Parse(time.RFC3339, v)
				if convErr != nil {
					fmt.Printf("Error converting timestamp to TIMESTAMP: %v\n", convErr)
					continue // Skip this row and continue with the next
				}
				rowInterfaces[i] = addTimestamp
			case 2: // TIMESTAMP
				updateTimestamp, convErr := time.Parse(time.RFC3339, v)
				if convErr != nil {
					fmt.Printf("Error converting timestamp to TIMESTAMP: %v\n", convErr)
					continue // Skip this row and continue with the next
				}
				rowInterfaces[i] = updateTimestamp
			case 9: // latitude column, FLOAT64
				latitudeFloat, convErr := strconv.ParseFloat(v, 64)
				if convErr != nil {
					fmt.Printf("Error converting rating to FLOAT64: %v\n", convErr)
					continue // Skip this row and continue with the next
				}
				rowInterfaces[i] = latitudeFloat
			case 10: // longitude column, FLOAT64
				longitudeFloat, convErr := strconv.ParseFloat(v, 64)
				if convErr != nil {
					fmt.Printf("Error converting rating to FLOAT64: %v\n", convErr)
					continue // Skip this row and continue with the next
				}
				rowInterfaces[i] = longitudeFloat
			default: // For STRING columns, no conversion needed
				rowInterfaces[i] = v
			}
		}

		// Assuming your Spanner table name is "my_table"
		mutations = append(mutations, spanner.Insert("Restaurants", columnNames, rowInterfaces))
		if len(mutations) >= batchSize {
			// Apply the current batch of mutations
			_, err := client.Apply(ctx, mutations)
			if err != nil {
				return fmt.Errorf("error applying mutations: %v", err)
			}
			// Reset mutations for the next batch
			mutations = []*spanner.Mutation{}
		}
	}

	// 3.2. Apply any remaining mutations that didn't reach the batchSize threshold
	if len(mutations) > 0 {
		_, err := client.Apply(ctx, mutations)
		if err != nil {
			return fmt.Errorf("error applying mutations: %v", err)
		}
	}
	fmt.Println("CSV data imported to Spanner successfully!")
	return nil
}
