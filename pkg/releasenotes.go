package releasenotes

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

// GetReleaseNotes retrieves release notes for a specific product.
//
// The function returns a slice of ReleaseNote structs containing the release
// note type and description, or an error if any occurs during the process.
func GetReleaseNotes(ctx context.Context, projectID string, product string, releaseNoteType string) ([]ReleaseNote, error) {

	// Create a BigQuery client to interact with the BigQuery service.
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	defer client.Close() // Close the client when the function exits.

	// Define the BigQuery query to retrieve release notes for the specified product and specific release note
	q := client.Query(`
	SELECT
		product_name,
		release_note_type,
		product_version_name,
		description,
		published_at
	FROM bigquery-public-data.google_cloud_release_notes.release_notes
	WHERE 
		product_name = @product AND 
		release_note_type = @releaseNoteType
	GROUP BY product_name, product_version_name, release_note_type, description, published_at
	ORDER BY published_at DESC
	LIMIT 10;
		`)

	// Set the query parameters for the product name.
	q.Parameters = []bigquery.QueryParameter{
		{
			Name:  "product",
			Value: product,
		},
		{
			Name:  "releaseNoteType",
			Value: releaseNoteType,
		},
	}
	// Set the query location to US.
	q.Location = "US"

	// Run the BigQuery query and wait for it to complete.
	job, err := q.Run(ctx)
	if err != nil {
		return nil, err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return nil, status.Err()
	}
	if err := status.Err(); err != nil {
		return nil, status.Err()
	}

	// Read the query results.
	it, err := job.Read(ctx)
	if err != nil {
		return nil, err
	}

	// Initialize a slice to store the retrieved release notes.
	var releaseNotes []ReleaseNote

	// Iterate over the query results and populate the releaseNotes slice.
	rowCount := 0
	for {
		var row []bigquery.Value
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		// Extract data from the rows.
		releaseNote := ReleaseNote{
			ProductName:        getStringValue(row[0]),
			ProductVersionName: getStringValue(row[2]),
			ReleaseNoteType:    getStringValue(row[1]),
			Description:        getStringValue(row[3]),
			PublishedAt:        getStringValue(row[4]),
		}

		// Append the release note to the releaseNotes slice.
		releaseNotes = append(releaseNotes, releaseNote)
		rowCount++
	}

	// Print the number of release notes found for informational purposes.
	if rowCount > 1 {
		fmt.Printf("\nFound %d entires for: %s\n", rowCount, product)
	} else {
		fmt.Printf("\nFound %d entry for : %s\n", rowCount, product)
	}

	// Return the slice of release notes.
	return releaseNotes, nil

}

// getStringValue returns the string value of a bigquery.Value.
func getStringValue(v bigquery.Value) string {
	if v == nil {
		return "NULL"
	}
	return fmt.Sprintf("%v", v)
}

// ReleaseNote represents a release note.
type ReleaseNote struct {
	ProductName        string `bigquery:"product_name" json:"product_name"`
	ProductVersionName string `bigquery:"product_version_name" json:"product_version_name"`
	ReleaseNoteType    string `bigquery:"release_note_type" json:"release_note_type"`
	Description        string `bigquery:"description" json:"description"`
	PublishedAt        string `bigquery:"published_at" json:"published_at"`
}
