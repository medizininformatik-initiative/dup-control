//go:build integration

package test

import (
	"git.smith.care/smith/uc-phep/polar/polarctl/cmd"
	. "git.smith.care/smith/uc-phep/polar/polarctl/internal/testing"
	"testing"
)

func TestUpgrade(t *testing.T) {
	exe := Run(cmd.NewRootCmd().Command(), "upgrade")

	exe.AssertSuccess(t)
	//assert.True(t, exe.OutStreamContains("No new version available")) TODO Find way to analyse logs
}

func TestOfflineUpgrade(t *testing.T) {
	exe := Run(cmd.NewRootCmd().Command(), "--offline", "upgrade")

	exe.AssertFailure(t)
	exe.AssertErrContains(t, "cannot upgrade in --offline mode")
	exe.AssertOutStreamContains(t, "Usage")
}
