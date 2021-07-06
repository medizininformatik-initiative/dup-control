package container

import (
	"context"
	"errors"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockClient struct {
	mock.Mock
}

type mockPullOpts struct {
	Repository string
	Tag        string
}

func (mock *mockClient) PullImage(opts docker.PullImageOptions, auth docker.AuthConfiguration) error {
	mockOpts := mockPullOpts{Repository: opts.Repository, Tag: opts.Tag}
	args := mock.Called(mockOpts, auth)
	return args.Error(0)
}

func (mock *mockClient) CreateContainer(opts docker.CreateContainerOptions) (*docker.Container, error) {
	args := mock.Called(opts)
	return args.Get(0).(*docker.Container), args.Error(1)
}

func (mock *mockClient) StartContainerWithContext(id string, hostConfig *docker.HostConfig, ctx context.Context) error {
	args := mock.Called(id, hostConfig, ctx)
	return args.Error(0)
}

type mockLockOpts struct {
	Container string
	Follow    bool
}

func (mock *mockClient) Logs(opts docker.LogsOptions) error {
	mockOpts := mockLockOpts{Container: opts.Container, Follow: opts.Follow}
	args := mock.Called(mockOpts)
	return args.Error(0)
}

type mockRemoveOpts struct {
	ID    string
	Force bool
}

func (mock *mockClient) RemoveContainer(opts docker.RemoveContainerOptions) error {
	mockOpts := mockRemoveOpts{ID: opts.ID, Force: opts.Force}
	args := mock.Called(mockOpts)
	return args.Error(0)
}

const pass = "some-pass"
const user = "some-user"
const dic = "dic-anywhere"
const id = "102583"

func TestPull(t *testing.T) {
	dockerMock := new(mockClient)
	runtime := NewRuntime(dockerMock, user, pass)
	dockerMock.On("PullImage",
		mockPullOpts{Repository: "registry.gitlab.com/smith-phep/polar/wp-0", Tag: dic},
		docker.AuthConfiguration{Username: user, Password: pass}).Return(nil)

	_ = runtime.Pull(PullOpts{Workpackage: "wp-0", Site: dic})

	dockerMock.AssertExpectations(t)
}

func TestPullError(t *testing.T) {
	dockerMock := new(mockClient)
	runtime := NewRuntime(dockerMock, user, pass)
	dockerMock.On("PullImage", mock.Anything, mock.Anything).Return(errors.New("unable to pull image"))

	err := runtime.Pull(PullOpts{Workpackage: "wp-6", Site: dic})

	assert.Error(t, err, "unable to pull image")
	dockerMock.AssertExpectations(t)
}

func TestRun(t *testing.T) {
	dockerMock := new(mockClient)
	runtime := NewRuntime(dockerMock, user, pass)

	dockerMock.On("CreateContainer",
		mock.Anything).Return(&docker.Container{ID: id}, nil)
	dockerMock.On("StartContainerWithContext",
		id, mock.Anything, mock.Anything).Return(nil)
	dockerMock.On("Logs",
		mockLockOpts{Container: id, Follow: true}).Return(nil)
	dockerMock.On("RemoveContainer",
		mockRemoveOpts{ID: id, Force: false}).Return(nil)

	_ = runtime.Run("prefix",
		PullOpts{Workpackage: "wp-0", Site: dic},
		RunOpts{User: "", Env: []string{}, Mounts: []docker.Mount{}})

	dockerMock.AssertExpectations(t)
}

func TestRunWithCreateError(t *testing.T) {
	dockerMock := new(mockClient)
	runtime := NewRuntime(dockerMock, user, pass)

	dockerMock.On("CreateContainer",
		mock.Anything).Return(&docker.Container{}, errors.New("unable to create container"))

	err := runtime.Run("prefix",
		PullOpts{Workpackage: "wp-0", Site: dic},
		RunOpts{User: "", Env: []string{}, Mounts: []docker.Mount{}})

	assert.Error(t, err, "unable to create container")

	dockerMock.AssertExpectations(t)
}

func TestRunWithStartError(t *testing.T) {
	dockerMock := new(mockClient)
	runtime := NewRuntime(dockerMock, user, pass)

	dockerMock.On("CreateContainer",
		mock.Anything).Return(&docker.Container{ID: id}, nil)
	dockerMock.On("StartContainerWithContext",
		id, mock.Anything, mock.Anything).Return(errors.New("unable to start container"))
	dockerMock.On("RemoveContainer",
		mockRemoveOpts{ID: id, Force: false}).Return(nil)

	err := runtime.Run("prefix",
		PullOpts{Workpackage: "wp-0", Site: dic},
		RunOpts{User: "", Env: []string{}, Mounts: []docker.Mount{}})

	assert.Error(t, err, "unable to start container")

	dockerMock.AssertExpectations(t)
}
