package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

// mockGitClient for commit_test
type mockCommitGitClient struct {
	commitAllowEmptyCalled bool
	commitCalled           bool
	commitMessage          string
	err                    error
}

func (m *mockCommitGitClient) CommitAllowEmpty() error {
	m.commitAllowEmptyCalled = true
	return m.err
}

func (m *mockCommitGitClient) Commit(message string) error {
	m.commitCalled = true
	m.commitMessage = message
	return m.err
}

// Implement all other required methods from git.Clienter interface
func (m *mockCommitGitClient) GetCurrentBranch() (string, error) { return "main", nil }
func (m *mockCommitGitClient) GetBranchName() (string, error)    { return "main", nil }
func (m *mockCommitGitClient) GetGitStatus() (string, error)     { return "", nil }
func (m *mockCommitGitClient) Status() (string, error)          { return "", nil }
func (m *mockCommitGitClient) StatusShort() (string, error)     { return "", nil }
func (m *mockCommitGitClient) StatusWithColor() (string, error) { return "", nil }
func (m *mockCommitGitClient) StatusShortWithColor() (string, error) { return "", nil }
func (m *mockCommitGitClient) Add(files ...string) error                   { return nil }
func (m *mockCommitGitClient) AddInteractive() error                       { return nil }
func (m *mockCommitGitClient) CommitAmend() error                          { return nil }
func (m *mockCommitGitClient) CommitAmendNoEdit() error                    { return nil }
func (m *mockCommitGitClient) CommitAmendWithMessage(message string) error { return nil }
func (m *mockCommitGitClient) Diff() (string, error)                       { return "", nil }
func (m *mockCommitGitClient) DiffStaged() (string, error)                 { return "", nil }
func (m *mockCommitGitClient) DiffHead() (string, error)                   { return "", nil }
func (m *mockCommitGitClient) ListLocalBranches() ([]string, error)        { return []string{}, nil }
func (m *mockCommitGitClient) ListRemoteBranches() ([]string, error)       { return []string{}, nil }
func (m *mockCommitGitClient) CheckoutNewBranch(name string) error { return nil }
func (m *mockCommitGitClient) CheckoutBranch(name string) error    { return nil }
func (m *mockCommitGitClient) CheckoutNewBranchFromRemote(localBranch, remoteBranch string) error {
	return nil
}
func (m *mockCommitGitClient) DeleteBranch(name string) error        { return nil }
func (m *mockCommitGitClient) ListMergedBranches() ([]string, error) { return []string{}, nil }
func (m *mockCommitGitClient) Push(force bool) error               { return nil }
func (m *mockCommitGitClient) Pull(rebase bool) error              { return nil }
func (m *mockCommitGitClient) Fetch(prune bool) error              { return nil }
func (m *mockCommitGitClient) RemoteList() error                   { return nil }
func (m *mockCommitGitClient) RemoteAdd(name, url string) error    { return nil }
func (m *mockCommitGitClient) RemoteRemove(name string) error      { return nil }
func (m *mockCommitGitClient) RemoteSetURL(name, url string) error { return nil }
func (m *mockCommitGitClient) LogSimple() error                    { return nil }
func (m *mockCommitGitClient) LogGraph() error                     { return nil }
func (m *mockCommitGitClient) LogOneline(from, to string) (string, error) { return "", nil }
func (m *mockCommitGitClient) RebaseInteractive(commitCount int) error { return nil }
func (m *mockCommitGitClient) GetUpstreamBranch(branch string) (string, error) {
	return "origin/main", nil
}
func (m *mockCommitGitClient) Stash() error                       { return nil }
func (m *mockCommitGitClient) StashList() (string, error)         { return "", nil }
func (m *mockCommitGitClient) StashShow(stash string) error       { return nil }
func (m *mockCommitGitClient) StashApply(stash string) error      { return nil }
func (m *mockCommitGitClient) StashPop(stash string) error        { return nil }
func (m *mockCommitGitClient) StashDrop(stash string) error       { return nil }
func (m *mockCommitGitClient) StashClear() error                  { return nil }
func (m *mockCommitGitClient) RestoreWorkingDir(paths ...string) error                { return nil }
func (m *mockCommitGitClient) RestoreStaged(paths ...string) error                    { return nil }
func (m *mockCommitGitClient) RestoreFromCommit(commit string, paths ...string) error { return nil }
func (m *mockCommitGitClient) RestoreAll() error                                      { return nil }
func (m *mockCommitGitClient) RestoreAllStaged() error                                { return nil }
func (m *mockCommitGitClient) ResetHardAndClean() error { return nil }
func (m *mockCommitGitClient) ResetHard(commit string) error        { return nil }
func (m *mockCommitGitClient) CleanFiles() error                   { return nil }
func (m *mockCommitGitClient) CleanDirs() error                    { return nil }
func (m *mockCommitGitClient) CleanDryRun() (string, error)        { return "", nil }
func (m *mockCommitGitClient) CleanFilesForce(files []string) error { return nil }
func (m *mockCommitGitClient) TagList(pattern []string) error                { return nil }
func (m *mockCommitGitClient) TagCreate(name string, commit string) error    { return nil }
func (m *mockCommitGitClient) TagCreateAnnotated(name, message string) error { return nil }
func (m *mockCommitGitClient) TagDelete(names []string) error                { return nil }
func (m *mockCommitGitClient) TagPush(remote, name string) error             { return nil }
func (m *mockCommitGitClient) TagPushAll(remote string) error                { return nil }
func (m *mockCommitGitClient) TagShow(name string) error                     { return nil }
func (m *mockCommitGitClient) GetLatestTag() (string, error)                 { return "", nil }
func (m *mockCommitGitClient) TagExists(name string) bool                    { return false }
func (m *mockCommitGitClient) GetTagCommit(name string) (string, error)      { return "abc123", nil }
func (m *mockCommitGitClient) ListFiles() (string, error)                    { return "", nil }
func (m *mockCommitGitClient) GetUpstreamBranchName(branch string) (string, error) {
	return "origin/main", nil
}
func (m *mockCommitGitClient) GetAheadBehindCount(branch, upstream string) (string, error) {
	return "0\t0", nil
}

