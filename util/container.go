package container

import (
	"context"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/op/go-logging"
	"os"
)

var log = logging.MustGetLogger("util.container")

type Runtime struct {
	context    context.Context
	client     *docker.Client
	authConfig docker.AuthConfiguration
}

func NewRuntime(registryUser string, registryPass string) (*Runtime, error) {
	cli, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, fmt.Errorf("cannot instantiate docker client, %w", err)
	}

	return &Runtime{
		context:    context.Background(),
		client:     cli,
		authConfig: docker.AuthConfiguration{Username: registryUser, Password: registryPass},
	}, nil
}

func SprintRepositoryName(workpackage string) string {
	return fmt.Sprintf("registry.gitlab.com/smith-phep/polar/%s", workpackage)
}

func SprintImageName(workpackage string, site string) string {
	return fmt.Sprintf("%s:%s", SprintRepositoryName(workpackage), site)
}

func (runtime *Runtime) Pull(workpackage string, site string) error {
	opts := docker.PullImageOptions{
		Context:      runtime.context,
		Repository:   SprintRepositoryName(workpackage),
		Tag:          site,
		OutputStream: os.Stdout,
	}
	return runtime.client.PullImage(opts, runtime.authConfig)
}

func (runtime *Runtime) Run(containerNamePrefix string, workpackage string, site string) error {
	containerOpts := docker.CreateContainerOptions{
		Context: runtime.context,
		Name:    fmt.Sprintf("polar-%s-%s-%s", containerNamePrefix, workpackage, site),
		Config: &docker.Config{
			Image: SprintImageName(workpackage, site),
		},
	}
	container, err := runtime.client.CreateContainer(containerOpts)
	if err == nil && container != nil {
		err = runtime.client.StartContainerWithContext(container.ID, nil, runtime.context)
	}

	if err == nil && container != nil {
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

	if container != nil {
		removeOpts := docker.RemoveContainerOptions{
			Context: runtime.context,
			ID:      container.ID,
		}
		if removeErr := runtime.client.RemoveContainer(removeOpts); removeErr != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Unable to remove container with ID %s, %s", container.ID, removeErr.Error())
		}
	}

	return err
}
