package filepathz

import "testing"

func TestExtractShortPath(t *testing.T) {
	t.Parallel()

	t.Run("success,path", func(t *testing.T) {
		t.Parallel()

		const (
			input  = "/path/to/directory/file"
			expect = "directory/file"
		)

		actual := ExtractShortPath(input)
		if expect != actual {
			t.Errorf("❌: expect(%s) != actual(%s)", expect, actual)
		}
	})

	t.Run("success,file_only", func(t *testing.T) {
		t.Parallel()

		const (
			input  = "file"
			expect = "file"
		)

		actual := ExtractShortPath(input)
		if expect != actual {
			t.Errorf("❌: expect(%s) != actual(%s)", expect, actual)
		}
	})

	t.Run("success,directory_file_only", func(t *testing.T) {
		t.Parallel()

		const (
			input  = "directory/file"
			expect = "directory/file"
		)

		actual := ExtractShortPath(input)
		if expect != actual {
			t.Errorf("❌: expect(%s) != actual(%s)", expect, actual)
		}
	})
}
