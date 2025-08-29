package cmd

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestAdder_Add_NoArgs_PrintsUsage(t *testing.T) {
	adder := NewAdder()
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	adder.Add([]string{})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	os.Stdout = oldStdout

	output := buf.String()
	if output == "" || output[:5] != "Usage" {
		t.Errorf("Usage not output: %s", output)
	}
}

// mockGitClient for testing
type mockAddGitClient struct {
	addCalled            bool
	addInteractiveCalled bool
	addFiles             []string
	addError             error
}

func (m *mockAddGitClient) Add(files ...string) error {
	m.addCalled = true
	m.addFiles = files
	return m.addError
}

func (m *mockAddGitClient) AddInteractive() error {
	m.addInteractiveCalled = true
	return m.addError
}

// Implement other required methods with no-ops for testing
func (m *mockAddGitClient) GetCurrentBranch() (string, error)                      { return "", nil }
func (m *mockAddGitClient) ListLocalBranches() ([]string, error)                   { return nil, nil }
func (m *mockAddGitClient) ListRemoteBranches() ([]string, error)                  { return nil, nil }
func (m *mockAddGitClient) CheckoutNewBranch(name string) error                    { return nil }
func (m *mockAddGitClient) Push(force bool) error                                  { return nil }
func (m *mockAddGitClient) Pull(rebase bool) error                                 { return nil }
func (m *mockAddGitClient) LogSimple() error                                       { return nil }
func (m *mockAddGitClient) LogGraph() error                                        { return nil }
func (m *mockAddGitClient) CommitAllowEmpty() error                                { return nil }
func (m *mockAddGitClient) ResetHardAndClean() error                               { return nil }
func (m *mockAddGitClient) CleanFiles() error                                      { return nil }
func (m *mockAddGitClient) CleanDirs() error                                       { return nil }
func (m *mockAddGitClient) GetGitStatus() (string, error)                          { return "", nil }
func (m *mockAddGitClient) GetBranchName() (string, error)                         { return "", nil }
func (m *mockAddGitClient) RestoreWorkingDir(paths ...string) error                { return nil }
func (m *mockAddGitClient) RestoreStaged(paths ...string) error                    { return nil }
func (m *mockAddGitClient) RestoreFromCommit(commit string, paths ...string) error { return nil }
func (m *mockAddGitClient) RestoreAll() error                                      { return nil }
func (m *mockAddGitClient) RestoreAllStaged() error                                { return nil }
func (m *mockAddGitClient) Commit(message string) error                            { return nil }
func (m *mockAddGitClient) CommitAmend() error                                     { return nil }
func (m *mockAddGitClient) CommitAmendNoEdit() error                               { return nil }
func (m *mockAddGitClient) CommitAmendWithMessage(message string) error            { return nil }
func (m *mockAddGitClient) Status() (string, error)                                { return "", nil }
func (m *mockAddGitClient) StatusShort() (string, error)                           { return "", nil }
func (m *mockAddGitClient) StatusWithColor() (string, error)                       { return "", nil }
func (m *mockAddGitClient) StatusShortWithColor() (string, error)                  { return "", nil }
func (m *mockAddGitClient) Diff() (string, error)                                  { return "", nil }
func (m *mockAddGitClient) DiffStaged() (string, error)                            { return "", nil }
func (m *mockAddGitClient) DiffHead() (string, error)                              { return "", nil }
func (m *mockAddGitClient) Fetch(prune bool) error                                 { return nil }
func (m *mockAddGitClient) Stash() error                                           { return nil }
func (m *mockAddGitClient) StashList() (string, error)                             { return "", nil }
func (m *mockAddGitClient) StashShow(stash string) error                           { return nil }
func (m *mockAddGitClient) StashApply(stash string) error                          { return nil }
func (m *mockAddGitClient) StashPop(stash string) error                            { return nil }
func (m *mockAddGitClient) StashDrop(stash string) error                           { return nil }
func (m *mockAddGitClient) StashClear() error                                      { return nil }

func TestAdder_Add_GitAddCalled(t *testing.T) {
	mockClient := &mockAddGitClient{}
	adder := &Adder{
		gitClient:    mockClient,
		outputWriter: &bytes.Buffer{},
	}
	adder.Add([]string{"hoge.txt"})
	if !mockClient.addCalled {
		t.Error("Add was not called")
	}
	if len(mockClient.addFiles) != 1 || mockClient.addFiles[0] != "hoge.txt" {
		t.Errorf("Expected files [hoge.txt], got %v", mockClient.addFiles)
	}
}

