package maybe_test

import (
	"errors"
	"testing"

	"github.com/conjurinc/conjur-preflight/pkg/maybe"
	"github.com/stretchr/testify/assert"
)

func TestSuccess(t *testing.T) {
	value := "test"

	// Construct a new successful Maybe instance
	success := maybe.NewSuccess(value)

	// It's error should be nil
	assert.Nil(t, success.Error())

	// It should hold the successful value
	assert.Equal(t, value, success.Value())

	// If we use the accessor with an error, it should match
	// the above state.
	returnedValue, returnedErr := success.ValueE()
	assert.Nil(t, returnedErr)
	assert.Equal(t, value, returnedValue)
}

func TestFailure(t *testing.T) {
	err := errors.New("Test error")

	// Construct a new failed Maybe instance
	failure := maybe.NewFailure[string](err)

	// It should return the error
	assert.Equal(t, err, failure.Error())

	// It should return the default value for the Maybe type
	assert.Equal(t, "", failure.Value())

	// If we use the accessor with an error, it should fail with an error
	returnedValue, returnedErr := failure.ValueE()
	assert.Equal(t, maybe.ErrorNoValue, returnedErr)
	assert.Equal(t, "", returnedValue)

}

func TestResult(t *testing.T) {
	value := "test"
	err := errors.New("Test error")

	// With an error
	result := maybe.Result(value, err)

	// It should return the error
	assert.Equal(t, err, result.Error())

	// It should return the default value for the Maybe type
	assert.Equal(t, "", result.Value())

	// Without an error
	result = maybe.Result(value, nil)

	// It's error should be nil
	assert.Nil(t, result.Error())

	// It should hold the successful value
	assert.Equal(t, value, result.Value())
}

func TestBind(t *testing.T) {
	var startingMaybe maybe.Maybe[string] = maybe.NewSuccess("ping")

	// When the bind is successful
	resultMaybe := maybe.Bind(
		startingMaybe,
		func(in string) (string, error) {
			return "pong", nil
		},
	)

	// It's error should be nil
	assert.Nil(t, resultMaybe.Error())

	// It should hold the successful value
	assert.Equal(t, "pong", resultMaybe.Value())

	// When the bind fails
	resultMaybe = maybe.Bind(
		startingMaybe,
		func(in string) (string, error) {
			return "", errors.New("pong error")
		},
	)

	// It should return the error
	assert.Error(t, resultMaybe.Error(), "pong error")

	// It should return the default value for the Maybe type
	assert.Equal(t, "", resultMaybe.Value())
}
