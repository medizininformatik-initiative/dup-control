package container

import (
	"context"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/op/go-logging"
	"os"
)

var log = logging.MustGetLogger("util.container")

type DockerClient interface {
	PullImage(opts docker.PullImageOptions, auth docker.AuthConfiguration) error
	CreateContainer(opts docker.CreateContainerOptions) (*docker.Container, error)
	StartContainerWithContext(id string, hostConfig *docker.HostConfig, ctx context.Context) error
	Logs(opts docker.LogsOptions) error
	RemoveContainer(opts docker.RemoveContainerOptions) error
}

type Runtime struct {
	context    context.Context
	client     DockerClient
	authConfig docker.AuthConfiguration
}

func NewRuntime(client DockerClient, registryUser string, registryPass string) *Runtime {
	return &Runtime{
		context:    context.Background(),
		client:     client,
		authConfig: docker.AuthConfiguration{Username: registryUser, Password: registryPass},
	}
}

func sprintRepositoryName(workpackage string) string {
	return fmt.Sprintf("registry.gitlab.com/smith-phep/polar/%s", workpackage)
}

func sprintImageName(workpackage string, site string) string {
	return fmt.Sprintf("%s:%s", sprintRepositoryName(workpackage), site)
}

func sprintContainerName(prefix string, workpackage string, site string) string {
	return fmt.Sprintf("polar-%s-%s-%s", prefix, workpackage, site)
}

type PullOpts struct {
	Workpackage string
	Site        string
}

func (runtime *Runtime) Pull(pullOpts PullOpts) error {
	opts := docker.PullImageOptions{
		Context:      runtime.context,
		Repository:   sprintRepositoryName(pullOpts.Workpackage),
		Tag:          pullOpts.Site,
		OutputStream: os.Stdout,
	}
	return runtime.client.PullImage(opts, runtime.authConfig)
}

type RunOpts struct {
	User   string
	Env    []string
	Mounts []docker.Mount
}

func dockerFromOpts(pullOpts PullOpts, runOpts RunOpts) *docker.Config {
	return &docker.Config{
		User:   runOpts.User,
		Env:    runOpts.Env,
		Image:  sprintImageName(pullOpts.Workpackage, pullOpts.Site),
		Mounts: runOpts.Mounts,
	}
}

func (runtime *Runtime) Run(containerNamePrefix string, pullOpts PullOpts, runOpts RunOpts) error {
	containerOpts := docker.CreateContainerOptions{
		Context: runtime.context,
		Name:    sprintContainerName(containerNamePrefix, pullOpts.Workpackage, pullOpts.Site),
		Config:  dockerFromOpts(pullOpts, runOpts),
	}
	container, err := runtime.client.CreateContainer(containerOpts)
	if err == nil {
		defer runtime.remove(container)
		err = runtime.client.StartContainerWithContext(container.ID, nil, runtime.context)
	}

	if err == nil {
		logOpts := docker.LogsOptions{
			Stdout:       true,
			Stderr:       true,
			Container:    container.ID,
			OutputStream: os.Stdout,
			ErrorStream:  os.Stderr,
			Follow:       true,
		}
		err = runtime.client.Logs(logOpts)
	}

	return err
}

func (runtime *Runtime) remove(container *docker.Container) {
	removeOpts := docker.RemoveContainerOptions{
		Context: runtime.context,
		ID:      container.ID,
	}
	if err := runtime.client.RemoveContainer(removeOpts); err != nil {
		log.Errorf("Unable to remove container with ID %s, %v", container.ID, err.Error())
	}
}
