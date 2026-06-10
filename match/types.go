package match

// Build represents one MinGW-w64 binary distribution.
type Build struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	Arch      string `json:"arch"`
	Thread    string `json:"thread"`
	Exception string `json:"exception"`
	CRT       string `json:"crt"`
	GCC       string `json:"gcc"`
	Priority  int    `json:"priority"`
}

// Rules define matching preferences read from builds.json.
type Rules struct {
	ArchMap             map[string]string   `json:"arch_map"`
	ThreadPreference    []string            `json:"thread_preference"`
	ExceptionPreference map[string][]string `json:"exception_preference"`
	CRTPreference       []string            `json:"crt_preference"`
}

// DimensionChoice explains why a particular value was chosen for one dimension.
type DimensionChoice struct {
	Dimension string `json:"dimension"`
	Choice    string `json:"choice"`
	Reason    string `json:"reason"`
	Manual    bool   `json:"manual"`
}

// MatchResult holds the best build, alternatives, and explanations.
type MatchResult struct {
	Build        Build             `json:"build"`
	Alternatives []Build           `json:"alternatives,omitempty"`
	Explanation  []DimensionChoice `json:"explanation"`
	Score        int               `json:"score"`
}
