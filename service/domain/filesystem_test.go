package domain_test

import (
	"github.com/google/go-cmp/cmp"
	"imgnheap/service/domain"
	"imgnheap/service/models"
	"testing"
	"time"
)

var fileNamesContainingParseableTimestamp = []string{
	"20180526140029",
	"20180526_140029",
	"20180526-140029",
	"Screenshot_20180526140029",
	"Screenshot_20180526_140029",
	"Screenshot_20180526-140029",
	"Screenshot_20180526140029_MyFaceSpace",
	"Screenshot_20180526_140029_MyFaceSpace",
	"Screenshot_20180526-140029_MyFaceSpace",
	"Screenshot 2018-05-26 at 14.00.29",
}

var fileNamesContainingNoParseableTimestamp = []string{
	"hello-world",
	"12345678",
	"true",

	// all garbage date strings!
	"2006660102150405",
	"2006011102150405",
	"2006010222150405",
	"2006010215550405",
	"2006010215044405",
	"2006010215040555",
	"2006660102_150405",
	"2006011102_150405",
	"2006010222_150405",
	"20060102_15550405",
	"20060102_15044405",
	"20060102_15040555",
	"2006660102-150405",
	"2006011102-150405",
	"2006010222-150405",
	"20060102-15550405",
	"20060102-15044405",
	"20060102-15040555",
	"Screenshot_2006660102150405",
	"Screenshot_2006011102150405",
	"Screenshot_2006010222150405",
	"Screenshot_2006010215550405",
	"Screenshot_2006010215044405",
	"Screenshot_2006010215040555",
	"Screenshot_2006660102_150405",
	"Screenshot_2006011102_150405",
	"Screenshot_2006010222_150405",
	"Screenshot_20060102_15550405",
	"Screenshot_20060102_15044405",
	"Screenshot_20060102_15040555",
	"Screenshot_2006660102-150405",
	"Screenshot_2006011102-150405",
	"Screenshot_2006010222-150405",
	"Screenshot_20060102-15550405",
	"Screenshot_20060102-15044405",
	"Screenshot_20060102-15040555",
	"Screenshot 200666-01-02 at 15.04.05",
	"Screenshot 2006-0111-02 at 15.04.05",
	"Screenshot 2006-01-0222 at 15.04.05",
	"Screenshot 2006-01-02 at 1555.04.05",
	"Screenshot 2006-01-02 at 15.0444.05",
	"Screenshot 2006-01-02 at 15.04.0555",
}

func TestParseNameAndExtensionFromFileName(t *testing.T) {
	t.Run("parsing name and extension from filename must provided expected result", func(t *testing.T) {
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

		for idx, tc := range testCases {
			fileName, ext := domain.ParseNameAndExtensionFromFileName(tc.input)
			actualOutput := []string{fileName, ext}

			diff := cmp.Diff(tc.expectedOutput, actualOutput)
			if diff != "" {
				t.Fatalf("tc %d: expected %+v, got %+v", idx, tc.expectedOutput, actualOutput)
			}
		}
	})
}

func TestParseTimestampFromFile(t *testing.T) {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		t.Fatal(err)
	}

	expectedTs := time.Date(2018, 5, 26, 14, 0, 29, 0, loc)

	t.Run("parsing timestamp from file that has a compatible filename must provide timestamp parsed from filename", func(t *testing.T) {
		testCases := fileNamesContainingParseableTimestamp

		for idx, tc := range testCases {
			file := models.File{Name: tc}
			actualTs := domain.ParseTimestampFromFile(file)

			if !expectedTs.Equal(actualTs) {
				t.Fatalf("tc %d: expected %+v, got %+v", idx, expectedTs, actualTs)
			}
		}
	})

	t.Run("parsing timestamp from file that does not have a compatible filename must provide file created at timestamp", func(t *testing.T) {
		testCases := fileNamesContainingNoParseableTimestamp

		for idx, tc := range testCases {
			file := models.File{Name: tc, CreatedAt: expectedTs}
			actualTs := domain.ParseTimestampFromFile(file)

			if !expectedTs.Equal(actualTs) {
				t.Fatalf("tc %d: expected %+v, got %+v", idx, expectedTs, actualTs)
			}
		}
	})
}

func TestGetDestinationDirWithTimestamp(t *testing.T) {
	sess := models.Session{
		BaseDir: "/base/dir",
		SubDir:  "subdir",
	}

	t.Run("get destination dir with timestamp using a timestamp-parseable filename must return the expected result", func(t *testing.T) {
		testCases := fileNamesContainingParseableTimestamp

		expectedOutput := "/base/dir/subdir/2018-05-26"

		for idx, tc := range testCases {
			file := models.NewFile(tc, "jpg", "/base/dir", time.Time{})

			destDir := domain.GetDestinationDirWithTimestamp(file, &sess)
			if destDir != expectedOutput {
				t.Fatalf("tc %d: expected %s, got %s", idx, expectedOutput, destDir)
			}
		}
	})

	t.Run("get destination dir with timestamp using a non-timestamp-parseable filename must return the expected result", func(t *testing.T) {
		testCases := fileNamesContainingNoParseableTimestamp

		expectedOutput := "/base/dir/subdir/0001-01-01"

		for idx, tc := range testCases {
			file := models.NewFile(tc, "jpg", "/base/dir", time.Time{})

			destDir := domain.GetDestinationDirWithTimestamp(file, &sess)
			if destDir != expectedOutput {
				t.Fatalf("tc %d: expected %s, got %s", idx, expectedOutput, destDir)
			}
		}
	})

	t.Run("get destination dir with timestamp using nil session must return blank string", func(t *testing.T) {
		file := models.NewFile("hello_world", "jpg", "/base/dir", time.Time{})

		destDir := domain.GetDestinationDirWithTimestamp(file, nil)
		if destDir != "" {
			t.Fatalf("expected empty string, got %s", destDir)
		}
	})
}

func TestGetDestinationDirWithTag(t *testing.T) {
	sess := models.Session{
		BaseDir: "/base/dir",
		SubDir:  "subdir",
	}

	t.Run("get destination dir with tag must return the expected result", func(t *testing.T) {
		testCases := []struct {
			input          string
			expectedOutput string
		}{
			{input: "helloWorld", expectedOutput: "/base/dir/subdir/helloWorld"},
			{input: "hello-world", expectedOutput: "/base/dir/subdir/hello-world"},
			{input: "hello_world", expectedOutput: "/base/dir/subdir/hello_world"},
			{input: "hello/world", expectedOutput: "/base/dir/subdir/hello/world"},
			{input: "hello/world//", expectedOutput: "/base/dir/subdir/hello/world"},
		}

		for idx, tc := range testCases {
			destDir := domain.GetDestinationDirWithTag(&sess, tc.input)
			if destDir != tc.expectedOutput {
				t.Fatalf("tc %d: expected %s, got %s", idx, tc.expectedOutput, destDir)
			}
		}
	})

	t.Run("get destination dir with tag using nil session must return blank string", func(t *testing.T) {
		destDir := domain.GetDestinationDirWithTag(nil, "helloWorld")
		if destDir != "" {
			t.Fatalf("expected empty string, got %s", destDir)
		}
	})
}
