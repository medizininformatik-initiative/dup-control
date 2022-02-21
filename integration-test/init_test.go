//go:build integration

package test

import (
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"os"
	"testing"
)

func cleanup() {
	if err := os.RemoveAll("outputGlobal"); err != nil {
		fmt.Println(err)
	}
	if err := os.RemoveAll("outputLocal"); err != nil {
		fmt.Println(err)
	}
}

func TestMain(m *testing.M) {
	cleanup()
	code := m.Run()
	os.Exit(code)
}

var (
	d, _ = docker.NewClientFromEnv()
)
