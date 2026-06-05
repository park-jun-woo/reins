package temporal

// Kind is whether a time spec is an absolute date (in some calendar) or relative to
// a reference instant (offset in days).
type Kind string

const (
	Absolute Kind = "absolute"
	Relative Kind = "relative"
)

// Calendar identifies the calendar an absolute spec's components are stated in. Only
// Gregorian is converted in v1; the others return Determined=false until v2 adds
// conversion tables.
type Calendar string

const (
	Gregorian Calendar = "gregorian"
	Persian   Calendar = "persian"
	Islamic   Calendar = "islamic"
	Hebrew    Calendar = "hebrew"
	Chinese   Calendar = "chinese"
)

// Spec is the structured time description a gate rule fills from an AI's reading. The
// verifier produces the normalized Value; the AI supplies the identification only.
type Spec struct {
	Kind       Kind     `json:"kind"`
	Calendar   Calendar `json:"calendar,omitempty"`
	Start      string   `json:"start,omitempty"`
	End        string   `json:"end,omitempty"`
	OffsetDays int      `json:"offset_days,omitempty"`
	Anchors    []string `json:"anchors,omitempty"`
}

// Result is the outcome of Resolve: the normalized Gregorian ISO Value (a single
// date, or "start/end" when IsInterval), and whether the spec could be Determined.
type Result struct {
	Value      string `json:"value,omitempty"`
	IsInterval bool   `json:"is_interval,omitempty"`
	Determined bool   `json:"determined"`
}
