// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/bmf-san/ggc/v4/git"
)

// Resetter handles reset operations.
type Resetter struct {
	outputWriter io.Writer
	helper       *Helper
	gitClient    git.Clienter
}

// NewResetter creates a new Resetter instance.
func NewResetter() *Resetter {
	return &Resetter{
		outputWriter: os.Stdout,
		helper:       NewHelper(),
		gitClient:    git.NewClient(),
	}
}

// NewResetterWithClient creates a new Resetter instance with a custom git client.
func NewResetterWithClient(client git.Clienter) *Resetter {
	return &Resetter{
		outputWriter: os.Stdout,
		helper:       NewHelper(),
		gitClient:    client,
	}
}

// Reset executes git reset commands.
func (r *Resetter) Reset(args []string) {
	if len(args) == 0 {
		if r.gitClient != nil {
			if err := r.gitClient.ResetHardAndClean(); err != nil {
				_, _ = fmt.Fprintf(r.outputWriter, "Error: %v\n", err)
				return
			}
			return
		}
		r.helper.ShowResetHelp()
		return
	}

	switch args[0] {
	case "clean":
		// Reset to HEAD and clean untracked files
		if err := r.gitClient.ResetHard("HEAD"); err != nil {
			_, _ = fmt.Fprintf(r.outputWriter, "Error: %v\n", err)
			return
		}

		if err := r.gitClient.CleanFiles(); err != nil {
			_, _ = fmt.Fprintf(r.outputWriter, "Error: %v\n", err)
			return
		}
		_, _ = fmt.Fprintf(r.outputWriter, "Reset and clean successful\n")
		return
	default:
		r.helper.ShowResetHelp()
	}
}
