package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/mpolski/gemini-function-calling/pkg/releasenotes"
)

func init() {
	functions.HTTP("tool", tool)
}

func tool(w http.ResponseWriter, r *http.Request) {

	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		fmt.Println("Set PROJET_ID= in environment variables")
		return
	}

	// model := os.Getenv("MODEL")
	// if model == "" {
	// 	fmt.Println("Set MODEL= in environment variables, e.g. gemini-pro")
	// 	return
	// }
	// modelLocation := os.Getenv("MODEL_LOCATION")
	// if modelLocation == "" {
	// 	fmt.Println("Set MODEL_LOCATION= in environment variables, e.g. us-central1")
	// 	return
	// }

	ctx := context.Background()

	// function to run GetReleaseNotes
	fmt.Println("Calling GetReleaseNotes function...")
	releaseNotes, err := releasenotes.GetReleaseNotes(ctx, projectID, "Compute Engine", "FEATURE")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error calling GetReleaseNotes: %v", err), http.StatusInternalServerError)
		return
	}

	// Write the release notes to the HTTP response.
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(releaseNotes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON response: %v", err), http.StatusInternalServerError)
		return
	}

}
