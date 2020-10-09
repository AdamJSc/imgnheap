package domain_test

import (
	"github.com/google/go-cmp/cmp"
	"imgnheap/service/domain"
	"os"
	"testing"
)

type mockFileInfo struct {
	os.FileInfo
	name string
}

func (m mockFileInfo) Name() string { return m.name }

func TestParseFileNameAndExtensionFromInfo(t *testing.T) {
	t.Run("parsing file name and extension from file info must provided expected result", func(t *testing.T) {
		testCases := []struct {
			input          string
			expectedOutput []string
		}{
			{input: "noExt", expectedOutput: []string{"noExt", ""}},
			{input: "no-ext", expectedOutput: []string{"no-ext", ""}},
			{input: "no_ext", expectedOutput: []string{"no_ext", ""}},
			{input: ".", expectedOutput: []string{"", ""}},
			{input: ".cotton", expectedOutput: []string{"", "cotton"}},
			{input: "juneBrownPlays.cotton", expectedOutput: []string{"juneBrownPlays", "cotton"}},
			{input: "june-brown-plays.cotton", expectedOutput: []string{"june-brown-plays", "cotton"}},
			{input: "june_brown_plays.cotton", expectedOutput: []string{"june_brown_plays", "cotton"}},
			{input: "june_brown.plays.cotton", expectedOutput: []string{"june_brown.plays", "cotton"}},
			{input: "june...brown...plays...cotton", expectedOutput: []string{"june...brown...plays..", "cotton"}},
		}

		for _, tc := range testCases {
			info := mockFileInfo{name: tc.input}

			fileName, ext := domain.ParseFileNameAndExtensionFromInfo(info)
			actualOutput := []string{fileName, ext}

			diff := cmp.Diff(tc.expectedOutput, actualOutput)
			if diff != "" {
				t.Fatalf("expected %+v, got %+v", tc.expectedOutput, actualOutput)
			}
		}
	})
}
