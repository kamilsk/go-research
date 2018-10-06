package gcp

import (
	"os"
	"testing"
)

var credsEnvVars = []string{
	"GOOGLE_CREDENTIALS",
	"GOOGLE_CLOUD_KEYFILE_JSON",
	"GCLOUD_KEYFILE_JSON",
	"GOOGLE_USE_DEFAULT_CREDENTIALS",
}

var projectEnvVars = []string{
	"GOOGLE_PROJECT",
	"GOOGLE_CLOUD_PROJECT",
	"GOOGLE_CLOUD_PROJECT_ID",
	"GCLOUD_PROJECT",
	"CLOUDSDK_CORE_PROJECT",
}

var regionEnvVars = []string{
	"GOOGLE_REGION",
	"GCLOUD_REGION",
	"CLOUDSDK_COMPUTE_REGION",
}

// GetGoogleCredentialsFromEnvVar returns the Credentials for use with testing.
func GetGoogleCredentialsFromEnvVar(t *testing.T) string {
	return getFirstNonEmptyValOrEmptyString(t, credsEnvVars)
}

// GetGoogleProjectIDFromEnvVar returns the Project Id for use with testing.
func GetGoogleProjectIDFromEnvVar(t *testing.T) string {
	return getFirstNonEmptyValOrFatal(t, projectEnvVars)
}

// GetGoogleRegionFromEnvVar returns the Region for use with testing.
func GetGoogleRegionFromEnvVar(t *testing.T) string {
	return getFirstNonEmptyValOrFatal(t, regionEnvVars)
}

// getFirstNonEmptyValOrFatal returns the first non-empty value from ks, or throws a fatal
func getFirstNonEmptyValOrFatal(t *testing.T, ks []string) string {
	v := getFirstNonEmptyValOrEmptyString(t, ks)
	if v == "" {
		t.Fatalf("All of the following env vars %v are empty. At least one must be non-empty.", ks)
	}

	return v
}

// getFirstNonEmptyValOrFatal returns the first non-empty value from ks, or returns the empty string
func getFirstNonEmptyValOrEmptyString(t *testing.T, ks []string) string {
	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}

	return ""
}
