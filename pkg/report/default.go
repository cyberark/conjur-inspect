package report

import (
	"github.com/conjurinc/conjur-preflight/pkg/checks"
	"github.com/conjurinc/conjur-preflight/pkg/framework"
)

func NewDefaultReport() framework.Report {
	return framework.Report{
		Sections: []framework.ReportSection{
			// TODO:
			// - Recent load
			{
				Title: "CPU",
				Checks: []framework.Check{
					&checks.Cpu{},
				},
			},
			// TODO
			// - IOPS
			{
				Title: "Disk",
				Checks: []framework.Check{
					&checks.DiskSpace{},
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
