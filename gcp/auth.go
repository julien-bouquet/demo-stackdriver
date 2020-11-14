package gcp

import (
	"context"

	"golang.org/x/oauth2/google"
)

func getProjectID(ctx context.Context) string {
	/// Return authenticated GCP ProjectID
	cred, _ := google.FindDefaultCredentials(ctx)

	return cred.ProjectID
}
