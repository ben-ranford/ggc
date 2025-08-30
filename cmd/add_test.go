package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestAdder_Add_NoArgs_PrintsUsage(t *testing.T) {
	mockClient := &mockAddGitClient{}
	var buf bytes.Buffer
	adder := &Adder{
		gitClient:    mockClient,
		outputWriter: &buf,
	}

	adder.Add([]string{})

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
	addInteractiveError  error
}

func (m *mockAddGitClient) Add(files ...string) error {
	m.addCalled = true
	m.addFiles = files
	return m.addError
}

func (m *mockAddGitClient) AddInteractive() error {
	m.addInteractiveCalled = true
	return m.addInteractiveError
}

// Repository Information methods
func (m *mockAddGitClient) GetCurrentBranch() (string, error) { return "main", nil }
func (m *mockAddGitClient) GetBranchName() (string, error)    { return "main", nil }
func (m *mockAddGitClient) GetGitStatus() (string, error)     { return "", nil }

// Status Operations methods
func (m *mockAddGitClient) Status() (string, error)               { return "", nil }
func (m *mockAddGitClient) StatusShort() (string, error)          { return "", nil }
func (m *mockAddGitClient) StatusWithColor() (string, error)      { return "", nil }
func (m *mockAddGitClient) StatusShortWithColor() (string, error) { return "", nil }

// Commit Operations methods
func (m *mockAddGitClient) Commit(message string) error                 { return nil }
func (m *mockAddGitClient) CommitAmend() error                          { return nil }
func (m *mockAddGitClient) CommitAmendNoEdit() error                    { return nil }
func (m *mockAddGitClient) CommitAmendWithMessage(message string) error { return nil }
func (m *mockAddGitClient) CommitAllowEmpty() error                     { return nil }

// Diff Operations methods
func (m *mockAddGitClient) Diff() (string, error)       { return "", nil }
func (m *mockAddGitClient) DiffStaged() (string, error) { return "", nil }
func (m *mockAddGitClient) DiffHead() (string, error)   { return "", nil }

// Branch Operations methods
func (m *mockAddGitClient) ListLocalBranches() ([]string, error) { return []string{"main"}, nil }
func (m *mockAddGitClient) ListRemoteBranches() ([]string, error) {
	return []string{"origin/main"}, nil
}
func (m *mockAddGitClient) CheckoutNewBranch(name string) error { return nil }
func (m *mockAddGitClient) CheckoutBranch(name string) error    { return nil }
func (m *mockAddGitClient) CheckoutNewBranchFromRemote(localBranch, remoteBranch string) error {
	return nil
}
func (m *mockAddGitClient) DeleteBranch(name string) error        { return nil }
func (m *mockAddGitClient) ListMergedBranches() ([]string, error) { return []string{}, nil }

// Remote Operations methods
func (m *mockAddGitClient) Push(force bool) error               { return nil }
func (m *mockAddGitClient) Pull(rebase bool) error              { return nil }
func (m *mockAddGitClient) Fetch(prune bool) error              { return nil }
func (m *mockAddGitClient) RemoteList() error                   { return nil }
func (m *mockAddGitClient) RemoteAdd(name, url string) error    { return nil }
func (m *mockAddGitClient) RemoteRemove(name string) error      { return nil }
func (m *mockAddGitClient) RemoteSetURL(name, url string) error { return nil }

// Tag Operations methods
func (m *mockAddGitClient) TagList(pattern []string) error                { return nil }
func (m *mockAddGitClient) TagCreate(name string, commit string) error    { return nil }
func (m *mockAddGitClient) TagCreateAnnotated(name, message string) error { return nil }
func (m *mockAddGitClient) TagDelete(names []string) error                { return nil }
func (m *mockAddGitClient) TagPush(remote, name string) error             { return nil }
func (m *mockAddGitClient) TagPushAll(remote string) error                { return nil }
func (m *mockAddGitClient) TagShow(name string) error                     { return nil }
func (m *mockAddGitClient) GetLatestTag() (string, error)                 { return "v1.0.0", nil }
func (m *mockAddGitClient) TagExists(name string) bool                    { return false }
func (m *mockAddGitClient) GetTagCommit(name string) (string, error)      { return "abc123", nil }

