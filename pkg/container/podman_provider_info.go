// Package container defines the providers for concrete container engines
// (e.g. Docker, Podman)
package container

import "github.com/cyberark/conjur-inspect/pkg/check"

// PodmanProviderInfo is the concrete implementation of ContainerProviderInfo
// for Podman
type PodmanProviderInfo struct {
	rawData []byte
	info    *PodmanInfo
}

// PodmanInfo contains the specific Podman runtime information for reporting
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

// Results returns the specific Podman runtime information for reporting
func (info *PodmanProviderInfo) Results() []check.Result {
	return []check.Result{
		podmanVersionResult(info.info),
		podmanDriverResult(info.info),
		podmanGraphRootResult(info.info),
		podmanRunRootResult(info.info),
		podmanVolumeRootResult(info.info),
	}
}

// RawData returns the raw JSON output from `podman info`
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
