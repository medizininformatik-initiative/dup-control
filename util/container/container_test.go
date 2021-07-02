package container

import (
	"context"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockClient struct {
	mock.Mock
}

type InterestingPullOpts struct {
	Repository string
	Tag        string
}

func (mock *MockClient) PullImage(opts docker.PullImageOptions, auth docker.AuthConfiguration) error {
	interestingOpts := InterestingPullOpts{Repository: opts.Repository, Tag: opts.Tag}
	args := mock.Called(interestingOpts, auth)
	return args.Error(0)
}

func (mock *MockClient) CreateContainer(opts docker.CreateContainerOptions) (*docker.Container, error) {
	args := mock.Called(opts)
	return args.Get(0).(*docker.Container), args.Error(1)
}

func (mock *MockClient) StartContainerWithContext(id string, hostConfig *docker.HostConfig, ctx context.Context) error {
	args := mock.Called(id, hostConfig, ctx)
	return args.Error(0)
}

func (mock *MockClient) Logs(opts docker.LogsOptions) error {
	args := mock.Called(opts)
	return args.Error(0)
}

func (mock *MockClient) RemoveContainer(opts docker.RemoveContainerOptions) error {
	args := mock.Called(opts)
	return args.Error(0)
}

const pass = "some-pass"
const user = "some-user"
const dic = "dic-anywhere"

func TestPull(t *testing.T) {
	dockerMock := new(MockClient)
	runtime := NewRuntime(dockerMock, user, pass)
	dockerMock.On("PullImage",
		InterestingPullOpts{Repository: "registry.gitlab.com/smith-phep/polar/wp-0", Tag: dic},
		docker.AuthConfiguration{Username: user, Password: pass}).Return(nil)

	_ = runtime.Pull(PullOpts{Workpackage: "wp-0", Site: dic})
}
