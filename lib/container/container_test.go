package container

import (
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

func (mock *mockClient) StartContainer(id string, hostConfig *docker.HostConfig) error {
	args := mock.Called(id, hostConfig)
	return args.Error(0)
}

func (mock *mockClient) StopContainer(id string, timeout uint) error {
	args := mock.Called(id, timeout)
	return args.Error(0)
}

type mockLogOpts struct {
	Container string
	Follow    bool
}

func (mock *mockClient) Logs(opts docker.LogsOptions) error {
	mockOpts := mockLogOpts{Container: opts.Container, Follow: opts.Follow}
	args := mock.Called(mockOpts)
	return args.Error(0)
}

type mockListOpts struct {
	Filters map[string][]string
}

func (mock *mockClient) ListContainers(opts docker.ListContainersOptions) ([]docker.APIContainers, error) {
	mockOpts := mockListOpts{Filters: opts.Filters}
	args := mock.Called(mockOpts)
	return args.Get(0).([]docker.APIContainers), args.Error(1)
}

const registry = "registry.gitlab.com/some-registry"
const namespace = "test"
const pass = "some-pass"
const user = "some-user"
const dic = "dic-anywhere"
const id = "102583"

func TestPull(t *testing.T) {
	dockerMock := new(mockClient)
	runtime := NewRuntime(dockerMock, registry, namespace, user, pass)
	dockerMock.On("PullImage",
		mockPullOpts{Repository: "registry.gitlab.com/some-registry/wp-0", Tag: dic},
		docker.AuthConfiguration{Username: user, Password: pass}).Return(nil)

	_ = runtime.Pull(PullOpts{Image: "wp-0", Tag: dic})

	dockerMock.AssertExpectations(t)
}

func TestPullError(t *testing.T) {
	dockerMock := new(mockClient)
	runtime := NewRuntime(dockerMock, registry, namespace, user, pass)
	dockerMock.On("PullImage", mock.Anything, mock.Anything).Return(errors.New("unable to pull image"))

	err := runtime.Pull(PullOpts{Image: "wp-6", Tag: dic})

	assert.Error(t, err, "unable to pull image")
	dockerMock.AssertExpectations(t)
}

func TestRun(t *testing.T) {
	dockerMock := new(mockClient)
	runtime := NewRuntime(dockerMock, registry, namespace, user, pass)

	dockerMock.On("CreateContainer",
		mock.Anything).Return(&docker.Container{ID: id}, nil)
	dockerMock.On("StartContainer",
		id, mock.Anything, mock.Anything).Return(nil)
	dockerMock.On("Logs",
		mockLogOpts{Container: id, Follow: true}).Return(nil)
	dockerMock.On("ListContainers",
		mockListOpts{Filters: map[string][]string{"id": {id}}}).Return([]docker.APIContainers{}, nil)

	_ = runtime.Run("prefix",
		PullOpts{Image: "wp-0", Tag: dic},
		RunOpts{User: "", Env: []string{}, Mounts: []docker.HostMount{}})

	dockerMock.AssertExpectations(t)
}

func TestRunWithCreateError(t *testing.T) {
	dockerMock := new(mockClient)
	runtime := NewRuntime(dockerMock, registry, namespace, user, pass)

	dockerMock.On("CreateContainer",
		mock.Anything).Return(&docker.Container{}, errors.New("unable to create container"))

	err := runtime.Run("prefix",
		PullOpts{Image: "wp-0", Tag: dic},
		RunOpts{User: "", Env: []string{}, Mounts: []docker.HostMount{}})

	assert.Error(t, err, "unable to create container")

	dockerMock.AssertExpectations(t)
}

func TestRunWithStartError(t *testing.T) {
	dockerMock := new(mockClient)
	runtime := NewRuntime(dockerMock, registry, namespace, user, pass)

	dockerMock.On("CreateContainer",
		mock.Anything).Return(&docker.Container{ID: id}, nil)
	dockerMock.On("StartContainer",
		id, mock.Anything, mock.Anything).Return(errors.New("unable to start container"))
	dockerMock.On("ListContainers",
		mockListOpts{Filters: map[string][]string{"id": {id}}}).Return([]docker.APIContainers{}, nil)

	err := runtime.Run("prefix",
		PullOpts{Image: "wp-0", Tag: dic},
		RunOpts{User: "", Env: []string{}, Mounts: []docker.HostMount{}})

	assert.Error(t, err, "unable to start container")

	dockerMock.AssertExpectations(t)
}

func TestRunWithLogError(t *testing.T) {
	dockerMock := new(mockClient)
	runtime := NewRuntime(dockerMock, registry, namespace, user, pass)

	dockerMock.On("CreateContainer",
		mock.Anything).Return(&docker.Container{ID: id}, nil)
	dockerMock.On("StartContainer",
		id, mock.Anything, mock.Anything).Return(nil)
	dockerMock.On("Logs",
		mockLogOpts{Container: id, Follow: true}).Return(errors.New("unable to get container logs"))
	dockerMock.On("ListContainers",
		mockListOpts{Filters: map[string][]string{"id": {id}}}).Return([]docker.APIContainers{{ID: id}}, nil)
	dockerMock.On("StopContainer",
		id, uint(0)).Return(nil)

	err := runtime.Run("prefix",
		PullOpts{Image: "wp-0", Tag: dic},
		RunOpts{User: "", Env: []string{}, Mounts: []docker.HostMount{}})

	assert.Error(t, err, "unable to get container logs")

	dockerMock.AssertExpectations(t)
}
