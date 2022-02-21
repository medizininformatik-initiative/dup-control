package util

import (
	"bytes"
	"errors"
	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockUpdater struct {
	mock.Mock
}

func (mock *mockUpdater) IsNewerVersionAvailable() (bool, string) {
	args := mock.Called()
	return args.Bool(0), args.String(1)
}
func (mock *mockUpdater) Upgrade() error {
	args := mock.Called()
	return args.Error(0)
}

func setupCommand() (*upgradeCommand, *mockUpdater) {
	updater := new(mockUpdater)
	logger := logging.MustGetLogger("test")
	cmd := NewUpgradeCommand(logger, updater)
	return cmd, updater
}

func TestUpgradeFails(t *testing.T) {
	cmd, updater := setupCommand()

	updater.On("Upgrade").Return(errors.New("some error"))

	command := cmd.Command()
	command.SetArgs([]string{})
	command.SetOut(bytes.NewBufferString(""))
	err := command.Execute()

	assert.Error(t, err, "Error updating polarctl: some error")
}
