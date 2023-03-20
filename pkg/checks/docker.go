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

var executeDockerInfoFunc func() (stderr, stdout []byte, err error) = executeDockerInfo

// Docker collects the information on the version of Docker on the system
type Docker struct{}

// Describe provides a textual description of what this check gathers info on
func (*Docker) Describe() string {
	return "Docker runtime"
}

// Run performs the Docker inspection checks
func (*Docker) Run(context *check.RunContext) <-chan []check.Result {
	future := make(chan []check.Result)

	go func() {
		dockerInfo, err := getDockerInfo(context.OutputStore)
		if err != nil {
			future <- []check.Result{
				{
					Title:   "Docker",
					Status:  check.StatusError,
					Value:   "N/A",
					Message: err.Error(),
				},
			}

			return
		}

		if len(dockerInfo.ServerErrors) > 0 {
			future <- []check.Result{
				{
					Title:   "Docker",
					Status:  check.StatusError,
					Value:   "N/A",
					Message: strings.Join(dockerInfo.ServerErrors, ", "),
				},
			}

			return
		}

		future <- []check.Result{
			dockerVersionResult(dockerInfo),
			dockerDriverResult(dockerInfo),
			dockerRootDirResult(dockerInfo),
		}
	}() // async

	return future
}

// DockerInfo contains the key fields from the `docker info` output
type DockerInfo struct {
	ServerErrors  []string `json:"ServerErrors"`
	ServerVersion string   `json:"ServerVersion"`
	Driver        string   `json:"Driver"`
	DockerRootDir string   `json:"DockerRootDir"`
}

func getDockerInfo(outputStore output.Store) (*DockerInfo, error) {
	stdout, stderr, err := executeDockerInfoFunc()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to inspect Docker runtime: %w (%s)",
			err,
			strings.TrimSpace(string(stderr)),
		)
	}

	// Save raw docker info output
	outputReader := bytes.NewReader(stdout)
	err = outputStore.Save("docker-info.json", outputReader)
	if err != nil {
		log.Warn("Failed to save docker info output: %s", err)
	}

	// Parse the JSON output
	dockerInfo := DockerInfo{}
	err = json.Unmarshal(stdout, &dockerInfo)
	if err != nil {
		return nil, err
	}

	return &dockerInfo, nil
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
