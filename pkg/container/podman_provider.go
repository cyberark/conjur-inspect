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
var executePodmanInfoFunc = executePodmanInfo

// PodmanProvider is a concrete implementation of the
// ContainerProvider interface for Podman
type PodmanProvider struct{}

type PodmanProviderInfo struct {
	rawData []byte
	info    *PodmanInfo
}

type PodmanInfo struct {
	Version VersionInfo `json:"version"`
	Store   StoreInfo   `json:"store"`
}

// VersionInfo contains the Podman version information
type VersionInfo struct {
	Version string
}

// StoreInfo contains the Podman storage information
type StoreInfo struct {
	GraphDriverName string `json:"graphDriverName"`
	GraphRoot       string `json:"graphRoot"`
	RunRoot         string `json:"runRoot"`
	VolumePath      string `json:"volumePath"`
}

// Name returns the name of the Podman provider
func (*PodmanProvider) Name() string {
	return "Podman"
}

// Info returns the Podman runtime info
func (*PodmanProvider) Info() (ContainerProviderInfo, error) {
	stdout, stderr, err := executePodmanInfoFunc()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to inspect Podman runtime: %w (%s)",
			err,
			strings.TrimSpace(string(stderr)),
		)
	}

	// Parse the JSON output
	podmanInfo := &PodmanInfo{}
	err = json.Unmarshal(stdout, podmanInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Podman info output: %w", err)
	}

	podmanProviderInfo := &PodmanProviderInfo{
		rawData: stdout,
		info:    podmanInfo,
	}

	return podmanProviderInfo, nil
}

func (info *PodmanProviderInfo) Results() []check.Result {
	return []check.Result{
		podmanVersionResult(info.info),
		podmanDriverResult(info.info),
		podmanGraphRootResult(info.info),
		podmanRunRootResult(info.info),
		podmanVolumeRootResult(info.info),
	}
}

func (info *PodmanProviderInfo) RawData() []byte {
	return info.rawData
}

func podmanVersionResult(podmanInfo *PodmanInfo) check.Result {
	return check.Result{
		Title:  "Podman Version",
		Status: check.StatusInfo,
		Value:  podmanInfo.Version.Version,
	}
}

func podmanDriverResult(podmanInfo *PodmanInfo) check.Result {
	return check.Result{
		Title:  "Podman Driver",
		Status: check.StatusInfo,
		Value:  podmanInfo.Store.GraphDriverName,
	}
}

func podmanGraphRootResult(podmanInfo *PodmanInfo) check.Result {
	return check.Result{
		Title:  "Podman Graph Root",
		Status: check.StatusInfo,
		Value:  podmanInfo.Store.GraphRoot,
	}
}

func podmanRunRootResult(podmanInfo *PodmanInfo) check.Result {
	return check.Result{
		Title:  "Podman Run Root",
		Status: check.StatusInfo,
		Value:  podmanInfo.Store.RunRoot,
	}
}

func podmanVolumeRootResult(podmanInfo *PodmanInfo) check.Result {
	return check.Result{
		Title:  "Podman Volume Path",
		Status: check.StatusInfo,
		Value:  podmanInfo.Store.VolumePath,
	}
}

func executePodmanInfo() (stdout, stderr []byte, err error) {
	return shell.NewCommandWrapper(
		"podman",
		"info",
		"--debug",
		"--format",
		"{{json .}}",
	).Run()
}
