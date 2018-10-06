package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/packer"
)

// An example of how to test the Packer template in examples/packer-basic-example using Terratest.
func TestPackerGCPBasicExample(t *testing.T) {
	t.Parallel()

	// Get the Project Id to use
	projectID := gcp.GetGoogleProjectIDFromEnvVar(t)

	// Pick a random GCP zone to test in. This helps ensure your code works in all regions.
	zone := gcp.GetRandomZone(t, projectID, nil, nil)

	packerOptions := &packer.Options{
		// The path to where the Packer template is located
		Template: "../examples/packer-basic-example/build.json",

		// Variables to pass to our Packer build using -var options
		Vars: map[string]string{
			"gcp_project_id": projectID,
			"gcp_zone":       zone,
		},

		// Only build the Google Compute Image
		Only: "googlecompute",
	}

	// Make sure the Packer build completes successfully
	imageName := packer.BuildArtifact(t, packerOptions)

	// Delete the Image after we're done
	image := gcp.FetchImage(t, projectID, imageName)
	defer image.DeleteImage(t)
}
