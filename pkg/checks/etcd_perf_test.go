package checks

import (
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/container"
	"github.com/cyberark/conjur-inspect/pkg/test"
	"github.com/stretchr/testify/assert"
)

// mockContainerProvider for EtcdPerfCheck tests
type mockContainerProvider struct {
	execMap     map[string]mockExecResult
	containerID string
}

type mockExecResult struct {
	stdout io.Reader
	stderr io.Reader
	err    error
}

func (m *mockContainerProvider) Name() string                                   { return "MockProvider" }
func (m *mockContainerProvider) Info() (container.ContainerProviderInfo, error) { return nil, nil }
func (m *mockContainerProvider) Container(id string) container.Container {
	return &mockContainer{
		execMap:     m.execMap,
		containerID: m.containerID,
	}
}

type mockContainer struct {
	execMap     map[string]mockExecResult
	containerID string
}

func (m *mockContainer) ID() string                  { return m.containerID }
func (m *mockContainer) Inspect() (io.Reader, error) { return nil, nil }
func (m *mockContainer) Exec(args ...string) (io.Reader, io.Reader, error) {
	key := strings.Join(args, " ")
	res, ok := m.execMap[key]
	if !ok {
		return nil, nil, errors.New("not found")
	}
	return res.stdout, res.stderr, res.err
}
func (m *mockContainer) ExecAsUser(user string, args ...string) (io.Reader, io.Reader, error) {
	// For tests, treat ExecAsUser the same as Exec
	return m.Exec(args...)
}
func (m *mockContainer) Logs(since time.Duration) (io.Reader, error) { return nil, nil }

// helper to build SUT and run context
func newEtcdPerfCheck(execMap map[string]mockExecResult, containerID string) (EtcdPerfCheck, *check.RunContext) {
	provider := &mockContainerProvider{execMap: execMap, containerID: containerID}
	sut := EtcdPerfCheck{Provider: provider}
	runCtx := &check.RunContext{ContainerID: "mock", OutputStore: test.NewOutputStore()}
	return sut, runCtx
}

func TestEtcdPerfCheck_Run_Success(t *testing.T) {
	execMap := map[string]mockExecResult{
		"echo":             {},
		"which etcd":       {},
		"which etcdctl":    {},
		"sv status conjur": {stdout: strings.NewReader("down: conjur")},
		"sv status pg":     {stdout: strings.NewReader("down: pg")},
		"sv status etcd":   {stdout: strings.NewReader("down: etcd")},
		"pgrep etcd":       {err: errors.New("not running")},
		"rm -rf /var/lib/conjur/etcd_performance_test":   {},
		"mkdir -p /var/lib/conjur/etcd_performance_test": {},
		"sh -c ETCD_DATA_DIR=/var/lib/conjur/etcd_performance_test ETCD_DEBUG=true etcd >/var/lib/conjur/etcd_performance_test/server.log 2>&1 & echo $!": {stdout: strings.NewReader("123")},
		"curl -s -o /dev/null -w %{http_code} http://127.0.0.1:2379/health":                                                                               {stdout: strings.NewReader("200")},
		"env ETCDCTL_API=3 etcdctl check perf --prefix /etcdctl-check-perf/":                                                                              {stdout: strings.NewReader("PASS: perf test passed\r\nFAIL: another passed")},
		"cat /var/lib/conjur/etcd_performance_test/server.log":                                                                                            {stdout: strings.NewReader("etcd log")},
		"kill -HUP 123": {},
	}
	sut, runCtx := newEtcdPerfCheck(execMap, "mock")
	results := sut.Run(runCtx)
	assert.Len(t, results, 2)
	assert.Equal(t, check.StatusPass, results[0].Status)
	assert.Equal(t, "GOOD", results[0].Value)
	assert.Equal(t, check.StatusFail, results[1].Status)
	assert.Equal(t, "BAD", results[1].Value)
}

func TestEtcdPerfCheck_Run_ValidationError(t *testing.T) {
	execMap := map[string]mockExecResult{
		"echo": {err: errors.New("fail")},
	}
	sut, runCtx := newEtcdPerfCheck(execMap, "mock")
	runCtx.VerboseErrors = true
	results := sut.Run(runCtx)
	assert.Len(t, results, 1)
	assert.Equal(t, check.StatusError, results[0].Status)
}

