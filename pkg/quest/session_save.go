//ff:func feature=quest type=helper control=sequence level=error
//ff:what 세션을 들여쓴 JSON으로 직렬화해 파일에 쓴다(diff 친화).

package quest

import (
	"encoding/json"
	"os"
)

// Save writes the session to a JSON file (pretty-printed for diff-friendliness).
func (s *Session) Save(path string) error {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}
