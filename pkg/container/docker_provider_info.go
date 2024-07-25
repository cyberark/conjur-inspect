// Package container defines the providers for concrete container engines
// (e.g. Docker, Podman)
package container

import (
	"io"

	"github.com/cyberark/conjur-inspect/pkg/check"
)

// DockerProviderInfo is the concrete implementation of ContainerProviderInfo
// for Docker
type DockerProviderInfo struct {
	rawData io.Reader
	info    *DockerInfo
}

// DockerInfo contains the specific Docker runtime information for reporting
type DockerInfo struct {
	ServerErrors  []string `json:"ServerErrors"`
	ServerVersion string   `json:"ServerVersion"`
	Driver        string   `json:"Driver"`
	DockerRootDir string   `json:"DockerRootDir"`
}

// Results returns the specific Docker runtime information for reporting
func (info *DockerProviderInfo) Results() []check.Result {
	return []check.Result{
		dockerVersionResult(info.info),
		dockerDriverResult(info.info),
		dockerRootDirResult(info.info),
	}
}

// RawData returns the raw JSON output from `docker info`
func (info *DockerProviderInfo) RawData() io.Reader {
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
