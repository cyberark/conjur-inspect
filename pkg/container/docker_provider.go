// Package container defines the providers for concrete container engines
// (e.g. Docker, Podman)
package container

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/shell"
)

// Function variable for dependency injection
var executeDockerInfoFunc = executeDockerInfo

// DockerProvider is a concrete implementation of the
// ContainerProvider interface for Docker
type DockerProvider struct {
}

type DockerProviderInfo struct {
	rawData []byte
	info    *DockerInfo
}

type DockerInfo struct {
	ServerErrors  []string `json:"ServerErrors"`
	ServerVersion string   `json:"ServerVersion"`
	Driver        string   `json:"Driver"`
	DockerRootDir string   `json:"DockerRootDir"`
}

// Name returns the name of the Docker provider
func (*DockerProvider) Name() string {
	return "Docker"
}

// Info returns the Docker runtime info
func (*DockerProvider) Info() (ContainerProviderInfo, error) {
	stdout, stderr, err := executeDockerInfoFunc()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to inspect Docker runtime: %w (%s)",
			err,
			strings.TrimSpace(string(stderr)),
		)
	}

	// Parse the JSON output
	dockerInfo := &DockerInfo{}
	err = json.Unmarshal(stdout, dockerInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Docker info output: %w", err)
	}

	// Check for server errors
	if len(dockerInfo.ServerErrors) > 0 {
		return nil, fmt.Errorf(
			"Docker runtime has server errors: %s",
			strings.Join(dockerInfo.ServerErrors, ", "),
		)
	}

	dockerProviderInfo := &DockerProviderInfo{
		rawData: stdout,
		info:    dockerInfo,
	}

	return dockerProviderInfo, nil
}

func (info *DockerProviderInfo) Results() []check.Result {
	return []check.Result{
		dockerVersionResult(info.info),
		dockerDriverResult(info.info),
		dockerRootDirResult(info.info),
	}
}

func (info *DockerProviderInfo) RawData() []byte {
	return info.rawData
}

func dockerVersionResult(dockerInfo *DockerInfo) check.Result {
	return check.Result{
		Title:  "Docker Version",
		Status: check.StatusInfo,
		Value:  dockerInfo.ServerVersion,
	}
}

func dockerDriverResult(dockerInfo *DockerInfo) check.Result {
	return check.Result{
		Title:  "Docker Driver",
		Status: check.StatusInfo,
		Value:  dockerInfo.Driver,
	}
}

func dockerRootDirResult(dockerInfo *DockerInfo) check.Result {
	return check.Result{
		Title:  "Docker Root Directory",
		Status: check.StatusInfo,
		Value:  dockerInfo.DockerRootDir,
	}
}

func executeDockerInfo() (stdout, stderr []byte, err error) {
	return shell.NewCommandWrapper(
		"docker",
		"--debug",
		"info",
		"--format",
		"{{json .}}",
	).Run()
}
