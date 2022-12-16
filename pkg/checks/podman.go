package checks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/log"
	"github.com/cyberark/conjur-inspect/pkg/output"
	"github.com/cyberark/conjur-inspect/pkg/shell"
)

var executePodmanInfoFunc func() (stderr, stdout []byte, err error) = executePodmanInfo

// Podman collects the information on the version of Podman on the system
type Podman struct{}

// Describe provides a textual description of what this check gathers info on
func (*Podman) Describe() string {
	return "Podman runtime"
}

// Run performs the Podman inspection checks
func (podman *Podman) Run(context *check.RunContext) <-chan []check.Result {
	future := make(chan []check.Result)

	go func() {
		podmanInfo, err := getPodmanInfo(context.OutputStore)
		if err != nil {
			future <- []check.Result{
				{
					Title:   "Podman",
					Status:  check.StatusError,
					Value:   "N/A",
					Message: err.Error(),
				},
			}

			return
		}

		future <- []check.Result{
			podmanVersionResult(podmanInfo),
			podmanDriverResult(podmanInfo),
			podmanGraphRootResult(podmanInfo),
			podmanRunRootResult(podmanInfo),
			podmanVolumeRootResult(podmanInfo),
		}
	}() // async

	return future
}

// PodmanInfo contains the key fields from the `podman info` output
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

func getPodmanInfo(outputStore output.Store) (*PodmanInfo, error) {
	stdout, stderr, err := executePodmanInfoFunc()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to inspect Podman runtime: %w (%s)",
			err,
			strings.TrimSpace(string(stderr)),
		)
	}

	// Save raw podman info output
	outputReader := bytes.NewReader(stdout)
	err = outputStore.Save("podman-info.json", outputReader)
	if err != nil {
		log.Warn("Failed to save podman info output: %s", err)
	}

	// Parse the podman info output
	podmanInfo := PodmanInfo{}
	err = json.Unmarshal(stdout, &podmanInfo)

	return &podmanInfo, err
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
