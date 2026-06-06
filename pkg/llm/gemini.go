//ff:type feature=llm type=adapter
//ff:what Gemini — Google Gemini generateContent endpoint용 Backend. Model과 BaseURL(빈 값이면 generativelanguage 호스트; 테스트 주입 seam)을 갖는다. API 키는 env GEMINI_API_KEY.

package llm

// defaultGeminiBaseURL is the Google generativelanguage host used when BaseURL is empty.
const defaultGeminiBaseURL = "https://generativelanguage.googleapis.com"

// Gemini is a Backend for the Google Gemini generateContent endpoint.
type Gemini struct {
	Model   string
	BaseURL string // empty ⇒ the Google generativelanguage host (test injection seam)
}
