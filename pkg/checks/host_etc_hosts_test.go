package checks

import (
	"testing"

	"github.com/cyberark/conjur-inspect/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHostEtcHostsDescribe(t *testing.T) {
	h := &HostEtcHosts{}
	assert.Equal(t, "Host /etc/hosts", h.Describe())
}

func TestHostEtcHostsRunSuccessful(t *testing.T) {
	// This test reads the actual /etc/hosts file on the system
	h := &HostEtcHosts{}
	runContext := test.NewRunContext("")
	results := h.Run(&runContext)

	// Should return empty results on success
	assert.Empty(t, results)

	// Verify the file was saved to the output store
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	require.Len(t, items, 1)
	info, err := items[0].Info()
	require.NoError(t, err)
	assert.Equal(t, "host-etc-hosts.txt", info.Name())
}

func TestHostEtcHostsRunWithVerboseErrors(t *testing.T) {
	// This test verifies that we return empty results even with verbose errors
	// when the file is readable (which /etc/hosts typically is)
	h := &HostEtcHosts{}
	runContext := test.NewRunContext("")
	runContext.VerboseErrors = true
	results := h.Run(&runContext)

	// Should still return empty results when file is readable
	assert.Empty(t, results)

	// Verify the file was saved to the output store
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	require.Len(t, items, 1)
}
