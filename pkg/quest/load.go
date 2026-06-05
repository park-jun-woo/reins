//ff:func feature=quest type=helper control=sequence level=error
//ff:what 세션 JSON 파일을 읽어 Session으로 역직렬화한다. 부재는 os.IsNotExist(err)로 보고한다.

package quest

import (
	"encoding/json"
	"fmt"
	"os"
)

// Load reads a session from a JSON file. A missing file is reported via
// os.IsNotExist(err) so callers can create a new session.
func Load(path string) (*Session, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var s Session
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, fmt.Errorf("session %s: %w", path, err)
	}
	return &s, nil
}
