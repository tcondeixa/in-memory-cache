package cache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {

	for _, tc := range []struct {
		testName       string
		errorType      string
		expectedOutput string
	}{
		{
			testName:       "InvalidCacheKeyError",
			errorType:      (&InvalidCacheKeyError{}).Error(),
			expectedOutput: invalidCacheKeyError,
		},
		{
			testName:       "OutdatedCacheEntryError",
			errorType:      (&OutdatedCacheEntryError{}).Error(),
			expectedOutput: outdatedCacheEntryError,
		},
		{
			testName:       "KeyAlreadyExistsError",
			errorType:      (&KeyAlreadyExistsError{}).Error(),
			expectedOutput: keyAlreadyExistsError,
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			require.Equal(t, tc.expectedOutput, tc.errorType)
		})
	}

}
