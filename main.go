// Sample helloworld on App Engine app with ProjectId.
package main

import (
	"context"
	"net/http"
	"os"

	"github.com/google/uuid"
	"golang.org/x/oauth2/google"
)

// Use to set task_id in Resource of logs
var headerNameRequestID = "X-Request-ID"

func main() {
	http.HandleFunc("/", indexHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	requestID := getOrCreateRequestID(r.Header)
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	client := createClientLogger(ctx)
	defer client.Close()

	ctx, logger := initializeLogging(ctx, client, requestID)

	debug(ctx, logger, "Request Recevied", map[string]interface{}{
		"url":    r.URL,
		"header": r.Header,
	})

	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "application/json")

}

func getProjectID(ctx context.Context) string {
	/// Return authenticated GCP ProjectID
	cred, _ := google.FindDefaultCredentials(ctx)

	return cred.ProjectID
}

func getOrCreateRequestID(header http.Header) string {
	requestID := header.Get(headerNameRequestID)
	if requestID != "" {
		return requestID
	}
	return uuid.New().String()
}
