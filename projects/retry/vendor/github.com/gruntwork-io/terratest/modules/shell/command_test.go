package shell

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCommandAndGetOutput(t *testing.T) {
	t.Parallel()

	text := "Hello, World"
	cmd := Command{
		Command: "echo",
		Args:    []string{text},
	}

	out := RunCommandAndGetOutput(t, cmd)
	assert.Equal(t, text, strings.TrimSpace(out))
}
