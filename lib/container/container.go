package container

import (
	"context"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/op/go-logging"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

var log = logging.MustGetLogger("util.container")

type DockerClient interface {
	PullImage(opts docker.PullImageOptions, auth docker.AuthConfiguration) error
	CreateContainer(opts docker.CreateContainerOptions) (*docker.Container, error)
	StartContainer(id string, hostConfig *docker.HostConfig) error
	Logs(opts docker.LogsOptions) error
	StopContainer(id string, timeout uint) error
	ListContainers(opts docker.ListContainersOptions) ([]docker.APIContainers, error)
}

type Runtime struct {
	context    context.Context
	client     DockerClient
	registry   string
	namespace  string
	authConfig docker.AuthConfiguration
}

func NewRuntime(client DockerClient, registry string, namespace string, registryUser string, registryPass string) *Runtime {
	return &Runtime{
		context:    context.Background(),
		client:     client,
		registry:   registry,
		namespace:  namespace,
		authConfig: docker.AuthConfiguration{Username: registryUser, Password: registryPass},
	}
}

func (runtime *Runtime) sprintRepositoryName(dup string) string {
	return fmt.Sprintf("%s/%s", runtime.registry, dup)
}

func (runtime *Runtime) sprintImageName(dup string, site string) string {
	return fmt.Sprintf("%s:%s", runtime.sprintRepositoryName(dup), site)
}

func (runtime *Runtime) sprintContainerName(prefix string, dup string, site string) string {
	return fmt.Sprintf("%s-%s-%s-%s", runtime.namespace, prefix, dup, site)
}

type PullOpts struct {
	Image string
	Tag   string
}

func (runtime *Runtime) Pull(pullOpts PullOpts) error {
	imageName := runtime.sprintRepositoryName(pullOpts.Image)
	opts := docker.PullImageOptions{
		Context:      runtime.context,
		Repository:   imageName,
		Tag:          pullOpts.Tag,
		OutputStream: os.Stdout,
	}
	if err := runtime.client.PullImage(opts, runtime.authConfig); err != nil {
		return fmt.Errorf("unable to pull image %s:%s, %w", imageName, pullOpts.Tag, err)
	} else {
		return nil
	}
}

type RunOpts struct {
	User   string
	Env    []string
	Mounts []docker.HostMount
}

func (runtime *Runtime) containerConfigFromOpts(pullOpts PullOpts, runOpts RunOpts) *docker.Config {
	return &docker.Config{
		User:  runOpts.User,
		Env:   runOpts.Env,
		Image: runtime.sprintImageName(pullOpts.Image, pullOpts.Tag),
	}
}

func containerHostConfigFromOpts(opts RunOpts) *docker.HostConfig {
	return &docker.HostConfig{
		Mounts:     opts.Mounts,
		AutoRemove: true,
	}
}

func (runtime *Runtime) Run(containerNamePrefix string, pullOpts PullOpts, runOpts RunOpts) error {
	containerName := runtime.sprintContainerName(containerNamePrefix, pullOpts.Image, pullOpts.Tag)
	containerOpts := docker.CreateContainerOptions{
		Context:    runtime.context,
		Name:       containerName,
		Config:     runtime.containerConfigFromOpts(pullOpts, runOpts),
		HostConfig: containerHostConfigFromOpts(runOpts),
	}
	container, err := runtime.client.CreateContainer(containerOpts)
	if err == nil {
		defer runtime.terminateIfRunning(container, 500*time.Millisecond)
		runtime.registerInterruptTermination(container)
		err = runtime.client.StartContainer(container.ID, nil)
	}

	if err == nil {
		logOpts := docker.LogsOptions{
			Context:      runtime.context,
			Stdout:       true,
			Stderr:       true,
			Container:    container.ID,
			OutputStream: os.Stdout,
			ErrorStream:  os.Stderr,
			Follow:       true,
		}
		err = runtime.client.Logs(logOpts)
	}
	if err != nil {
		return fmt.Errorf("unable to start container %s, %w", containerName, err)
	} else {
		return nil
	}
}

func (runtime *Runtime) terminateIfRunning(container *docker.Container, wait time.Duration) {
	time.Sleep(wait)
	if runtime.isContainerRunning(container) {
		runtime.terminate(container)
	}
}

func (runtime *Runtime) registerInterruptTermination(container *docker.Container) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Infof("Interrupted. Terminating container %s", container.Name)
		runtime.terminateIfRunning(container, 0)
		os.Exit(130)
	}()
}

func (runtime *Runtime) terminate(container *docker.Container) {
	if err := runtime.client.StopContainer(container.ID, 0); err != nil {
		log.Errorf("Unable to stop container %s, %v", container.Name, err)
	}
}

func (runtime *Runtime) isContainerRunning(container *docker.Container) bool {
	if list, err := runtime.client.ListContainers(docker.ListContainersOptions{
		Filters: map[string][]string{"id": {container.ID}},
	}); err != nil {
		log.Errorf("Unable to list containers %v", err)
		return false
	} else {
		return len(list) > 0
	}
}

func LocalMount(dirName string, rw bool) docker.HostMount {
	workdir, _ := os.Getwd()
	path, _ := filepath.Abs(workdir)
	dir := filepath.Join(path, dirName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.Mkdir(dir, os.ModePerm)
	}
	return docker.HostMount{
		Source:   dir,
		Target:   fmt.Sprintf("/mnt/%s", dirName),
		Type:     "bind",
		ReadOnly: !rw,
	}
}

type RuntimeProvider interface {
	CreateRuntime() (*Runtime, error)
}
