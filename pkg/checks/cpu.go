package checks

import (
	"runtime"
	"strconv"

	"github.com/conjurinc/conjur-preflight/pkg/framework"
)

type Cpu struct {
}

// Describe provides a textual description of what this check gathers info on
func (*Cpu) Describe() string {
	return "CPU"
}

func (cpu *Cpu) Run() <-chan []framework.CheckResult {
	future := make(chan []framework.CheckResult)

	// TODO: Can we return avg recent utilization?
	go func() {
		future <- []framework.CheckResult{
			{
				Title:   "CPU Cores",
				Status:  framework.STATUS_INFO,
				Value:   strconv.Itoa(runtime.NumCPU()),
				Message: "",
			},
			{
				Title:  "CPU Architecture",
				Status: framework.STATUS_INFO,
				Value:  runtime.GOARCH,
			},
		}
	}() // async

	return future
}