// Log Operations methods
func (m *mockAddGitClient) LogSimple() error                           { return nil }
func (m *mockAddGitClient) LogGraph() error                            { return nil }
func (m *mockAddGitClient) LogOneline(from, to string) (string, error) { return "", nil }

// Rebase Operations methods
func (m *mockAddGitClient) RebaseInteractive(commitCount int) error { return nil }
func (m *mockAddGitClient) GetUpstreamBranch(branch string) (string, error) {
	return "origin/main", nil
}

// Stash Operations methods
func (m *mockAddGitClient) Stash() error                  { return nil }
func (m *mockAddGitClient) StashList() (string, error)    { return "", nil }
func (m *mockAddGitClient) StashShow(stash string) error  { return nil }
func (m *mockAddGitClient) StashApply(stash string) error { return nil }
func (m *mockAddGitClient) StashPop(stash string) error   { return nil }
func (m *mockAddGitClient) StashDrop(stash string) error  { return nil }
func (m *mockAddGitClient) StashClear() error             { return nil }

// Restore Operations methods
func (m *mockAddGitClient) RestoreWorkingDir(paths ...string) error                { return nil }
func (m *mockAddGitClient) RestoreStaged(paths ...string) error                    { return nil }
func (m *mockAddGitClient) RestoreFromCommit(commit string, paths ...string) error { return nil }
func (m *mockAddGitClient) RestoreAll() error                                      { return nil }
func (m *mockAddGitClient) RestoreAllStaged() error                                { return nil }

// Reset and Clean Operations methods
func (m *mockAddGitClient) ResetHardAndClean() error             { return nil }
func (m *mockAddGitClient) ResetHard(commit string) error        { return nil }
func (m *mockAddGitClient) CleanFiles() error                    { return nil }
func (m *mockAddGitClient) CleanDirs() error                     { return nil }
func (m *mockAddGitClient) CleanDryRun() (string, error)         { return "", nil }
func (m *mockAddGitClient) CleanFilesForce(files []string) error { return nil }

// Utility Operations methods
func (m *mockAddGitClient) ListFiles() (string, error) { return "", nil }
func (m *mockAddGitClient) GetUpstreamBranchName(branch string) (string, error) {
	return "origin/main", nil
}
func (m *mockAddGitClient) GetAheadBehindCount(branch, upstream string) (string, error) {
	return "0	0", nil
}

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
	mockClient := &mockAddGitClient{addInteractiveError: errors.New("interactive add failed")}
	var buf bytes.Buffer
	adder := &Adder{
		gitClient:    mockClient,
		outputWriter: &buf,
	}
	adder.Add([]string{"-p"})

	output := buf.String()
	if output == "" || output[:5] != "Error" {
		t.Errorf("Error output not generated with -p option: %s", output)
	}
}

func TestAdder_Add_Interactive(t *testing.T) {
	mockClient := &mockAddGitClient{}
	var buf bytes.Buffer
	adder := &Adder{
		gitClient:    mockClient,
		outputWriter: &buf,
	}

	adder.Add([]string{"-p"})

	// Check that AddInteractive was called
	if !mockClient.addInteractiveCalled {
		t.Error("AddInteractive should be called for -p option")
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
			mockClient := &mockAddGitClient{}
			var buf bytes.Buffer
			a := &Adder{
				gitClient:    mockClient,
				outputWriter: &buf,
			}

			a.Add(tc.args)

			// Check output for no args case
			if len(tc.args) == 0 {
				output := buf.String()
				if !strings.Contains(output, "Usage:") {
					t.Errorf("expected usage message, got %q", output)
				}
			}
		})
	}
}

func TestAdder_Add_Error(t *testing.T) {
	mockClient := &mockAddGitClient{addError: errors.New("git add failed")}
	var buf bytes.Buffer
	a := &Adder{
		gitClient:    mockClient,
		outputWriter: &buf,
	}

	a.Add([]string{"file.txt"})

	output := buf.String()
	if !strings.Contains(output, "Error:") {
		t.Errorf("expected error message, got %q", output)
	}
}
