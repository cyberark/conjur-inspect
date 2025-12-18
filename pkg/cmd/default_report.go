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
				&checks.CommandHistory{},
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
				// Check container runtime availability first to cache the results
				&checks.ContainerAvailability{},

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

				// Container command history
				&checks.ContainerCommandHistory{
					Provider: &container.DockerProvider{},
				},
				&checks.ContainerCommandHistory{
					Provider: &container.PodmanProvider{},
				},

				// Container processes
				&checks.ContainerProcesses{
					Provider: &container.DockerProvider{},
				},
				&checks.ContainerProcesses{
					Provider: &container.PodmanProvider{},
				},

				// Container top
				&checks.ContainerTop{
					Provider: &container.DockerProvider{},
				},
				&checks.ContainerTop{
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

				// Runit services
				&checks.RunItServices{
					Provider: &container.DockerProvider{},
				},
				&checks.RunItServices{
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

				// Ruby thread dumps
				&checks.RubyThreadDump{
					Provider: &container.DockerProvider{},
				},
				&checks.RubyThreadDump{
					Provider: &container.PodmanProvider{},
				},

				// PostgreSQL pg_stat_activity
				&checks.PgStatActivity{
					Provider: &container.DockerProvider{},
				},
				&checks.PgStatActivity{
					Provider: &container.PodmanProvider{},
				},
			},
		},
		{
			Title: "Etcd",
			Checks: []check.Check{
				// Etcd Perf
				&checks.EtcdPerfCheck{
					Provider: &container.DockerProvider{},
				},
				&checks.EtcdPerfCheck{
					Provider: &container.PodmanProvider{},
				},

				// Etcd Cluster Members
				&checks.EtcdClusterMembers{
					Provider: &container.DockerProvider{},
				},
				&checks.EtcdClusterMembers{
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
