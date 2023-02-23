package report

import (
	"github.com/conjurinc/conjur-preflight/pkg/checks"
	"github.com/conjurinc/conjur-preflight/pkg/checks/disk"
	"github.com/conjurinc/conjur-preflight/pkg/framework"
)

// NewDefaultReport returns a report containing the standard pre-flight checks
func NewDefaultReport(debug bool) Report {
	return Report{
		Sections: []Section{
			// TODO:
			// - Recent load
			{
				Title: "CPU",
				Checks: []framework.Check{
					&checks.Cpu{},
				},
			},
			{
				Title: "Disk",
				Checks: []framework.Check{
					&disk.SpaceCheck{},
					disk.NewIopsCheck(debug),
					disk.NewLatencyCheck(debug),
				},
			},
			{
				Title: "Memory",
				Checks: []framework.Check{
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
				Checks: []framework.Check{
					&checks.Host{},
				},
			},
			{
				Title: "Follower",
				Checks: []framework.Check{
					&checks.Follower{},
				},
			},
			// TODO:
			// - Podman version
			// - Docker version
			// {
			// 	Title:  "Container Runtime",
			// 	Checks: []framework.Check{},
			// },
		},
	}
}