func TestEtcdPerfCheck_Run_NoContainerId(t *testing.T) {
	execMap := map[string]mockExecResult{
		"echo": {err: errors.New("fail")},
	}
	sut, runCtx := newEtcdPerfCheck(execMap, "")
	runCtx.VerboseErrors = true
	results := sut.Run(runCtx)
	assert.Len(t, results, 1)
	assert.Equal(t, check.StatusError, results[0].Status)
	assert.Equal(t, "N/A", results[0].Value)
}

func TestEtcdPerfCheck_Run_ServiceRunning(t *testing.T) {
	execMap := map[string]mockExecResult{
		"echo":             {},
		"which etcd":       {},
		"which etcdctl":    {},
		"sv status conjur": {stdout: strings.NewReader("run: conjur")},
	}
	sut, runCtx := newEtcdPerfCheck(execMap, "mock")
	runCtx.VerboseErrors = true
	results := sut.Run(runCtx)
	assert.Len(t, results, 1)
	assert.Equal(t, check.StatusError, results[0].Status)
	assert.Contains(t, results[0].Message, "service is running")
}

func TestEtcdPerfCheck_Run_ServiceRunningNoVerboseErrors(t *testing.T) {
	execMap := map[string]mockExecResult{
		"echo":             {},
		"which etcd":       {},
		"which etcdctl":    {},
		"sv status conjur": {stdout: strings.NewReader("run: conjur")},
	}
	sut, runCtx := newEtcdPerfCheck(execMap, "mock")
	runCtx.VerboseErrors = false
	results := sut.Run(runCtx)
	// Should return empty since the validation error is suppressed
	assert.Len(t, results, 0)
}

func TestEtcdPerfCheck_parse(t *testing.T) {
	type exp struct {
		len      int
		statuses []string
		values   []string
		msgs     []string
	}
	tests := []struct {
		name   string
		input  string
		expect exp
	}{
		{
			name:  "ErrorOutput",
			input: "ERROR: something went wrong",
			expect: exp{
				len:      1,
				statuses: []string{check.StatusError},
				msgs:     []string{"something went wrong"},
			},
		},
		{
			name:  "ErrorOutputWithOtherLines",
			input: "PASS: test\r\nFAIL: test\r\nERROR: test",
			expect: exp{
				len:      1,
				statuses: []string{check.StatusError},
				msgs:     []string{"unexpected output"},
			},
		},
		{
			name:  "UnexpectedOutput",
			input: "unexpected output",
			expect: exp{
				len:      1,
				statuses: []string{check.StatusError},
				msgs:     []string{"unexpected output"},
			},
		},
		{
			name:  "SuccessMixedPassFail",
			input: "00%test_test$test\r\nPASS: test1\r\nFAIL: test2\r\nPASS: test3",
			expect: exp{
				len:      3,
				statuses: []string{check.StatusPass, check.StatusFail, check.StatusPass},
				values:   []string{"GOOD", "BAD", "GOOD"},
				msgs:     []string{"test1", "test2", "test3"},
			},
		},
		{
			name:  "MultipleFailIndicators",
			input: "PASS: ok\r\nSlowest request took too long: slow details here\r\nStddev too high: stddev details here\r\nFAIL: failure",
			expect: exp{
				len:      4,
				statuses: []string{check.StatusPass, check.StatusFail, check.StatusFail, check.StatusFail},
				values:   []string{"GOOD", "BAD", "BAD", "BAD"},
				msgs:     []string{"ok", "slow details here", "stddev details here", "failure"},
			},
		},
	}

	provider := &mockContainerProvider{}
	sut := EtcdPerfCheck{Provider: provider}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			results := sut.parse(tc.input)
			assert.Len(t, results, tc.expect.len)
			for i, st := range tc.expect.statuses {
				assert.Equal(t, st, results[i].Status)
			}
			for i, v := range tc.expect.values {
				assert.Contains(t, results[i].Value, v)
			}
			for i, sub := range tc.expect.msgs {
				assert.Contains(t, results[i].Message, sub)
			}
		})
	}
}