func TestCommitter_Commit_AllowEmpty(t *testing.T) {
	mockClient := &mockCommitGitClient{}
	var buf bytes.Buffer
	c := &Committer{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"allow-empty"})
	if !mockClient.commitAllowEmptyCalled {
		t.Error("CommitAllowEmpty should be called")
	}
}

func TestCommitter_Commit_Help(t *testing.T) {
	var buf bytes.Buffer
	c := &Committer{
		gitClient:    &mockCommitGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{})

	output := buf.String()
	if output == "" || !strings.Contains(output, "Usage") {
		t.Errorf("Usage should be displayed, but got: %s", output)
	}
}

func TestCommitter_Commit_AllowEmpty_Error(t *testing.T) {
	var buf bytes.Buffer
	c := &Committer{
		gitClient:    &mockCommitGitClient{err: errors.New("fail")},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"allow-empty"})

	output := buf.String()
	if output != "Error: fail\n" {
		t.Errorf("unexpected output: got %q", output)
	}
}

func TestCommitter_Commit_Normal(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockCommitGitClient{}
	c := &Committer{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"test message"})

	if !mockClient.commitCalled {
		t.Error("git commit command should be called with the correct message")
	}
	if mockClient.commitMessage != "test message" {
		t.Errorf("expected commit message 'test message', got '%s'", mockClient.commitMessage)
	}
}

func TestCommitter_Commit_Normal_WithBrackets(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockCommitGitClient{}
	c := &Committer{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"[update]", "test", "message"})

	if !mockClient.commitCalled {
		t.Error("git commit command should be called with the correct message including brackets")
	}
	expectedMessage := "[update] test message"
	if mockClient.commitMessage != expectedMessage {
		t.Errorf("expected commit message '%s', got '%s'", expectedMessage, mockClient.commitMessage)
	}
}

func TestCommitter_Commit_Normal_Error(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockCommitGitClient{err: errors.New("commit failed")}
	c := &Committer{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"test message"})

	output := buf.String()
	if !strings.Contains(output, "Error:") {
		t.Errorf("Expected error message, got: %s", output)
	}
}

func TestCommitter_Commit_Amend_WithMessage(t *testing.T) {
	var buf bytes.Buffer
	commandCalled := false
	c := &Committer{
		gitClient:    &mockCommitGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"amend", "updated message"})
	if !commandCalled {
		t.Error("git commit --amend -m command should be called")
	}
}

func TestCommitter_Commit_Amend_NoEdit(t *testing.T) {
	var buf bytes.Buffer
	commandCalled := false
	c := &Committer{
		gitClient:    &mockCommitGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"amend", "--no-edit"})
	if !commandCalled {
		t.Error("git commit --amend --no-edit command should be called")
	}
}

func TestCommitter_Commit_Amend_Error(t *testing.T) {
	var buf bytes.Buffer
	c := &Committer{
		gitClient:    &mockCommitGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"amend", "test message"})
	output := buf.String()
	if !strings.Contains(output, "Error:") {
		t.Errorf("Expected error message, got: %s", output)
	}
}

func TestCommitter_Commit_Normal_NoSpace(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockCommitGitClient{}
	c := &Committer{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"[update]hoge"})

	if !mockClient.commitCalled {
		t.Error("git commit command should be called with the correct message without spaces")
	}
	if mockClient.commitMessage != "[update]hoge" {
		t.Errorf("expected commit message '[update]hoge', got '%s'", mockClient.commitMessage)
	}
}

func TestCommitter_Commit_Amend_WithMultiWordMessage(t *testing.T) {
	var buf bytes.Buffer
	commandCalled := false
	c := &Committer{
		gitClient:    &mockCommitGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"amend", "[update]", "message", "with", "spaces"})
	if !commandCalled {
		t.Error("git commit --amend -m command should be called with the complete message")
	}
}
