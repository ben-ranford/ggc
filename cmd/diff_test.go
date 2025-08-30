package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestDiffer_Diff(t *testing.T) {
	cases := []struct {
		name           string
		args           []string
		expectedCmds   []string
		mockOutput     []byte
		mockError      error
		expectedOutput string
	}{
		{
			name:           "diff no args",
			args:           []string{},
			expectedCmds:   []string{"git diff HEAD"},
			mockOutput:     []byte("diff --git a/file.go b/file.go\nindex 1234567..abcdefg 100644\n--- a/file.go\n+++ b/file.go\n@@ -1,3 +1,4 @@\n package main\n+import \"fmt\"\n func main() {\n"),
			mockError:      nil,
			expectedOutput: "diff --git a/file.go b/file.go",
		},
		{
			name:           "diff unstaged",
			args:           []string{"unstaged"},
			expectedCmds:   []string{"git diff"},
			mockOutput:     []byte("diff --git a/unstaged.go b/unstaged.go\nindex 1234567..abcdefg 100644\n--- a/unstaged.go\n+++ b/unstaged.go\n@@ -1,3 +1,4 @@\n package main\n+import \"fmt\"\n func main() {\n"),
			mockError:      nil,
			expectedOutput: "diff --git a/unstaged.go b/unstaged.go",
		},
		{
			name:           "diff staged",
			args:           []string{"staged"},
			expectedCmds:   []string{"git diff --staged"},
			mockOutput:     []byte("diff --git a/staged.go b/staged.go\nindex 1234567..abcdefg 100644\n--- a/staged.go\n+++ b/staged.go\n@@ -1,3 +1,4 @@\n package main\n+import \"fmt\"\n func main() {\n"),
			mockError:      nil,
			expectedOutput: "diff --git a/staged.go b/staged.go",
		},
		{
			name:           "invalid command",
			args:           []string{"invalid"},
			expectedCmds:   nil,
			mockOutput:     nil,
			mockError:      nil,
			expectedOutput: "Usage: ggc diff",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			d := &Differ{
				outputWriter: &buf,
				helper:       NewHelper(),
			}
			d.helper.outputWriter = &buf
			d.Diff(tc.args)

			output := buf.String()
			if !strings.Contains(output, tc.expectedOutput) {
				t.Errorf("expected output to contain %q, got %q", tc.expectedOutput, output)
			}
		})
	}
}
