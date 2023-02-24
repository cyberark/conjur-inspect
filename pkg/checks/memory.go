package checks

import (
	"fmt"

	"github.com/conjurinc/conjur-preflight/pkg/framework"
	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/v3/mem"
)

type Memory struct {
}

// Describe provides a textual description of what this check gathers info on
func (*Memory) Describe() string {
	return "memory"
}

func (memory *Memory) Run() <-chan []framework.CheckResult {
	future := make(chan []framework.CheckResult)

	go func() {

		v, _ := mem.VirtualMemory()

		future <- []framework.CheckResult{
			{
				Title:  "Memory Total",
				Status: framework.STATUS_INFO,
				Value:  humanize.Bytes(v.Total),
			},
			{
				Title:  "Memory Free",
				Status: framework.STATUS_INFO,
				Value:  humanize.Bytes(v.Free),
			},
			{
				Title:  "Memory Used",
				Status: framework.STATUS_INFO,
				Value: fmt.Sprintf(
					"%s (%.1f %%)",
					humanize.Bytes(v.Used),
					v.UsedPercent,
				),
			},
		}
	}() // async

	return future
}
