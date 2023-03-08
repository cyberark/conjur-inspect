package checks

import (
	"fmt"
	"net"
	"os"

	"github.com/cyberark/conjur-inspect/pkg/check"
)

// Follower implements a inspection check for the FOLLOWER environment variable.
type Follower struct {
}

// Describe provides a textual description of what this check gathers info on
func (*Follower) Describe() string {
	return "follower"
}

// LeaderPort is the port that the follower listens on
type LeaderPort struct {
	PortName string
	Port     string
	IsOpen   bool
}

// Run executes the check
func (follower *Follower) Run() <-chan []check.Result {
	// Create a channel to communicate with the check check
	future := make(chan []check.Result)

	go func() {
		hostname := os.Getenv("MASTER_HOSTNAME")

		if hostname == "" {
			future <- []check.Result{
				{
					Title:   "Leader Hostname",
					Status:  check.STATUS_ERROR,
					Value:   "N/A",
					Message: "Leader hostname is not set. Set the 'MASTER_HOSTNAME' environment variable to run this check",
				},
			}

			return
		}

		// Initialize ports
		leaderPorts := []LeaderPort{
			{
				PortName: "Leader API Port",
				Port:     "443",
			},
			{
				PortName: "Leader Replication Port",
				Port:     "5432",
			},
			{
				PortName: "Leader Audit Forwarding Port",
				Port:     "1999",
			},
		}

		// a slice (array) of all port reports
		results := []check.Result{}

		for _, leaderPort := range leaderPorts {
			result := check.Result{
				Title: leaderPort.PortName,
			}

			leaderPort, err := checkPort(hostname, &leaderPort)
			if err != nil {
				result.Status = check.STATUS_ERROR
				result.Value = "N/A"
				result.Message = err.Error()
				results = append(results, result)

				continue
			}

			result.Status = check.STATUS_INFO
			result.Value = net.JoinHostPort(hostname, leaderPort.Port)
			result.Message = fmt.Sprintf("Port: %s is open", leaderPort.Port)

			results = append(results, result)
		}

		future <- results
	}() // async

	return future
}

func checkPort(host string, leaderPort *LeaderPort) (*LeaderPort, error) {
	leaderPort.IsOpen = false

	url := fmt.Sprintf("%s:%s", host, leaderPort.Port)

	conn, err := net.Dial("tcp", url)
	if err != nil {
		return leaderPort, fmt.Errorf("connection falied on port: %s", leaderPort.Port)
	}

	conn.Close()

	leaderPort.IsOpen = true
	return leaderPort, nil
}
