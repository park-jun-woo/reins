//ff:func feature=cli type=helper control=sequence
//ff:what 테스트 헬퍼. 주어진 stdin을 가진 빈 cobra 명령을 만든다(readSubmission의 stdin 분기 자극용).

package cli

import (
	"strings"

	"github.com/spf13/cobra"
)

// newReadCmd builds a bare cobra command with the given stdin, for exercising
// readSubmission's stdin branch via cmd.InOrStdin.
func newReadCmd(in string) *cobra.Command {
	cmd := &cobra.Command{}
	cmd.SetIn(strings.NewReader(in))
	return cmd
}
