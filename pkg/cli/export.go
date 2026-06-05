//ff:func feature=cli type=command control=sequence level=error
//ff:what `export` 명령. quest.Export로 종단 아이템을 JSONL sink에 증분 방출(원본 보존)하고, 방출 래칫을 저장한 뒤 새로 방출한 레코드 수를 알린다.

package cli

import (
	"fmt"

	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/spf13/cobra"
)

// newExportCmd emits terminal items to the JSONL sink and saves the export ratchet.
func newExportCmd(sessionPath, outPath *string, load sessionLoader) *cobra.Command {
	return &cobra.Command{
		Use:   "export",
		Short: "export terminal results to JSONL (originals preserved)",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := load()
			if err != nil {
				return err
			}
			sink, err := newJSONLSink(*outPath)
			if err != nil {
				return err
			}
			n, err := quest.Export(s, sink)
			if err != nil {
				return err
			}
			if err := s.Save(*sessionPath); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "exported %d new record(s) to %s\n", n, *outPath)
			return nil
		},
	}
}
