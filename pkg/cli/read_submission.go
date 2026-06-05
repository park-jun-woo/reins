//ff:func feature=cli type=helper control=sequence level=error
//ff:what 제출물 raw 바이트를 읽는다. path가 ""이거나 "-"면 stdin(cmd.InOrStdin)에서, 그 외엔 파일 경로에서 읽는다.

package cli

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

// readSubmission reads the raw submission from a file path, or from stdin when path
// is "-" (or empty).
func readSubmission(cmd *cobra.Command, path string) ([]byte, error) {
	if path == "" || path == "-" {
		return io.ReadAll(cmd.InOrStdin())
	}
	return os.ReadFile(path)
}
