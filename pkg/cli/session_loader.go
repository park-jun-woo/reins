//ff:type feature=cli type=model
//ff:what 작업 세션을 로드(없으면 생성)하는 함수 타입. 모든 서브커맨드가 공유해 세션 IO를 한 곳에 모은다.

package cli

import "github.com/park-jun-woo/reins/pkg/quest"

// sessionLoader loads (or creates) the working session; shared by every subcommand.
type sessionLoader func() (*quest.Session, error)
