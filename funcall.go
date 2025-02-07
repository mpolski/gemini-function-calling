package funcall

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	_ "github.com/mpolski/gemini-function-calling/pkg/fetch"
)

func init() {
	functions.HTTP("funcall", funcall)
}

func funcall(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	// function to run GetReleaseNotes
	fmt.Println("Calling GetReleaseNotes function...")
	releaseNotes, err := releasenotes.FetchReleaseNotes(ctx, "Compute Engine", "FEATURE")
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
