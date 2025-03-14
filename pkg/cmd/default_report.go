// Package cmd is the entry point for the conjur-inspect command line tool.
package cmd

import (
	"os"
	"path"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/checks"
	"github.com/cyberark/conjur-inspect/pkg/checks/disk"
	"github.com/cyberark/conjur-inspect/pkg/container"
	"github.com/cyberark/conjur-inspect/pkg/output"
	"github.com/cyberark/conjur-inspect/pkg/report"
	"github.com/cyberark/conjur-inspect/pkg/reports"
)

// NewDefaultReport returns a report containing the standard inspection checks
func NewDefaultReport(
	id string,
	rawDataDir string,
) (report.Report, error) {

	storeDirectory := path.Join(rawDataDir, id)

	err := os.MkdirAll(storeDirectory, 0755)
	if err != nil {
		return nil, err
	}

	outputStore := output.NewDirectoryStore(storeDirectory)
	outputArchive := &output.TarGzipArchive{OutputDir: rawDataDir}

	return reports.NewStandardReport(
		id,
		defaultReportSections(),
		outputStore,
		outputArchive,
	), nil
}

func defaultReportSections() []report.Section {
	return []report.Section{
		// TODO:
		// - Recent load
		{
			Title: "CPU",
			Checks: []check.Check{
				&checks.Cpu{},
			},
		},
		{
			Title: "Disk",
			Checks: []check.Check{
				&disk.SpaceCheck{},
				disk.NewIopsCheck(),
				disk.NewLatencyCheck(),
			},
		},
		{
			Title: "Memory",
			Checks: []check.Check{
				&checks.Memory{},
			},
		},
		// TODO:
		// - ipv6 status
		// {
		// 	Title:  "Network",
		// 	Checks: []framework.Check{},
		// },
		{
			Title: "Host",
			Checks: []check.Check{
				&checks.Host{},
			},
		},
		{
			Title: "Follower",
			Checks: []check.Check{
				&checks.Follower{},
			},
		},
		{
			Title: "Container",
			Checks: []check.Check{
				// Runtime
				&checks.ContainerRuntime{
					Provider: &container.DockerProvider{},
				},
				&checks.ContainerRuntime{
					Provider: &container.PodmanProvider{},
				},

				// Container inspect
				&checks.ContainerInspect{
					Provider: &container.DockerProvider{},
				},
				&checks.ContainerInspect{
					Provider: &container.PodmanProvider{},
				},

				// Container logs
				&checks.ContainerLogs{
					Provider: &container.DockerProvider{},
				},
				&checks.ContainerLogs{
					Provider: &container.PodmanProvider{},
				},

				// Container config
				&checks.ConjurConfig{
					Provider: &container.DockerProvider{},
				},
				&checks.ConjurConfig{
					Provider: &container.PodmanProvider{},
				},

				// Container config
				&checks.ConjurConfigPermissions{
					Provider: &container.DockerProvider{},
				},
				&checks.ConjurConfigPermissions{
					Provider: &container.PodmanProvider{},
				},
			},
		},
		{
			Title: "Conjur",
			Checks: []check.Check{
				// Health
				&checks.ConjurHealth{
					Provider: &container.DockerProvider{},
				},
				&checks.ConjurHealth{
					Provider: &container.PodmanProvider{},
				},

				// Info
				&checks.ConjurInfo{
					Provider: &container.DockerProvider{},
				},
				&checks.ConjurInfo{
					Provider: &container.PodmanProvider{},
				},
			},
		},
		{
			Title: "Ulimits",
			Checks: []check.Check{
				&checks.Ulimit{},
			},
		},
	}
}
