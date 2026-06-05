//ff:func feature=graph type=helper control=selection
//ff:what Validate — 빈 ID를 거부한다(Spec 인터페이스 요구).

package graph

import "errors"

// Validate rejects an empty ID.
func (s idSpec) Validate() error {
	if s.ID == "" {
		return errors.New("graph: node id must not be empty")
	}
	return nil
}
