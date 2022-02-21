package testing

import (
	docker "github.com/fsouza/go-dockerclient"
	"strings"
)

func FindContainerByName(d *docker.Client, name string) *docker.APIContainers {
	containers, _ := d.ListContainers(docker.ListContainersOptions{})
	for _, c := range containers {
		for _, n := range c.Names {
			if strings.Contains(n, name) {
				return &c
			}
		}
	}
	return nil
}

func KillContainer(d *docker.Client, c *docker.APIContainers) {
	if c != nil {
		_ = d.KillContainer(docker.KillContainerOptions{ID: c.ID})
	}
}

func RemoveImage(d *docker.Client, name string) {
	_ = d.RemoveImage(name)
}
