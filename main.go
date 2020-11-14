// Sample helloworld on App Engine app with ProjectId.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/julien-bouquet/demo-stackdriver/gcp"
)

// Use to set task_id in Resource of logs
var headerNameRequestID = "X-Request-ID"

func main() {
	http.HandleFunc("/", indexHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Test")
	ctx := context.Background()
	requestID := getOrCreateRequestID(r.Header)
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	client := gcp.CreateClientLogger(ctx)
	defer client.Close()

	ctx = gcp.InitializeLogger(ctx, client, requestID)

	requestMetadata := map[string]interface{}{
		"url":    r.URL.Path,
		"header": r.Header,
	}

	gcp.Debug(ctx, "Request Recevied", requestMetadata)

	doSomething(ctx)

	writeResponse(w, *r, requestMetadata)
}

func doSomething(ctx context.Context) {
	gcp.Warn(ctx, "Do something", nil)
}

func getOrCreateRequestID(header http.Header) string {
	requestID := header.Get(headerNameRequestID)
	if requestID != "" {
		return requestID
	}
	return uuid.New().String()
}

func writeResponse(w http.ResponseWriter, r http.Request, metadata map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	metadata["message"] = "Hello, world !"
	metadataJSON, _ := json.Marshal(metadata)
	io.WriteString(w, string(metadataJSON))
}
