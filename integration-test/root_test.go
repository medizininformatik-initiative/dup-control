//go:build integration

package test

import (
	"git.smith.care/smith/uc-phep/dupctl/cmd"
	. "git.smith.care/smith/uc-phep/dupctl/internal/testing"
	"testing"
)

func TestNoArgs(t *testing.T) {
	exe := Run(cmd.NewRootCmd().Command())

	exe.AssertSuccess(t)
}

func TestHelp(t *testing.T) {
	exe := Run(cmd.NewRootCmd().Command(), "--help")

	exe.AssertSuccess(t)
	exe.AssertOutStreamContains(t, "Usage")
}

func TestH(t *testing.T) {
	exe := Run(cmd.NewRootCmd().Command(), "-h")

	exe.AssertSuccess(t)
	exe.AssertOutStreamContains(t, "Usage")
}

func TestVersion(t *testing.T) {
	exe := Run(cmd.NewRootCmd().Command(), "--version")

	exe.AssertSuccess(t)
	exe.AssertOutStreamContains(t, "version")
}

func TestV(t *testing.T) {
	exe := Run(cmd.NewRootCmd().Command(), "-v")

	exe.AssertSuccess(t)
	exe.AssertOutStreamContains(t, "version")
}

func TestConfig(t *testing.T) {
	exe := Run(cmd.NewRootCmd().Command(), "--config", "/dev/null", "completion", "bash")

	exe.AssertFailure(t)
}

func TestCompletion(t *testing.T) {
	exe := Run(cmd.NewRootCmd().Command(), "completion", "bash")

	exe.AssertSuccess(t)
}
