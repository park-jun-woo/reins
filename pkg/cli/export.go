//ff:func feature=cli type=command control=sequence level=error
//ff:what `export` 명령. exportAndSave로 종단 아이템을 JSONL sink에 증분 방출(원본 보존)하고 — Emit이 중간에 실패해도 세션을 Save해 기방출 Emitted 래칫을 보존 — 새로 방출한 레코드 수를 알린다.

package cli

import (
	"fmt"

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
			n, err := exportAndSave(s, sink, *sessionPath)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "exported %d new record(s) to %s\n", n, *outPath)
			return nil
		},
	}
}