func TestAdder_Add_GitAddArgs(t *testing.T) {
	mockClient := &mockAddGitClient{}
	adder := &Adder{
		gitClient:    mockClient,
		outputWriter: &bytes.Buffer{},
	}
	adder.Add([]string{"foo.txt", "bar.txt"})

	if !mockClient.addCalled {
		t.Error("Add was not called")
	}

	wantFiles := []string{"foo.txt", "bar.txt"}
	if len(mockClient.addFiles) != len(wantFiles) {
		t.Errorf("Expected %d files, got %d", len(wantFiles), len(mockClient.addFiles))
		return
	}

	for i, expected := range wantFiles {
		if mockClient.addFiles[i] != expected {
			t.Errorf("Expected file %s at index %d, got %s", expected, i, mockClient.addFiles[i])
		}
	}
}

func TestAdder_Add_RunError_PrintsError(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockAddGitClient{
		addError: errors.New("git add failed"),
	}
	adder := &Adder{
		gitClient:    mockClient,
		outputWriter: &buf,
	}
	adder.Add([]string{"foo.txt"})

	output := buf.String()
	if !strings.Contains(output, "Error") {
		t.Errorf("Error message not output: %s", output)
	}
	if !strings.Contains(output, "git add failed") {
		t.Errorf("Expected error message not found: %s", output)
	}
}

func TestAdder_Add_POption_CallsGitAddP(t *testing.T) {
	mockClient := &mockAddGitClient{}
	adder := &Adder{
		gitClient:    mockClient,
		outputWriter: &bytes.Buffer{},
	}
	adder.Add([]string{"-p"})

	if !mockClient.addInteractiveCalled {
		t.Error("AddInteractive was not called")
	}
	if mockClient.addCalled {
		t.Error("Add should not be called for -p option")
	}
}

func TestAdder_Add_POption_Error(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	adder := &Adder{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			cmd := exec.Command("false") // command that always returns error
			return cmd
		},
	}
	adder.Add([]string{"-p"})
	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	os.Stdout = oldStdout
	output := buf.String()
	if output == "" || output[:5] != "error" {
		t.Errorf("Error output not generated with -p option: %s", output)
	}
}

func TestAdder_Add_Interactive(t *testing.T) {
	var buf bytes.Buffer
	adder := &Adder{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			cmd := exec.Command("echo", "interactive add")
			return cmd
		},
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	adder.Add([]string{"-p"})

	_ = w.Close()
	os.Stdout = oldStdout

	_, _ = buf.ReadFrom(r)
	output := buf.String()
	if !strings.Contains(output, "interactive add") {
		t.Errorf("expected interactive add output, got %q", output)
	}
}

func TestAdder_Add(t *testing.T) {
	cases := []struct {
		name        string
		args        []string
		expectedCmd string
		expectError bool
	}{
		{
			name:        "add all files",
			args:        []string{"."},
			expectedCmd: "git add .",
			expectError: false,
		},
		{
			name:        "add specific file",
			args:        []string{"file.txt"},
			expectedCmd: "git add file.txt",
			expectError: false,
		},
		{
			name:        "add multiple files",
			args:        []string{"file1.txt", "file2.txt"},
			expectedCmd: "git add file1.txt file2.txt",
			expectError: false,
		},
		{
			name:        "no args",
			args:        []string{},
			expectedCmd: "",
			expectError: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var actualCmd string
			a := &Adder{
				execCommand: func(name string, args ...string) *exec.Cmd {
					if tc.expectedCmd != "" {
						actualCmd = strings.Join(append([]string{name}, args...), " ")
						if actualCmd != tc.expectedCmd {
							t.Errorf("expected command %q, got %q", tc.expectedCmd, actualCmd)
						}
					}
					return exec.Command("echo")
				},
			}

			// Capture stdout for no args case
			if len(tc.args) == 0 {
				oldStdout := os.Stdout
				r, w, _ := os.Pipe()
				os.Stdout = w

				a.Add(tc.args)

				_ = w.Close()
				os.Stdout = oldStdout

				var buf bytes.Buffer
				_, _ = buf.ReadFrom(r)
				output := buf.String()
				if !strings.Contains(output, "Usage:") {
					t.Errorf("expected usage message, got %q", output)
				}
			} else {
				a.Add(tc.args)
			}
		})
	}
}

func TestAdder_Add_Error(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	a := &Adder{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("false") // Command fails
		},
	}

	a.Add([]string{"file.txt"})

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()
	if !strings.Contains(output, "error:") {
		t.Errorf("expected error message, got %q", output)
	}
}
