package checks

import (
	"errors"
	"testing"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestDockerRunSuccess(t *testing.T) {
	// Mock dependencies
	oldFunc := executeDockerInfoFunc
	executeDockerInfoFunc = func() (stdout, stderr []byte, err error) {
		stdout = []byte(`{"ServerVersion":"20.10.7","Driver":"overlay2","DockerRootDir":"/var/lib/docker"}`)
		return stdout, stderr, err
	}
	defer func() {
		executeDockerInfoFunc = oldFunc
	}()

	// Run the check
	docker := &Docker{}
	context := test.NewRunContext()
	results := <-docker.Run(&context)

	// Check the results
	expected := []check.Result{
		{
			Title:  "Docker Version",
			Status: check.StatusInfo,
			Value:  "20.10.7",
		},
		{
			Title:  "Docker Driver",
			Status: check.StatusInfo,
			Value:  "overlay2",
		},
		{
			Title:  "Docker Root Directory",
			Status: check.StatusInfo,
			Value:  "/var/lib/docker",
		},
	}
	assert.Equal(t, expected, results)
}

func TestDockerRunFailure(t *testing.T) {
	// Mock dependencies
	oldFunc := executeDockerInfoFunc
	executeDockerInfoFunc = func() (stdout, stderr []byte, err error) {
		err = errors.New("fake error")
		return stdout, stderr, err
	}
	defer func() {
		executeDockerInfoFunc = oldFunc
	}()

	// Run the check
	docker := &Docker{}
	context := test.NewRunContext()
	results := <-docker.Run(&context)

	// Check the results
	expected := []check.Result{
		{
			Title:   "Docker",
			Status:  check.StatusError,
			Value:   "N/A",
			Message: "failed to inspect Docker runtime: fake error ()",
		},
	}
	assert.Equal(t, expected, results)
}

func TestDockerRunServerError(t *testing.T) {
	// Mock dependencies
	oldFunc := executeDockerInfoFunc
	executeDockerInfoFunc = func() (stdout, stderr []byte, err error) {
		stdout = []byte(`{"ServerErrors": ["Test error"]}`)
		return stdout, stderr, err
	}
	defer func() {
		executeDockerInfoFunc = oldFunc
	}()

	// Run the check
	docker := &Docker{}
	context := test.NewRunContext()
	results := <-docker.Run(&context)

	// Check the results
	expected := []check.Result{
		{
			Title:   "Docker",
			Status:  check.StatusError,
			Value:   "N/A",
			Message: "Test error",
		},
	}
	assert.Equal(t, expected, results)
}
