// Package git provides a high-level interface to git commands.
package git

import (
	"os"
)

// LogSimple shows simple log.
func (c *Client) LogSimple() error {
	cmd := c.execCommand("git", "log", "--oneline", "--graph", "--decorate", "-10")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewError("log simple", "git log --oneline --graph --decorate -10", err)
	}
	return nil
}

// LogGraph shows log with graph.
func (c *Client) LogGraph() error {
	cmd := c.execCommand("git", "log", "--graph", "--oneline", "--decorate", "--all")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewError("log graph", "git log --graph --oneline --decorate --all", err)
	}
	return nil
}
