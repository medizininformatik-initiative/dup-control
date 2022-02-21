//go:build integration

package test

import (
	"git.smith.care/smith/uc-phep/polar/polarctl/cmd"
	. "git.smith.care/smith/uc-phep/polar/polarctl/internal/testing"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

func TestAnalyze(t *testing.T) {
	RemoveImage(d, "registry.gitlab.com/smith-phep/polar/ci-test-dummy-analysis:latest")
	exeF := RunAsync(cmd.NewRootCmd().Command(), "analyze", "--wp", "ci-test-dummy")

	time.Sleep(10 * time.Second)
	foundContainer := FindContainerByName(d, "polar-analysis-ci-test-dummy-analysis-latest")
	KillContainer(d, foundContainer)

	exe := <-exeF
	exe.AssertSuccess(t)
	assert.True(t, foundContainer != nil, "No running container could be found")
}

func TestAnalyzeOfflineFailsWithoutImage(t *testing.T) {
	RemoveImage(d, "registry.gitlab.com/smith-phep/polar/ci-test-dummy-analysis:latest")
	exe := Run(cmd.NewRootCmd().Command(), "--offline", "analyze", "--wp", "ci-test-dummy")

	exe.AssertFailure(t)
	exe.AssertErrContains(t, "no such image")
}

func TestAnalyzeWritesToOutputLocal(t *testing.T) {
	exe := Run(cmd.NewRootCmd().Command(), "analyze", "--wp", "ci-test-dummy", "-e", "SLEEP=0")

	doneFile, _ := ioutil.ReadFile("outputLocal/analysis-done")

	exe.AssertSuccess(t)
	assert.NotNil(t, doneFile, "Container didn't write to outputLocal")
}

func TestAnalyzeWritesToOutputGlobal(t *testing.T) {
	exe := Run(cmd.NewRootCmd().Command(), "analyze", "--wp", "ci-test-dummy", "-e", "SLEEP=0")

	doneFile, _ := ioutil.ReadFile("outputGlobal/analysis-done")

	exe.AssertSuccess(t)
	assert.NotNil(t, doneFile, "Container didn't write to outputGlobal")
}

func TestAnalyzeRecognizesEnv(t *testing.T) {
	exe := Run(cmd.NewRootCmd().Command(), "analyze", "--wp", "ci-test-dummy", "-e", "SLEEP=0", "-e", "FOO=BAR")

	doneFile, _ := ioutil.ReadFile("outputGlobal/analysis-done")

	exe.AssertSuccess(t)
	assert.True(t, strings.Contains(string(doneFile), "FOO=BAR"), "Couldn't find given env var in output file")
}

func TestAnalyzeRecognizesVersion(t *testing.T) {
	exe := Run(cmd.NewRootCmd().Command(), "analyze", "--wp", "ci-test-dummy", "--version", "absent")

	exe.AssertFailure(t)
	exe.AssertErrContains(t, "unable to pull image registry.gitlab.com/smith-phep/polar/ci-test-dummy-analysis:absent")
}
