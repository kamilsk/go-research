// Package packer allows to interact with Packer.
package packer

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/shell"
)

// Options are the options for Packer.
type Options struct {
	Template string            // The path to the Packer template
	Vars     map[string]string // The custom vars to pass when running the build command
	Only     string            // If specified, only run the build of this name
	Env      map[string]string // Custom environment variables to set when running Packer
}

// BuildArtifact builds the given Packer template and return the generated Artifact ID.
func BuildArtifact(t *testing.T, options *Options) string {
	artifactID, err := BuildArtifactE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return artifactID
}

// BuildArtifactE builds the given Packer template and return the generated Artifact ID.
func BuildArtifactE(t *testing.T, options *Options) (string, error) {
	logger.Logf(t, "Running Packer to generate a custom artifact for template %s", options.Template)

	cmd := shell.Command{
		Command: "packer",
		Args:    formatPackerArgs(options),
		Env:     options.Env,
	}

	output, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		return "", err
	}

	return extractArtifactID(output)
}

// BuildAmi builds the given Packer template and return the generated AMI ID.
//
// Deprecated: Use BuildArtifact instead.
func BuildAmi(t *testing.T, options *Options) string {
	return BuildArtifact(t, options)
}

// BuildAmiE builds the given Packer template and return the generated AMI ID.
//
// Deprecated: Use BuildArtifactE instead.
func BuildAmiE(t *testing.T, options *Options) (string, error) {
	return BuildArtifactE(t, options)
}

// The Packer machine-readable log output should contain an entry of this format:
//
// AWS: <timestamp>,<builder>,artifact,<index>,id,<region>:<image_id>
// GCP: <timestamp>,<builder>,artifact,<index>,id,<image_id>
//
// For example:
//
// 1456332887,amazon-ebs,artifact,0,id,us-east-1:ami-b481b3de
// 1533742764,googlecompute,artifact,0,id,terratest-packer-example-2018-08-08t15-35-19z
//
func extractArtifactID(packerLogOutput string) (string, error) {
	re := regexp.MustCompile(`.+artifact,\d+?,id,(?:.+?:|)(.+)`)
	matches := re.FindStringSubmatch(packerLogOutput)

	if len(matches) == 2 {
		return matches[1], nil
	}
	return "", errors.New("Could not find Artifact ID pattern in Packer output")
}

// Convert the inputs to a format palatable to packer. The build command should have the format:
//
// packer build [OPTIONS] template
func formatPackerArgs(options *Options) []string {
	args := []string{"build", "-machine-readable"}

	for key, value := range options.Vars {
		args = append(args, "-var", fmt.Sprintf("%s=%s", key, value))
	}

	if options.Only != "" {
		args = append(args, fmt.Sprintf("-only=%s", options.Only))
	}

	return append(args, options.Template)
}
