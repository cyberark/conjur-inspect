// Package test defines utilities and mock implementations for testing
package test

import (
	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/container"
)

// ContainerProvider is a mock implementation of the ContainerProvider interface
// for testing
type ContainerProvider struct {
	InfoError   error
	InfoRawData []byte
	InfoResults []check.Result
}

// ContainerProviderInfo is a mock implementation of the ContainerProviderInfo
// interface for testing
type ContainerProviderInfo struct {
	InfoRawData []byte
	InfoResults []check.Result
}

// Container is a mock implementation of the Container interface for testing
type Container struct {
	ContainerID string
}

// Name returns the name of the container provider
func (provider *ContainerProvider) Name() string {
	return "Test Container Provider"
}

// Info returns the container provider info
func (provider *ContainerProvider) Info() (container.ContainerProviderInfo, error) {
	if provider.InfoError != nil {
		return nil, provider.InfoError
	}

	return &ContainerProviderInfo{
		InfoRawData: provider.InfoRawData,
		InfoResults: provider.InfoResults,
	}, nil
}

// Container returns a container instance for the given ID
func (provider *ContainerProvider) Container(
	containerID string,
) container.Container {
	return &Container{
		ContainerID: containerID,
	}
}

// Results returns the check results
func (providerInfo *ContainerProviderInfo) Results() []check.Result {
	return providerInfo.InfoResults
}

// RawData returns the raw data
func (providerInfo *ContainerProviderInfo) RawData() []byte {
	return providerInfo.InfoRawData
}

// ID returns the container ID
func (container *Container) ID() string {
	return container.ContainerID
}
