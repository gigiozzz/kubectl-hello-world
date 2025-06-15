package cmd

import (
	"bytes"
	"strings"
	"testing"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func TestVersionCommand(t *testing.T) {
	tests := []struct {
		name      string
		short     bool
		wantShort bool
		wantLong  bool
	}{
		{
			name:      "short version",
			short:     true,
			wantShort: true,
			wantLong:  false,
		},
		{
			name:      "full version",
			short:     false,
			wantShort: false,
			wantLong:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original values
			originalVersion := Version
			originalCommitHash := CommitHash
			originalBuildDate := BuildDate

			// Set test values
			Version = "test-version"
			CommitHash = "test-commit"
			BuildDate = "test-date"

			// Restore original values after test
			defer func() {
				Version = originalVersion
				CommitHash = originalCommitHash
				BuildDate = originalBuildDate
			}()

			// Create test streams
			out := &bytes.Buffer{}
			streams := genericclioptions.IOStreams{
				Out: out,
			}

			// Create and run command
			options := NewVersionOptions(streams)
			options.Short = tt.short

			err := options.Run()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			output := out.String()

			if tt.wantShort {
				expected := "test-version\n"
				if output != expected {
					t.Errorf("expected short output %q, got %q", expected, output)
				}
			}

			if tt.wantLong {
				expectedParts := []string{
					"kubectl-hello-world version: test-version",
					"Commit: test-commit",
					"Built: test-date",
					"Go version:",
					"OS/Arch:",
				}

				for _, part := range expectedParts {
					if !strings.Contains(output, part) {
						t.Errorf("expected output to contain %q, got %q", part, output)
					}
				}
			}
		})
	}
}

func TestNewCmdVersion(t *testing.T) {
	streams := genericclioptions.IOStreams{}
	cmd := NewCmdVersion(streams)

	if cmd.Use != "version" {
		t.Errorf("expected Use to be 'version', got %q", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("expected Short description to be set")
	}

	// Check that the short flag exists
	flag := cmd.Flags().Lookup("short")
	if flag == nil {
		t.Error("expected --short flag to exist")
	}
}
