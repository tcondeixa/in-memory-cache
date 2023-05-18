package cache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileType(t *testing.T) {

	for _, tc := range []struct {
		testName       string
		fileFormat     FileFormat
		expectedOutput string
	}{
		{
			testName:       "json format",
			fileFormat:     FileFormat(Json),
			expectedOutput: jsonOutput,
		},
		{
			testName:       "yaml format",
			fileFormat:     FileFormat(Yaml),
			expectedOutput: yamlOutput,
		},
		{
			testName:       "not supported format",
			fileFormat:     FileFormat(5),
			expectedOutput: unknownOutput,
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			require.Equal(t, tc.expectedOutput, tc.fileFormat.String())
		})
	}

}
