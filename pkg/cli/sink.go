//ff:type feature=cli type=model
//ff:what 기본 export sink. quest.Sink을 구현해 아이템 하나당 JSON 한 줄을 파일에 append(JSONL)한다. quest.Export가 종단·미방출 아이템만 방출하므로 증분·멱등이다.

package cli

import "github.com/park-jun-woo/reins/pkg/quest"

// jsonlSink is a quest.Sink that appends one JSON line per item to a file (JSONL).
// It is the default export sink; export is incremental and idempotent because
// quest.Export only emits terminal, not-yet-emitted items.
type jsonlSink struct {
	path string
}

var _ quest.Sink = (*jsonlSink)(nil)
