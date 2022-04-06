package testing

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// A cmdExecution represents a completed execution of a single cobra.Command run
type cmdExecution struct {
	err       error
	outStream string
	errStream string
}

func Run(command *cobra.Command, args ...string) *cmdExecution {
	command.SetArgs(args)

	outStream := bytes.NewBufferString("")
	command.SetOut(outStream)

	errStream := bytes.NewBufferString("")
	command.SetErr(errStream)

	err := command.Execute()
	return &cmdExecution{
		err:       err,
		outStream: outStream.String(),
		errStream: errStream.String(),
	}
}

func RunAsync(command *cobra.Command, args ...string) chan *cmdExecution {
	result := make(chan *cmdExecution, 1)

	go func() {
		result <- Run(command, args...)
	}()

	return result
}

func (exe *cmdExecution) Success() bool {
	return exe.err == nil
}

func (exe *cmdExecution) Failure() bool {
	return !exe.Success()
}

func (exe *cmdExecution) Err() error {
	return exe.err
}

func (exe *cmdExecution) OutStream() string {
	return exe.outStream
}

func (exe *cmdExecution) ErrStream() string {
	return exe.errStream
}

func (exe *cmdExecution) AssertSuccess(t *testing.T) bool {
	if !exe.Success() {
		t.Log(exe.Err())
		return assert.Fail(t, "Command should have succeeded")
	}
	return true
}

func (exe *cmdExecution) AssertFailure(t *testing.T) bool {
	if !exe.Failure() {
		return assert.Fail(t, "Command should have failed")
	}
	return true
}

func (exe *cmdExecution) AssertOutStreamContains(t *testing.T, substr string) bool {
	if !strings.Contains(exe.outStream, substr) {
		return assert.Fail(t, fmt.Sprintf("OutStream '%s' should have contained '%s'", exe.outStream, substr))
	}
	return true
}

func (exe *cmdExecution) AssertErrContains(t *testing.T, substr string) bool {
	if !strings.Contains(exe.err.Error(), substr) {
		return assert.Fail(t, fmt.Sprintf("Err '%s' should have contained '%s'", exe.outStream, substr))
	}
	return true
}
