// Sample helloworld on App Engine app with ProjectId.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2/google"
)

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
	ctx := context.Background()

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	cred, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		log.Printf("Error to find credentials %s", err)
	}
	message := "Hello, World!, welcome on " + cred.ProjectID + " project."
	fmt.Fprint(w, message)
}
