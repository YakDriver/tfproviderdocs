package check

import (
	"os"
	"testing"
)

func TestFileSizeCheck(t *testing.T) {
	testCases := []struct {
		Name        string
		Size        int64
		ExpectError bool
	}{
		{
			Name: "under limit",
			Size: RegistryMaximumSizeOfFile - 1,
		},
		{
			Name:        "on limit",
			Size:        RegistryMaximumSizeOfFile,
			ExpectError: true,
		},
		{
			Name:        "over limit",
			Size:        RegistryMaximumSizeOfFile + 1,
			ExpectError: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			file, err := os.CreateTemp(os.TempDir(), "TestFileSizeCheck")

			if err != nil {
				t.Fatalf("error creating temporary file: %s", err)
			}

			defer os.Remove(file.Name())

			if err := file.Truncate(testCase.Size); err != nil {
				t.Fatalf("error writing temporary file: %s", err)
			}

			got := FileSizeCheck(file.Name())

			if got == nil && testCase.ExpectError {
				t.Errorf("expected error, got no error")
			}

			if got != nil && !testCase.ExpectError {
				t.Errorf("expected no error, got error: %s", got)
			}
		})
	}
}

func TestFullPath(t *testing.T) {
	testCases := []struct {
		Name        string
		FileOptions *FileOptions
		Path        string
		Expect      string
	}{
		{
			Name:        "without base path",
			FileOptions: &FileOptions{},
			Path:        "docs/resources/thing.md",
			Expect:      "docs/resources/thing.md",
		},
		{
			Name: "without base path",
			FileOptions: &FileOptions{
				BasePath: "/full/path/to",
			},
			Path:   "docs/resources/thing.md",
			Expect: "/full/path/to/docs/resources/thing.md",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			got := testCase.FileOptions.FullPath(testCase.Path)
			want := testCase.Expect

			if got != want {
				t.Errorf("expected %s, got %s", want, got)
			}
		})
	}
}

func TestFileIgnoreCheck(t *testing.T) {
	testCases := []struct {
		Name   string
		Path   string
		Expect bool
	}{
		{
			Name:   "ignore .DS_Store file",
			Path:   "/path/to/.DS_Store",
			Expect: true,
		},
		{
			Name:   "do not ignore other files",
			Path:   "/path/to/otherfile",
			Expect: false,
		},
		{
			Name:   "ignore .DS_Store file in nested directory",
			Path:   "/another/path/.DS_Store",
			Expect: true,
		},
		{
			Name:   "do not ignore hidden files other than .DS_Store",
			Path:   "/path/to/.hiddenfile",
			Expect: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			got := FileIgnoreCheck(testCase.Path)
			want := testCase.Expect

			if got != want {
				t.Errorf("expected %v, got %v", want, got)
			}
		})
	}
}
