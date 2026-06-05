//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what readSubmission이 ""·"-" 경로에서 cmd.InOrStdin으로 읽는지 검증한다.

package cli

import (
	"testing"
)

// TestReadSubmissionStdin: an empty or "-" path reads from cmd.InOrStdin.
func TestReadSubmissionStdin(t *testing.T) {
	for _, path := range []string{"", "-"} {
		got, err := readSubmission(newReadCmd("from stdin"), path)
		if err != nil {
			t.Fatalf("readSubmission(%q): %v", path, err)
		}
		if string(got) != "from stdin" {
			t.Fatalf("readSubmission(%q) = %q, want stdin content", path, got)
		}
	}
}
